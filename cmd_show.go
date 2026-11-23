package main

import "fmt"

type CmdShow struct{}

func init() {
	_, err := parser.AddCommand("show",
		"show monitors, IDs and input devices",
		"The show command lists all connected monitors with their IDs and the input devices",
		&CmdShow{})
	if err != nil {
		panic(err)
	}
}

func ListOutput(output Output) {
	fmt.Printf("%- 10s %s\n", output.Name, output.MonitorID)
}

func (cmd CmdShow) Execute(_ []string) error {
	outputs, err := DetectOutputs()
	if err != nil {
		return err
	}
	for _, output := range outputs {
		if output.Connected {
			ListOutput(output)
		}
	}

	// input device names for inputs_present/inputs_absent
	seen := make(map[string]bool)
	for _, name := range ReadInputDevices() {
		if !seen[name] {
			seen[name] = true
			fmt.Printf("%- 10s %s\n", "input", name)
		}
	}
	return nil
}
