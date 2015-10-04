package main

import (
	"fmt"
	"time"

	"github.com/proxypoke/i3ipc"
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

func MatchRules(rules []Rule, outputs Outputs) error {
	for _, rule := range rules {
		if rule.Match(outputs) {
			verbosePrintf("found matching rule (name %v)\n", rule.Name)
			if err := ApplyRule(outputs, rule); err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func (cmd CmdWatch) Execute(args []string) error {
	globalOpts.ReadConfigfile()

	ch, err := i3ipc.Subscribe(i3ipc.I3OutputEvent)
	if err != nil {
		return fmt.Errorf("unable to connect to i3: %v", err)
	}

	verbosePrintf("successfully subscribed to the i3 IPC socket\n")

	var tickerCh <-chan time.Time
	if globalOpts.PollInterval > 0 {
		tickerCh = time.NewTicker(time.Duration(globalOpts.PollInterval) * time.Second).C
	}

	var backoffCh <-chan time.Time
	var disablePoll bool

	var lastOutputs Outputs
	for {
		if !disablePoll {
			newOutputs, err := GetOutputs()
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
					verbosePrintf("disable polling for %d seconds\n", globalOpts.Pause)
					disablePoll = true
					backoffCh = time.After(time.Duration(globalOpts.Pause) * time.Second)
				}
			}
		}

		select {
		case <-ch:
			verbosePrintf("new output change event from i3 received\n")
		case <-tickerCh:
			verbosePrintf("regularly checking xrandr\n")
		case <- backoffCh:
			verbosePrintf("reenable polling\n")
			backoffCh = nil
			disablePoll = false
		}
	}
}
