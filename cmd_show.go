package main

import "fmt"

type CmdShow struct{}

func init() {
	_, err := parser.AddCommand("show",
		"show monitors and IDs",
		"The show command lists all connected monitors with their IDs",
		&CmdShow{})
	if err != nil {
		panic(err)
	}
}

func ListOutput(output Output) {
	fmt.Printf("%- 10s %s\n", output.Name, output.MonitorId)
}

func (cmd CmdShow) Execute(args []string) error {
	outputs, err := DetectOutputs()
	if err != nil {
		return err
	}
	for _, output := range outputs {
		if output.Connected {
			ListOutput(output)
		}
	}
	return nil
}
