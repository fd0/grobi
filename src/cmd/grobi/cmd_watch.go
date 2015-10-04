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

	ticker := time.NewTicker(2 * time.Second)

	var lastOutputs Outputs
	for {
		select {
		case <-ch:
			verbosePrintf("new output change event from i3 received\n")
		case <-ticker.C:
			verbosePrintf("regularly checking xrandr\n")
		}

		newOutputs, err := GetOutputs()
		if err != nil {
			return err
		}

		if lastOutputs.Equals(newOutputs) {
			verbosePrintf("nothing has changed, continuing\n")
			continue
		}

		err = MatchRules(globalOpts.cfg.Rules, newOutputs)
		if err != nil {
			return err
		}

		lastOutputs = newOutputs
	}
}
