package main

import (
	"log"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

type CmdWatch struct{}

func init() {
	_, err := parser.AddCommand("watch",
		"watch for changes",
		"The watch command listens for changes and configures the outputs accordingly",
		&CmdWatch{})
	if err != nil {
		panic(err)
	}
}

type Event struct {
	Event xgb.Event
	Error error
}

const eventSendTimeout = 500 * time.Millisecond

func subscribeXEvents(ch chan<- Event, done <-chan struct{}) {
	X, err := xgb.NewConn()
	if err != nil {
		ch <- Event{Error: err}
		return
	}

	defer X.Close()
	if err = randr.Init(X); err != nil {
		ch <- Event{Error: err}
		return
	}

	root := xproto.Setup(X).DefaultScreen(X).Root

	eventMask := randr.NotifyMaskScreenChange |
		randr.NotifyMaskCrtcChange |
		randr.NotifyMaskOutputChange |
		randr.NotifyMaskOutputProperty

	err = randr.SelectInputChecked(X, root, uint16(eventMask)).Check()
	if err != nil {
		ch <- Event{Error: err}
		return
	}

	for {
		ev, err := X.WaitForEvent()
		select {
		case ch <- Event{Event: ev, Error: err}:
		case <-time.After(eventSendTimeout):
			continue
		case <-done:
			return
		}
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (cmd CmdWatch) Execute(args []string) error {
	globalOpts.ReadConfigfile()

	done := make(chan struct{})
	defer close(done)

	ch := make(chan Event)
	go subscribeXEvents(ch, done)

	V("successfully subscribed to X RANDR change events\n")

	var tickerCh <-chan time.Time
	if globalOpts.PollInterval > 0 {
		tickerCh = time.NewTicker(time.Duration(globalOpts.PollInterval) * time.Second).C
	}

	var backoffCh <-chan time.Time
	var disablePoll bool
	var eventReceived bool

	var lastOutputs Outputs
	for {
		if !disablePoll {
			var newOutputs Outputs
			var err error

			if eventReceived {
				newOutputs, err = DetectOutputs()
				eventReceived = false
			} else {
				newOutputs, err = GetOutputs()
			}

			if err != nil {
				return err
			}

			if !lastOutputs.Equals(newOutputs) {
				err = MatchRules(globalOpts.cfg.Rules, newOutputs)
				if err != nil {
					return err
				}

				lastOutputs = newOutputs

				if globalOpts.Pause > 0 {
					V("disable polling for %d seconds\n", globalOpts.Pause)
					disablePoll = true
					backoffCh = time.After(time.Duration(globalOpts.Pause) * time.Second)
				}
			}
		}

		select {
		case ev := <-ch:
			V("new RANDR change event received:\n")
			V("  %v\n", ev)
			if ev.Error != nil {
				return ev.Error
			}

			eventReceived = true
		case <-tickerCh:
			V("regularly checking xrandr\n")
		case <-backoffCh:
			V("reenable polling\n")
			backoffCh = nil
			disablePoll = false
		}
	}
}
