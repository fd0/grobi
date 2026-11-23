package main

import (
	"strings"

	"golang.org/x/sys/unix"
)

const (
	devdSocket = "/var/run/devd.seqpacket.pipe"

	lidSysctl       = "dev.acpi_lid.0.state"
	lidSysctlClosed = 0 // any other value means open

	devdSystemACPI   = "system=ACPI"
	devdSubsystemLid = "subsystem=Lid"
	devdNotifyClosed = "notify=0x00"
	devdNotifyOpen   = "notify=0x01"
)

// readLidState reads the lid sysctl.
// Not verified on real hardware; without it the last devd event is used.
func readLidState() LidState {
	v, err := unix.SysctlUint32(lidSysctl)
	if err != nil {
		return LidUnknown
	}
	if v != lidSysctlClosed {
		return LidOpen
	}
	return LidClosed
}

// parseDevdEvent parses lines like
// "!system=ACPI subsystem=Lid type=\_SB_.LID_ notify=0x00".
func parseDevdEvent(line string) LidState {
	var acpi, lid bool
	l := LidUnknown
	for _, f := range strings.Fields(strings.TrimPrefix(line, "!")) {
		switch f {
		case devdSystemACPI:
			acpi = true
		case devdSubsystemLid:
			lid = true
		case devdNotifyClosed:
			l = LidClosed
		case devdNotifyOpen:
			l = LidOpen
		}
	}
	if !acpi || !lid {
		return LidUnknown
	}
	return l
}

// subscribeLidEvents reads lid events from devd.
func subscribeLidEvents(ch chan<- struct{}, done <-chan struct{}) {
	watchLidSocket("unixpacket", devdSocket, parseDevdEvent, ch, done)
}
