package main

import (
	"os/exec"
	"strings"
)

// readLidState scans the acpibtn sensors for the lid, e.g.
// "hw.sensors.acpibtn1.indicator0=On (lid open)". The sensor index varies.
// shortcut: forks sysctl per poll; not verified on real hardware.
func readLidState() LidState {
	out, err := exec.Command("sysctl", "hw.sensors").Output()
	if err != nil {
		return LidUnknown
	}
	return parseOpenBSDSensors(string(out))
}

func parseOpenBSDSensors(s string) LidState {
	for _, line := range strings.Split(s, "\n") {
		if !strings.Contains(line, ".acpibtn") || !strings.Contains(line, "(lid open)") {
			continue
		}
		_, v, _ := strings.Cut(line, "=")
		switch {
		case strings.HasPrefix(v, "On"):
			return LidOpen
		case strings.HasPrefix(v, "Off"):
			return LidClosed
		}
	}
	return LidUnknown
}

// subscribeLidEvents does nothing, OpenBSD has no lid event stream. Lid
// changes are seen on the next poll.
func subscribeLidEvents(_ chan<- struct{}, _ <-chan struct{}) {}
