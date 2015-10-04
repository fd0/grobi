package main

import (
	"bytes"
	"reflect"
	"testing"
)

var randrTestOutputs = []struct {
	str     string
	outputs []Output
}{
	{
		`Screen 0: minimum 320 x 200, current 3280 x 1200, maximum 8192 x 8192
LVDS1 connected (normal left inverted right x axis y axis)
   1366x768      60.10 +
   1024x768      60.00
   800x600       60.32    56.25
   640x480       59.94
VGA1 disconnected (normal left inverted right x axis y axis)
HDMI1 disconnected (normal left inverted right x axis y axis)
DP1 disconnected (normal left inverted right x axis y axis)
HDMI2 connected 1600x1200+0+0 (normal left inverted right x axis y axis) 408mm x 306mm
   1600x1200     60.00*+
   1280x1024     75.02    60.02
   1280x960      60.00
   1152x864      75.00
   1024x768      75.08    70.07    60.00
   832x624       74.55
   800x600       72.19    75.00    60.32    56.25
   640x480       75.00    72.81    66.67    60.00
   720x400       70.08
DP2 disconnected (normal left inverted right x axis y axis)
DP3 disconnected (normal left inverted right x axis y axis)`,
		[]Output{
			Output{
				Name: "LVDS1",
				Modes: []Mode{
					{Name: "1366x768", Default: true},
					{Name: "1024x768"},
					{Name: "800x600"},
					{Name: "640x480"},
				},
				Connected: true,
			},
			Output{Name: "VGA1"},
			Output{Name: "HDMI1"},
			Output{Name: "DP1"},
			Output{
				Name: "HDMI2",
				Modes: []Mode{
					{Name: "1600x1200", Default: true, Active: true},
					{Name: "1280x1024"},
					{Name: "1280x960"},
					{Name: "1152x864"},
					{Name: "1024x768"},
					{Name: "832x624"},
					{Name: "800x600"},
					{Name: "640x480"},
					{Name: "720x400"},
				},
				Connected: true,
			},
			Output{Name: "DP2"},
			Output{Name: "DP3"},
		},
	},
	{
		`Screen 0: minimum 320 x 200, current 3280 x 1200, maximum 8192 x 8192
LVDS1 connected (normal left inverted right x axis y axis)
   1366x768      60.10 +
   1024x768      60.00
   800x600       60.32    56.25
   640x480       59.94`,
		[]Output{
			Output{
				Name: "LVDS1",
				Modes: []Mode{
					{Name: "1366x768", Default: true},
					{Name: "1024x768"},
					{Name: "800x600"},
					{Name: "640x480"},
				},
				Connected: true,
			},
		},
	},
	{
		`Screen 0: minimum 8 x 8, current 4480 x 1440, maximum 32767 x 32767
eDP1 connected 1920x1080+2560+0 (normal left inverted right x axis y axis) 276mm x 156mm
   1920x1080     60.04*+
   1400x1050     59.98
   1600x900      60.00
   1280x1024     60.02
   1280x960      60.00
   1368x768      60.00
   1280x720      60.00
   1024x768      60.00
   1024x576      60.00
   960x540       60.00
   800x600       60.32    56.25
   864x486       60.00
   640x480       59.94
   720x405       60.00
   640x360       60.00
DP1 disconnected (normal left inverted right x axis y axis)
DP2 disconnected (normal left inverted right x axis y axis)
DP2-1 disconnected (normal left inverted right x axis y axis)
DP2-2 connected primary 2560x1440+0+0 (normal left inverted right x axis y axis) 597mm x 336mm
   2560x1440     59.95*+
   2048x1152     60.00
   1920x1200     59.88
   1920x1080     60.00    50.00    59.94    30.00    25.00    24.00    29.97    23.98
   1600x1200     60.00
   1680x1050     59.95
   1280x1024     75.02    60.02
   1200x960      59.99
   1152x864      75.00
   1280x720      60.00    50.00    59.94
   1024x768      75.08    60.00
   800x600       75.00    60.32
   720x576       50.00
   720x480       60.00    59.94
   640x480       75.00    60.00    59.94
   720x400       70.08
DP2-3 disconnected (normal left inverted right x axis y axis)
HDMI1 disconnected (normal left inverted right x axis y axis)
HDMI2 disconnected (normal left inverted right x axis y axis)
VIRTUAL1 disconnected (normal left inverted right x axis y axis)`,
		[]Output{
			Output{
				Name: "eDP1",
				Modes: []Mode{
					{Name: "1920x1080", Default: true, Active: true},
					{Name: "1400x1050"},
					{Name: "1600x900"},
					{Name: "1280x1024"},
					{Name: "1280x960"},
					{Name: "1368x768"},
					{Name: "1280x720"},
					{Name: "1024x768"},
					{Name: "1024x576"},
					{Name: "960x540"},
					{Name: "800x600"},
					{Name: "864x486"},
					{Name: "640x480"},
					{Name: "720x405"},
					{Name: "640x360"},
				},
				Connected: true,
			},
			Output{Name: "DP1"},
			Output{Name: "DP2"},
			Output{Name: "DP2-1"},
			Output{
				Name: "DP2-2",
				Modes: []Mode{
					{Name: "2560x1440", Default: true, Active: true},
					{Name: "2048x1152"},
					{Name: "1920x1200"},
					{Name: "1920x1080"},
					{Name: "1600x1200"},
					{Name: "1680x1050"},
					{Name: "1280x1024"},
					{Name: "1200x960"},
					{Name: "1152x864"},
					{Name: "1280x720"},
					{Name: "1024x768"},
					{Name: "800x600"},
					{Name: "720x576"},
					{Name: "720x480"},
					{Name: "640x480"},
					{Name: "720x400"},
				},
				Connected: true,
			},
			Output{Name: "DP2-3"},
			Output{Name: "HDMI1"},
			Output{Name: "HDMI2"},
			Output{Name: "VIRTUAL1"},
		},
	},
	{
		`Screen 0: minimum 320 x 200, current 3280 x 1200, maximum 8192 x 8192
LVDS1 connected (normal left inverted right x axis y axis)
   1366x768      60.10 +
   1024x768      60.00  
   800x600       60.32    56.25  
   640x480       59.94  
HDMI2 disconnected 1600x1200+0+0 (normal left inverted right x axis y axis) 0mm x 0mm
HDMI3 disconnected 1680x1050+1600+0 (normal left inverted right x axis y axis) 0mm x 0mm`,
		[]Output{
			Output{
				Name: "LVDS1",
				Modes: []Mode{
					{Name: "1366x768", Default: true},
					{Name: "1024x768"},
					{Name: "800x600"},
					{Name: "640x480"},
				},
				Connected: true,
			},
			Output{
				Name: "HDMI2",
				Modes: []Mode{
					{Name: "1600x1200", Active: true},
				},
			},
			Output{
				Name: "HDMI3",
				Modes: []Mode{
					{Name: "1680x1050", Active: true},
				},
			},
		},
	},
}

