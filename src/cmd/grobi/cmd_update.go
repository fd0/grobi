package main

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

	outputs, err := GetOutputs()
	if err != nil {
		return err
	}

	for _, rule := range globalOpts.cfg.Rules {
		if rule.Match(outputs) {
			verbosePrintf("found matching rule (name %v)\n", rule.Name)
			return ApplyRule(outputs, rule)
		}
	}

	verbosePrintf("no rules match the current configuration\n")
	return nil
}
