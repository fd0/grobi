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
	{`Screen 0: minimum 320 x 200, current 3280 x 1200, maximum 8192 x 8192
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
				Modes: []string{
					"1366x768",
					"1024x768",
					"800x600",
					"640x480",
				},
			},
			Output{
				Name: "HDMI2",
				Modes: []string{
					"1600x1200",
					"1280x1024",
					"1280x960",
					"1152x864",
					"1024x768",
					"832x624",
					"800x600",
					"640x480",
					"720x400",
				},
				Active:     true,
				ActiveMode: "1600x1200",
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

		if !reflect.DeepEqual(out, test.outputs) {
			t.Errorf("result is not equal to expected result")
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
			Name:   "LVDS1",
			Active: true,
		},
	},
	{
		"VGA1 disconnected (normal left inverted right x axis y axis)",
		Output{
			Name:   "VGA1",
			Active: false,
		},
	},
}

func TestParseOutputLine(t *testing.T) {
	for i, test := range TestOutputLines {
		out, err := ParseOutputLine(test.line)
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
