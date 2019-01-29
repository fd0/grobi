package main

import "fmt"

type CmdLayouts struct{}

func init() {
	_, err := parser.AddCommand("layouts",
		"list layouts",
		"The layouts command lists the configured layouts",
		&CmdLayouts{})
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

func (cmd CmdLayouts) Execute(args []string) error {
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
