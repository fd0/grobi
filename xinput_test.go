package main

import (
	"reflect"
	"testing"

	"github.com/jezek/xgb"
	"go.yaml.in/yaml/v4"
)

// queryDeviceReply builds an XIQueryDevice reply.
func queryDeviceReply(devs []struct {
	use     uint16
	name    string
	classes int
}) []byte {
	const (
		xReply      = 1 // reply type
		replyLength = 4 // in xUnit, after the 32 byte header
		classSize   = 3 * xUnit
	)

	buf := make([]byte, queryDeviceReplySize)
	buf[0] = xReply
	xgb.Put16(buf[queryDeviceReplyNumDevices:], uint16(len(devs)))
	for _, d := range devs {
		info := make([]byte, deviceInfoSize+xgb.Pad(len(d.name)))
		xgb.Put16(info[deviceInfoUse:], d.use)
		xgb.Put16(info[deviceInfoNumClasses:], uint16(d.classes))
		xgb.Put16(info[deviceInfoNameLen:], uint16(len(d.name)))
		copy(info[deviceInfoSize:], d.name)
		buf = append(buf, info...)

		for range d.classes {
			class := make([]byte, classSize)
			xgb.Put16(class[classLength:], classSize/xUnit)
			buf = append(buf, class...)
		}
	}
	xgb.Put32(buf[replyLength:], uint32((len(buf)-queryDeviceReplySize)/xUnit))
	return buf
}

func TestParseQueryDeviceReply(t *testing.T) {
	reply := queryDeviceReply([]struct {
		use     uint16
		name    string
		classes int
	}{
		{1, "Virtual core pointer", 3},
		{2, "Virtual core keyboard", 1},
		{3, "Virtual core XTEST pointer", 2},
		{4, "Lenovo ThinkPad Thunderbolt 3 Dock USB Audio", 1},
		{4, "HD Pro Webcam C920", 1},
		{5, "Wacom", 0},
	})

	devs, err := parseQueryDeviceReply(reply)
	if err != nil {
		t.Fatal(err)
	}
	want := InputDevices{
		"Virtual core XTEST pointer",
		"Lenovo ThinkPad Thunderbolt 3 Dock USB Audio",
		"HD Pro Webcam C920",
		"Wacom",
	}
	if !reflect.DeepEqual(devs, want) {
		t.Errorf("wanted %q, got %q", want, devs)
	}

	for n := 0; n < len(reply); n++ {
		if _, err := parseQueryDeviceReply(reply[:n]); err == nil {
			t.Errorf("truncated reply of %d bytes: no error", n)
		}
	}
}

func TestParseQueryDeviceReplyEmpty(t *testing.T) {
	devs, err := parseQueryDeviceReply(queryDeviceReply(nil))
	if err != nil {
		t.Fatal(err)
	}
	if devs == nil || len(devs) != 0 {
		t.Errorf("wanted empty non-nil list, got %#v", devs)
	}
}

func TestRuleMatchInputs(t *testing.T) {
	inputs := InputDevices{"Lenovo ThinkPad Thunderbolt 3 Dock USB Audio", "HD Pro Webcam C920"}
	tests := []struct {
		present, absent []string
		inputs          InputDevices
		match           bool
	}{
		{nil, nil, nil, true},
		{nil, nil, inputs, true},
		{[]string{"*Webcam C920"}, nil, inputs, true},
		{[]string{"*Webcam C920", "*Dock*"}, nil, inputs, true},
		{[]string{"*Webcam C920", "Yubikey*"}, nil, inputs, false},
		{nil, []string{"*Dock*"}, inputs, false},
		{nil, []string{"Yubikey*"}, inputs, true},
		{nil, []string{"Yubikey*"}, InputDevices{}, true},
		{[]string{"*Webcam C920"}, nil, nil, false},
		{nil, []string{"Yubikey*"}, nil, false},
	}

	for i, test := range tests {
		r := Rule{OutputsConnected: []string{"HDMI"}, InputsPresent: test.present, InputsAbsent: test.absent}
		m := r.Match(State{Outputs: testOutputs, Inputs: test.inputs})
		if m != test.match {
			t.Errorf("test %d: wanted %v, got %v", i, test.match, m)
		}
	}
}

