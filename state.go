package main

import (
	"fmt"
	"sync/atomic"
)

// LidState is the state of the laptop lid.
type LidState int

const (
	// LidUnknown means there is no lid or its state can't be read.
	LidUnknown LidState = iota
	LidOpen
	LidClosed
)

func (l LidState) String() string {
	switch l {
	case LidOpen:
		return "open"
	case LidClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// State is the machine state rules are matched against.
type State struct {
	Outputs Outputs
	Lid     LidState
}

// CurrentState combines outputs with the current lid state.
func CurrentState(outputs Outputs) State {
	return State{Outputs: outputs, Lid: ReadLidState()}
}

func validLid(s string) error {
	switch s {
	case "", "open", "closed":
		return nil
	}
	return fmt.Errorf("lid %q invalid, use \"open\" or \"closed\"", s)
}

// lastLidEvent holds the lid state from the most recent lid event. It is the
// fallback when the platform can't read the lid state directly.
var lastLidEvent atomic.Int32

// ReadLidState returns the current lid state.
func ReadLidState() LidState {
	if l := readLidState(); l != LidUnknown {
		return l
	}
	return LidState(lastLidEvent.Load())
}

// lidWakeupBuffer is the capacity of the lid wake-up channel: one pending
// wake-up is enough, more are dropped.
const lidWakeupBuffer = 1
