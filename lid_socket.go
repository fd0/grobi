//go:build linux || freebsd

package main

import (
	"bufio"
	"net"
	"time"
)

const lidReconnectDelay = 10 * time.Second

// watchLidSocket reads lid events from a unix socket and sends a wake-up into
// ch for each one. parse returns LidUnknown for lines which aren't lid events.
// A missing or closed socket is retried until done is closed, polling keeps
// working in the meantime.
func watchLidSocket(network, path string, parse func(string) LidState, ch chan<- struct{}, done <-chan struct{}) {
	connected := true // log the first failure
	for {
		conn, err := net.Dial(network, path)
		if err != nil {
			if connected {
				V("lid events: %v, retrying every %v\n", err, lidReconnectDelay)
			}
			connected = false
		} else {
			V("lid events: connected to %v\n", path)
			connected = true
			readLidEvents(conn, parse, ch, done)
			V("lid events: connection to %v lost\n", path)
		}

		select {
		case <-done:
			return
		case <-time.After(lidReconnectDelay):
		}
	}
}

func readLidEvents(conn net.Conn, parse func(string) LidState, ch chan<- struct{}, done <-chan struct{}) {
	stop := make(chan struct{})
	defer close(stop)
	go func() {
		select {
		case <-done:
		case <-stop:
		}
		_ = conn.Close()
	}()

	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		l := parse(sc.Text())
		if l == LidUnknown {
			continue
		}

		V("lid event: %v\n", l)
		lastLidEvent.Store(int32(l))

		// wake up the watch loop, one pending wake-up is enough
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
