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

	return MatchRules(globalOpts.cfg.Rules, outputs)
}
