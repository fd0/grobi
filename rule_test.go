package main

import "testing"

var testRules = []struct {
	rule  Rule
	match bool
}{
	{
		Rule{
			OutputsConnected: []string{"HDMI", "VGA"},
			OutputsAbsent:    []string{"DP2-2"},
		},
		true,
	},
	{
		Rule{
			OutputsConnected: []string{"LVDS1"},
			OutputsAbsent:    []string{"HDMI"},
		},
		false,
	},
	{
		Rule{
			OutputsConnected:    []string{"LVDS1"},
			OutputsDisconnected: []string{"HDMI"},
		},
		false,
	},
	{
		Rule{
			OutputsPresent: []string{"DP2-1"},
		},
		true,
	},
	{
		Rule{
			OutputsPresent: []string{"DP2-1"},
		},
		true,
	},
	{
		Rule{
			OutputsConnected:    []string{"HDMI*", "VGA"},
			OutputsDisconnected: []string{"DP2-?"},
		},
		true,
	},
	{
		Rule{
			OutputsConnected: []string{"HDMI*", "VGA"},
			OutputsAbsent:    []string{"DP2-?"},
		},
		false,
	},
	{
		Rule{
			OutputsPresent: []string{"HDMI-SAM-2618-808661557"},
		},
		true,
	},
	{
		Rule{
			OutputsPresent: []string{"*-UNK-123-456"},
		},
		false,
	},
	{
		Rule{
			OutputsDisconnected: []string{"HDMI-UNK-123-456"},
		},
		true,
	},
}

var testOutputs = []Output{
	{
		Name:      "LVDS",
		Connected: true,
		Modes: []Mode{
			{"1377x768", true, true},
			{"1024x768", false, false},
		},
		MonitorID: "CMN-5297-0",
	},
	{
		Name:      "VGA",
		Connected: true,
		Modes: []Mode{
			{"1280x1024", true, false},
			{"1024x768", false, true},
		},
	},
	{
		Name:      "HDMI",
		Connected: true,
		Modes: []Mode{
			{"1920x1080", true, true},
			{"1024x768", false, false},
		},
		MonitorID: "SAM-2618-808661557",
	},
	{
		Name: "DP2-1",
	},
}

func TestRuleMatch(t *testing.T) {
	for i, test := range testRules {
		m := test.rule.Match(testOutputs)
		if m != test.match {
			t.Errorf("test rule %d wrong match: wanted %v, got %v", i, test.match, m)
			continue
		}
	}
}
