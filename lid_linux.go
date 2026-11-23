package main

import (
	"os"
	"path/filepath"
	"strings"
)

// acpidSockets are tried in order: systemd socket activation uses /run, acpid
// itself defaults to /var/run. --acpid-socket replaces both.
var acpidSockets = []string{"/run/acpid.socket", "/var/run/acpid.socket"}

// Fields of an acpid event line, e.g. "button/lid LID close".
const (
	acpidFieldClass  = 0 // "button/lid"
	acpidFieldDevice = 1 // "LID"
	acpidFieldAction = 2 // "open" or "close"
	acpidMinFields   = acpidFieldAction + 1

	acpidClassLid = "button/lid"
)

func readLidState() LidState {
	files, _ := filepath.Glob("/proc/acpi/button/lid/*/state")
	for _, f := range files {
		buf, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		// "state:      open"
		if l := parseLidWord(strings.TrimPrefix(strings.TrimSpace(string(buf)), "state:")); l != LidUnknown {
			return l
		}
	}
	return LidUnknown
}

func parseLidWord(s string) LidState {
	switch strings.TrimSpace(s) {
	case "open":
		return LidOpen
	case "close", "closed":
		return LidClosed
	}
	return LidUnknown
}

// parseAcpidEvent parses lines like "button/lid LID close".
func parseAcpidEvent(line string) LidState {
	f := strings.Fields(line)
	if len(f) < acpidMinFields || f[acpidFieldClass] != acpidClassLid {
		return LidUnknown
	}
	return parseLidWord(f[acpidFieldAction])
}

// subscribeLidEvents reads lid events from acpid, if it runs.
func subscribeLidEvents(ch chan<- struct{}, done <-chan struct{}) {
	paths := acpidSockets
	if globalOpts.AcpidSocket != "" {
		paths = []string{globalOpts.AcpidSocket}
	}
	watchLidSocket("unix", paths, parseAcpidEvent, ch, done)
}
