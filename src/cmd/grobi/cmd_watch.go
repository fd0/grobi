package main

import "fmt"

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

	fmt.Printf("watch\n")

	return nil
}