// dockRules tells a desk with a dock apart from the same monitor plugged in
// directly. Monitor IDs and device names are made up.
const dockRules = `
rules:
  - name: desk-lid-closed
    lid: closed
    inputs_present: ["ACME*Dock*"]
    outputs_connected:
      - DP-1-AAA-1111-0-Monitor-SERIAL0001
      - DP-2-AAA-1111-0-Monitor-SERIAL0002
    configure_row: [DP-1, DP-2]

  - name: desk-lid-open
    lid: open
    inputs_present: ["ACME*Dock*"]
    outputs_connected:
      - eDP-1-BBB-2222-*
      - DP-1-AAA-1111-*-SERIAL0001
      - DP-2-AAA-1111-*-SERIAL0002
    configure_row: [DP-1, DP-2, eDP-1]

  - name: hdmi-lid-closed
    lid: closed
    inputs_absent: ["ACME*Dock*"]
    outputs_connected: [HDMI-1-AAA-1111-*]
    configure_single: HDMI-1

  - name: hdmi-lid-open
    lid: open
    inputs_absent: ["ACME*Dock*"]
    outputs_connected: [eDP-1-BBB-2222-*, HDMI-1-AAA-1111-*]
    configure_row: [HDMI-1, eDP-1]

  - name: laptop
    configure_single: eDP-1
`

func TestDockRules(t *testing.T) {
	var cfg Config
	if err := yaml.Unmarshal([]byte(dockRules), &cfg); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Valid(); err != nil {
		t.Fatal(err)
	}

	panel := Output{Name: "eDP-1", Connected: true, MonitorID: "BBB-2222-0--"}
	left := Output{Name: "DP-1", Connected: true, MonitorID: "AAA-1111-0-Monitor-SERIAL0001"}
	right := Output{Name: "DP-2", Connected: true, MonitorID: "AAA-1111-0-Monitor-SERIAL0002"}
	hdmiLeft := Output{Name: "HDMI-1", Connected: true, MonitorID: left.MonitorID}
	hdmiRight := Output{Name: "HDMI-1", Connected: true, MonitorID: right.MonitorID}
	off := func(name string) Output { return Output{Name: name} }
	docked := Outputs{panel, left, right, off("HDMI-1")}

	laptop := InputDevices{"Virtual core XTEST pointer", "Example Touchpad", "Power Button", "AT Translated Set 2 keyboard"}
	withDock := append(append(InputDevices{}, laptop...),
		"Example Mouse", "Example Webcam Consumer Control", "ACME,Inc. ACME Dock Audio - Rear", "USB Keyboard")

	tests := []struct {
		name    string
		outputs Outputs
		lid     LidState
		inputs  InputDevices
		want    string
	}{
		{"dock, lid closed", docked, LidClosed, withDock, "desk-lid-closed"},
		{"dock, lid open", docked, LidOpen, withDock, "desk-lid-open"},
		{"monitors, no dock", docked, LidOpen, laptop, "laptop"},
		{"hdmi left, lid closed", Outputs{panel, off("DP-1"), hdmiLeft}, LidClosed, laptop, "hdmi-lid-closed"},
		{"hdmi right, lid open", Outputs{panel, off("DP-1"), hdmiRight}, LidOpen, laptop, "hdmi-lid-open"},
		{"hdmi with dock", Outputs{panel, hdmiLeft}, LidOpen, withDock, "laptop"},
		{"inputs unknown", docked, LidOpen, nil, "laptop"},
		{"lid unknown", docked, LidUnknown, withDock, "laptop"},
		{"mobile", Outputs{panel, off("DP-1"), off("HDMI-1")}, LidOpen, laptop, "laptop"},
	}

	for _, test := range tests {
		r, err := MatchRules(cfg.Rules, State{Outputs: test.outputs, Lid: test.lid, Inputs: test.inputs})
		if err != nil {
			t.Errorf("%s: %v", test.name, err)
			continue
		}
		if r.Name != test.want {
			t.Errorf("%s: want rule %s, got %s", test.name, test.want, r.Name)
		}
	}
}

func TestConfigValidInputs(t *testing.T) {
	err := Config{Rules: []Rule{{Name: "x", InputsPresent: []string{"[abc"}}}}.Valid()
	if err == nil {
		t.Error("malformed inputs_present pattern accepted")
	}
}
