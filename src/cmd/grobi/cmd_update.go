package main

import "fmt"

type CmdUpdate struct{}

func init() {
	_, err := parser.AddCommand("update",
		"update outputs",
		"The update command updates the outputs as configured",
		&CmdUpdate{})
	if err != nil {
		panic(err)
	}
}

func (cmd CmdUpdate) Execute(args []string) error {
	globalOpts.ReadConfigfile()

	fmt.Printf("update\n")

	return nil
}
