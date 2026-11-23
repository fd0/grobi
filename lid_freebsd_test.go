package main

import "testing"

func TestParseDevdEvent(t *testing.T) {
	for line, want := range map[string]LidState{
		`!system=ACPI subsystem=Lid type=\_SB_.LID_ notify=0x00`:    LidClosed,
		`!system=ACPI subsystem=Lid type=\_SB_.LID_ notify=0x01`:    LidOpen,
		`!system=ACPI subsystem=Button type=\_SB_.PWRB notify=0x00`: LidUnknown,
		`+uhid0 at bus=0`: LidUnknown,
	} {
		if got := parseDevdEvent(line); got != want {
			t.Errorf("%q: wanted %v, got %v", line, want, got)
		}
	}
}