func TestRandrParse(t *testing.T) {
	for _, test := range randrTestOutputs {
		out, err := RandrParse(bytes.NewReader([]byte(test.str)))
		if err != nil {
			t.Errorf("unable to parse test string: %v", err)
			continue
		}

		if len(test.outputs) != len(out) {
			t.Errorf("number of connected outputs does not mach: expected %d, got %d",
				len(test.outputs), len(out))
			continue
		}

		for i := range test.outputs {
			out1 := test.outputs[i]
			out2 := out[i]

			if out1.Name != out2.Name {
				t.Errorf("output %d: name not equal: want %q, got %q", i,
					out1.Name, out2.Name)
			}

			if out1.Connected != out2.Connected {
				t.Errorf("output %d: connected not equal: want %v, got %v", i,
					out1.Connected, out2.Connected)
			}

			if !reflect.DeepEqual(out1.Modes, out2.Modes) {
				t.Errorf("output %d: list of modes not equal: want %v, got %v", i,
					out1.Modes, out2.Modes)
			}
		}
	}
}

var TestOutputLines = []struct {
	line   string
	output Output
}{
	{
		"LVDS1 connected (normal left inverted right x axis y axis)",
		Output{
			Name:      "LVDS1",
			Connected: true,
		},
	},
	{
		"VGA1 disconnected (normal left inverted right x axis y axis)",
		Output{
			Name:      "VGA1",
			Connected: false,
		},
	},
	{
		"HDMI3 disconnected 1680x1050+1600+0 (normal left inverted right x axis y axis) 0mm x 0mm`",
		Output{
			Name:  "HDMI3",
			Modes: []Mode{{Name: "1680x1050", Active: true}},
		},
	},
}

func TestParseOutputLine(t *testing.T) {
	for i, test := range TestOutputLines {
		out, err := parseOutputLine(test.line)
		if err != nil {
			t.Errorf("test %d returned error: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(out, test.output) {
			t.Errorf("test %d failed: expected Output not found", i)
			continue
		}
	}
}

var TestModeLines = []struct {
	line string
	mode Mode
}{
	{
		"  1152x864      75.00",
		Mode{
			Name: "1152x864",
		},
	},
	{
		"  1024x768      75.08    70.07    60.00",
		Mode{
			Name: "1024x768",
		},
	},
	{
		"  1600x1200     60.00*+",
		Mode{
			Name:    "1600x1200",
			Active:  true,
			Default: true,
		},
	},
	{
		"  1366x768      60.10 +",
		Mode{
			Name:    "1366x768",
			Default: true,
		},
	},
	{
		"  832x624       74.55",
		Mode{
			Name: "832x624",
		},
	},
}

func TestParseModeLine(t *testing.T) {
	for i, test := range TestModeLines {
		mode, err := parseModeLine(test.line)
		if err != nil {
			t.Errorf("test %d returned error: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(mode, test.mode) {
			t.Errorf("test %d failed: expected Mode not found", i)
			continue
		}
	}
}
