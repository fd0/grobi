package main

import "fmt"

type CmdRules struct{}

func init() {
	_, err := parser.AddCommand("rules",
		"list rules",
		"The rules command lists the configured rules",
		&CmdRules{})
	if err != nil {
		panic(err)
	}
}

func printList(label string, args []string) {
	if len(args) > 0 {
		fmt.Printf("  %s: %v\n", label, args)
	}
}

func printOne(label string, arg string) {
	if len(arg) > 0 {
		fmt.Printf("  %s: %v\n", label, arg)
	}
}

func (cmd CmdRules) Execute(args []string) error {
	globalOpts.ReadConfigfile()

	for _, rule := range globalOpts.cfg.Rules {
		fmt.Printf("%v\n", rule.Name)

		if globalOpts.Verbose {
			printList("Connected", rule.OutputsConnected)
			printList("Disconnected", rule.OutputsDisconnected)
			printList("Present", rule.OutputsPresent)
			printList("Absent", rule.OutputsAbsent)
			printList("ConfigureRow", rule.ConfigureRow)
			printOne("ConfigureSingle", rule.ConfigureSingle)
			printOne("ConfigureCommand", rule.ConfigureCommand)
			printList("ExecuteAfter", rule.ExecuteAfter)
		}
	}

	return nil
}
