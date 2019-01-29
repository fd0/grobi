package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type CmdApply struct{}

func init() {
	_, err := parser.AddCommand("apply",
		"apply a rule",
		"The apply command configures the outputs as described in the given",
		&CmdApply{})
	if err != nil {
		panic(err)
	}
}

func (cmd CmdApply) Usage() string {
	return "RULE"
}

func ApplyRule(outputs Outputs, rule Rule) error {
	var cmds []*exec.Cmd
	var err error

	switch {
	case rule.ConfigureSingle != "" || len(rule.ConfigureRow) > 0:
		cmds, err = BuildCommandOutputRow(rule, outputs)
	case rule.ConfigureCommand != "":
		cmds = []*exec.Cmd{exec.Command("sh", "-c", rule.ConfigureCommand)}
	default:
		return fmt.Errorf("no output configuration for rule %v", rule.Name)
	}

	after := append(globalOpts.cfg.ExecuteAfter, rule.ExecuteAfter...)
	if len(after) > 0 {
		for _, cmd := range after {
			cmds = append(cmds, exec.Command("sh", "-c", cmd))
		}
	}

	if err != nil {
		return err
	}
	for _, cmd := range cmds {
		err = RunCommand(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "executing command for rule %v failed: %v\n", rule.Name, err)
		}
	}

	return nil
}

func (cmd CmdApply) Execute(args []string) error {
	globalOpts.ReadConfigfile()

	if len(args) != 1 {
		return errors.New("need exactly one rule name as the parameter")
	}

	outputs, err := DetectOutputs()
	if err != nil {
		return err
	}

	ruleName := strings.ToLower(args[0])
	for _, rule := range globalOpts.cfg.Rules {
		if strings.ToLower(rule.Name) == ruleName {
			V("found matching rule (name %v)\n", rule.Name)
			return ApplyRule(outputs, rule)
		}
	}

	return fmt.Errorf("rule %q not found", ruleName)
}
