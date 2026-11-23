//go:build linux || freebsd

package main

import (
	"net"
	"path/filepath"
	"testing"
	"time"
)

func TestWatchLidSocket(t *testing.T) {
	path := filepath.Join(t.TempDir(), "acpid.socket")
	ln, err := net.Listen("unix", path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()

	ch := make(chan struct{}, 1)
	done := make(chan struct{})
	defer close(done)
	go watchLidSocket("unix", []string{filepath.Join(t.TempDir(), "missing"), path}, parseTestEvent, ch, done)

	conn, err := ln.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	_, err = conn.Write([]byte("other\nclosed\n"))
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-ch:
	case <-time.After(2 * time.Second):
		t.Fatal("no wake-up received")
	}

	if l := LidState(lastLidEvent.Load()); l != LidClosed {
		t.Errorf("wanted closed, got %v", l)
	}
	lastLidEvent.Store(int32(LidUnknown))
}

func parseTestEvent(line string) LidState {
	if line == "closed" {
		return LidClosed
	}
	return LidUnknown
}
