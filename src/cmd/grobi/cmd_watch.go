package main

import (
	"fmt"

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

func (cmd CmdWatch) Execute(args []string) error {
	globalOpts.ReadConfigfile()

	ch, err := i3ipc.Subscribe(i3ipc.I3OutputEvent)
	if err != nil {
		return fmt.Errorf("unable to connect to i3: %v", err)
	}

	verbosePrintf("successfully subscribed to the i3 IPC socket\n")

nextEvent:
	for range ch {
		fmt.Printf("received output change event\n")

		outputs, err := GetOutputs()
		if err != nil {
			return err
		}

		for _, rule := range globalOpts.cfg.Rules {
			if rule.Match(outputs) {
				verbosePrintf("found matching rule (name %v)\n", rule.Name)
				if err = ApplyRule(outputs, rule); err != nil {
					return err
				}
				continue nextEvent
			}
		}
	}

	return nil
}
