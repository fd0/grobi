package main

import "testing"

func TestParseOpenBSDSensors(t *testing.T) {
	in := "hw.sensors.acpibtn0.indicator0=Off (power supply)\nhw.sensors.acpibtn1.indicator0=On (lid open)\n"
	if got := parseOpenBSDSensors(in); got != LidOpen {
		t.Errorf("wanted open, got %v", got)
	}
	in = "hw.sensors.acpibtn1.indicator0=Off (lid open)\n"
	if got := parseOpenBSDSensors(in); got != LidClosed {
		t.Errorf("wanted closed, got %v", got)
	}
	if got := parseOpenBSDSensors("hw.sensors.cpu0.temp0=45.00 degC\n"); got != LidUnknown {
		t.Errorf("wanted unknown, got %v", got)
	}
}
