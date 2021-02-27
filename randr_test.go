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
			{
				Name: "LVDS1",
				Modes: []Mode{
					{Name: "1366x768", Default: true},
					{Name: "1024x768"},
					{Name: "800x600"},
					{Name: "640x480"},
				},
				Connected: true,
			},
			{Name: "VGA1"},
			{Name: "HDMI1"},
			{Name: "DP1"},
			{
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
			{Name: "DP2"},
			{Name: "DP3"},
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
			{
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
			{
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
			{Name: "DP1"},
			{Name: "DP2"},
			{Name: "DP2-1"},
			{
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
				Primary:   true,
			},
			{Name: "DP2-3"},
			{Name: "HDMI1"},
			{Name: "HDMI2"},
			{Name: "VIRTUAL1"},
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
			{
				Name: "LVDS1",
				Modes: []Mode{
					{Name: "1366x768", Default: true},
					{Name: "1024x768"},
					{Name: "800x600"},
					{Name: "640x480"},
				},
				Connected: true,
			},
			{
				Name: "HDMI2",
				Modes: []Mode{
					{Name: "1600x1200", Active: true},
				},
			},
			{
				Name: "HDMI3",
				Modes: []Mode{
					{Name: "1680x1050", Active: true},
				},
			},
		},
	},
	{
		`Screen 0: minimum 8 x 8, current 1920 x 1080, maximum 32767 x 32767
eDP1 connected primary 1920x1080+0+0 (normal left inverted right x axis y axis) 308mm x 173mm
   1920x1080     60.01*+
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
DP3 disconnected (normal left inverted right x axis y axis)
HDMI1 disconnected (normal left inverted right x axis y axis)
HDMI2 disconnected (normal left inverted right x axis y axis)
HDMI3 disconnected (normal left inverted right x axis y axis)
VIRTUAL1 disconnected (normal left inverted right x axis y axis)`,
		[]Output{
			{
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
				Primary:   true,
			},
			{
				Name:      "DP1",
				Connected: false,
			},
			{
				Name:      "DP2",
				Connected: false,
			},
			{
				Name:      "DP3",
				Connected: false,
			},
			{
				Name:      "HDMI1",
				Connected: false,
			},
			{
				Name:      "HDMI2",
				Connected: false,
			},
			{
				Name:      "HDMI3",
				Connected: false,
			},
			{
				Name:      "VIRTUAL1",
				Connected: false,
			},
		},
	},
	{
		`Screen 0: minimum 8 x 8, current 3840 x 1080, maximum 32767 x 32767
eDP1 connected primary 1920x1080+1920+0 (normal left inverted right x axis y axis) 310mm x 170mm
	EDID: 
		00ffffffffffff000daeb11400000000
		0c190104951f117802ff359255529529
		25505400000001010101010101010101
		010101010101b43b804a71383440503c
		680034ad10000018000000fe004e3134
		304843452d4541410a20000000fe0043
		4d4e0a202020202020202020000000fe
		004e3134304843452d4541410a2000a2
	BACKLIGHT: 332 
		range: (0, 852)
	Backlight: 332 
		range: (0, 852)
	scaling mode: Full aspect 
		supported: None, Full, Center, Full aspect
	Broadcast RGB: Automatic 
		supported: Automatic, Full, Limited 16:235
	audio: auto 
		supported: force-dvi, off, auto, on
   1920x1080     60.01*+
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
	Broadcast RGB: Automatic 
		supported: Automatic, Full, Limited 16:235
	audio: auto 
		supported: force-dvi, off, auto, on
DP2 disconnected (normal left inverted right x axis y axis)
	Broadcast RGB: Automatic 
		supported: Automatic, Full, Limited 16:235
	audio: auto 
		supported: force-dvi, off, auto, on
HDMI1 disconnected (normal left inverted right x axis y axis)
	aspect ratio: Automatic 
		supported: Automatic, 4:3, 16:9
	Broadcast RGB: Automatic 
		supported: Automatic, Full, Limited 16:235
	audio: auto 
		supported: force-dvi, off, auto, on
HDMI2 connected 1920x1080+0+0 (normal left inverted right x axis y axis) 530mm x 300mm
	EDID: 
		00ffffffffffff004c2d3a0a35323330
		2417010380351e782af711a3564f9e28
		0f5054bfef80714f81c0810081809500
		a9c0b3000101023a801871382d40582c
		4500132b2100001e011d007251d01e20
		6e285500132b2100001e000000fd0032
		4b1e5111000a202020202020000000fc
		00533234433335300a20202020200118
		02031af14690041f1303122309070783
		01000066030c00100080011d00bc52d0
		1e20b8285540132b2100001e8c0ad090
		204031200c405500132b210000188c0a
		d08a20e02d10103e9600132b21000018
		00000000000000000000000000000000
		00000000000000000000000000000000
		00000000000000000000000000000099
	aspect ratio: Automatic 
		supported: Automatic, 4:3, 16:9
	Broadcast RGB: Automatic 
		supported: Automatic, Full, Limited 16:235
	audio: auto 
		supported: force-dvi, off, auto, on
   1920x1080     60.00*+  50.00    59.94  
   1680x1050     59.88  
   1600x900      60.00  
   1280x1024     75.02    60.02  
   1440x900      59.90  
   1280x800      59.91  
   1152x864      75.00  
   1280x720      60.00    50.00    59.94  
   1024x768      75.03    70.07    60.00  
   832x624       74.55  
   800x600       72.19    75.00    60.32    56.25  
   720x576       50.00  
   720x480       60.00    59.94  
   640x480       75.00    72.81    66.67    60.00    59.94  
   720x400       70.08  
VIRTUAL1 disconnected (normal left inverted right x axis y axis)`,
		[]Output{
			{
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
				Primary:   true,
				MonitorID: "CMN-5297-0--",
			},
			{Name: "DP1"},
			{Name: "DP2"},
			{Name: "HDMI1"},
			{Name: "HDMI2",
				Modes: []Mode{
					{Name: "1920x1080", Default: true, Active: true},
					{Name: "1680x1050"},
					{Name: "1600x900"},
					{Name: "1280x1024"},
					{Name: "1440x900"},
					{Name: "1280x800"},
					{Name: "1152x864"},
					{Name: "1280x720"},
					{Name: "1024x768"},
					{Name: "832x624"},
					{Name: "800x600"},
					{Name: "720x576"},
					{Name: "720x480"},
					{Name: "640x480"},
					{Name: "720x400"},
				},
				Connected: true,
				MonitorID: "SAM-2618-808661557-S24C350-",
			},
			{Name: "VIRTUAL1"},
		},
	},
}

