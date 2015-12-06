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

func (cmd CmdUpdate) Execute(args []string) error {
	globalOpts.ReadConfigfile()

	outputs, err := DetectOutputs()
	if err != nil {
		return err
	}

	return MatchRules(globalOpts.cfg.Rules, outputs)
}
