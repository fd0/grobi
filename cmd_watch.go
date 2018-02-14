package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
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

	V("grobi %s, compiled with %v on %v\n", version, runtime.Version(), runtime.GOOS)
	V("successfully subscribed to X RANDR change events\n")

	var tickerCh <-chan time.Time
	if globalOpts.PollInterval > 0 {
		tickerCh = time.NewTicker(time.Duration(globalOpts.PollInterval) * time.Second).C
	}

	var backoffCh <-chan time.Time
	var disablePoll bool
	var eventReceived bool

	var lastRule Rule
	var lastOutputs Outputs
	for {
		if !disablePoll {
			var outputs Outputs
			var err error

			if eventReceived {
				outputs, err = DetectOutputs()
				eventReceived = false
			} else {
				outputs, err = GetOutputs()
			}

			if err != nil {
				return err
			}

			V("got outputs: %v", outputs)

			// disable outputs which have a changed display
			var off Outputs
			for _, o := range outputs {
				for _, last := range lastOutputs {
					if o.Name != last.Name {
						continue
					}

					if o.Active() != last.Active() {
						V("  output %v: monitor active has changed, disabling", o.Name)
						off = append(off, o)
						continue
					}

					if o.MonitorId != last.MonitorId {
						V("  output %v: monitor has changed, disabling", o.Name)
						off = append(off, o)
						continue
					}
				}
			}

			cmd, err := DisableOutputs(off)
			if err != nil {
				return err
			}

			if cmd != nil {
				// forget the last rule set, something has changed for sure
				lastRule = Rule{}

				err = RunCommand(cmd)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error disabling outputs: %v", err)
				}
			}

			rule, err := MatchRules(globalOpts.cfg.Rules, outputs)
			if err != nil {
				return err
			}

			if rule.Name != lastRule.Name {
				V("outputs: %v", outputs)
				V("new rule found: %v", rule.Name)

				err = ApplyRule(outputs, rule)
				if err != nil {
					return err
				}

				lastRule = rule

				if globalOpts.Pause > 0 {
					V("disable polling for %d seconds\n", globalOpts.Pause)
					disablePoll = true
					backoffCh = time.After(time.Duration(globalOpts.Pause) * time.Second)
				}
			}

			lastOutputs = outputs
		}

		select {
		case ev := <-ch:
			V("new RANDR change event received\n")
			if ev.Error != nil {
				return ev.Error
			}

			eventReceived = true
		case <-tickerCh:
		case <-backoffCh:
			V("reenable polling\n")
			backoffCh = nil
			disablePoll = false
		}
	}
}