func TestRandrParse(t *testing.T) {
	for ti, test := range randrTestOutputs {
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
				t.Errorf("test %d, output %d: name not equal: want %q, got %q",
					ti, i, out1.Name, out2.Name)
			}

			if out1.Connected != out2.Connected {
				t.Errorf("test %d, output %d: connected not equal: want %v, got %v",
					ti, i,
					out1.Connected, out2.Connected)
			}

			if out1.Primary != out2.Primary {
				t.Errorf("test %d, output %d: primary not equal: want %v, got %v",
					ti, i,
					out1.Primary, out2.Primary)
			}

			if !reflect.DeepEqual(out1.Modes, out2.Modes) {
				t.Errorf("test %d, output %d: list of modes not equal:\n  want %v\n  got  %v",
					ti, i,
					out1.Modes, out2.Modes)
			}

			if out1.MonitorID != out2.MonitorID {
				t.Errorf("test %d, output %d: Monitor IDs not equal:\n  want %v\n  got  %v",
					ti, i,
					out1.MonitorID, out2.MonitorID)
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
	{
		"DP3-1-8 connected primary 2560x1440+0+0 (normal left inverted right x axis y axis) 553mm x 311mm",
		Output{
			Name:      "DP3-1-8",
			Connected: true,
			Primary:   true,
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
			t.Errorf("test %d failed: expected Output not found, want %v, got %v",
				i, test.output, out)
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

func TestGenerateMonitorId(t *testing.T) {
	var tests = []struct {
		edid      string
		failure   bool
		monitorID string
	}{
		{
			"00ffffffffffff004c2d3a0a353233302417010380351e782af711a3564f9e280f5054bfef80714f81c0810081809500a9c0b3000101023a801871382d40582c4500132b2100001e011d007251d01e206e285500132b2100001e000000fd00324b1e5111000a202020202020000000fc00533234433335300a2020202020011802031af14690041f130312230907078301000066030c00100080011d00bc52d01e20b8285540132b2100001e8c0ad090204031200c405500132b210000188c0ad08a20e02d10103e9600132b21000018000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000099",
			false,
			"SAM-2618-808661557-S24C350-",
		},
		{
			"00ffffffffffff000daeb114000000000c190104951f117802ff35925552952925505400000001010101010101010101010101010101b43b804a71383440503c680034ad10000018000000fe004e3134304843452d4541410a20000000fe00434d4e0a202020202020202020000000fe004e3134304843452d4541410a2000a2",
			false,
			"CMN-5297-0--",
		},
		{
			"",
			true,
			"",
		},
		{
			"00ffffffffffff004c2d3a0a3532333",
			true,
			"",
		},
		{
			"00ffffeeffffff000daeb114000000000c190104951f117802ff35925552952925505400000001010101010101010101010101010101b43b804a71383440503c680034ad10000018000000fe004e3134304843452d4541410a20000000fe00434d4e0a202020202020202020000000fe004e3134304843452d4541410a2000a2",
			true,
			"",
		},
		{
			"00ffffffffffff004c2dd4043432494b0b14010380351e782a6041a6564a9c251250542308008100814081809500a940b300010101011a3680a070381f4030203500132a2100001a000000fd00383c1e5110000a202020202020000000fc0053796e634d61737465720a2020000000ff004839585a3330353131380a202000f7",
			false,
			"SAM-1236-1263088180-SyncMaster-H9XZ305118",
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			monitorID, err := GenerateMonitorID(test.edid)
			if test.failure {
				if err == nil {
					t.Fatal("did not return an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			if monitorID != test.monitorID {
				t.Fatalf("wanted monitor ID %q, got %q", test.monitorID, monitorID)
			}
		})
	}
}
