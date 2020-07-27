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

func MatchRules(rules []Rule, outputs Outputs) (Rule, error) {
	for _, rule := range rules {
		if rule.Match(outputs) {
			return rule, nil
		}
	}

	return Rule{}, nil
}

func (cmd CmdUpdate) Execute(args []string) (err error) {
	err = globalOpts.ReadConfigfile()
	if err != nil {
		return err
	}

	// install panic handler if commands are given
	defer RunCommandsOnFailure(&err, globalOpts.cfg.OnFailure)()

	outputs, err := DetectOutputs()
	if err != nil {
		return err
	}

	rule, err := MatchRules(globalOpts.cfg.Rules, outputs)
	if err != nil {
		return err
	}

	V("rule %q matches\n", rule.Name)

	return ApplyRule(outputs, rule)
}
