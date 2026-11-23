package main

import "testing"

func TestParseAcpidEvent(t *testing.T) {
	for line, want := range map[string]LidState{
		"button/lid LID close":                LidClosed,
		"button/lid LID open":                 LidOpen,
		"button/lid LID0 00000080 00000001":   LidUnknown,
		"button/power PBTN 00000080 00000000": LidUnknown,
		"":                                    LidUnknown,
	} {
		if got := parseAcpidEvent(line); got != want {
			t.Errorf("%q: wanted %v, got %v", line, want, got)
		}
	}
}
