package main

import "strings"

// Rule is a rule to configure outputs.
type Rule struct {
	Name string

	OutputsConnected    []string `yaml:"outputs_connected"`
	OutputsDisconnected []string `yaml:"outputs_disconnected"`
	OutputsPresent      []string `yaml:"outputs_present"`
	OutputsAbsent       []string `yaml:"outputs_absent"`

	ConfigureRow     []string `yaml:"configure_row"`
	ConfigureColumn  []string `yaml:"configure_column"`
	ConfigureSingle  string   `yaml:"configure_single"`
	ConfigureCommand string   `yaml:"configure_command"`

	Primary string `yaml:"primary"`

	DisableOrder []string `yaml:"disable_order"`

	Atomic bool `yaml:"atomic"`

	ExecuteAfter []string `yaml:"execute_after"`
}

func (r Rule) OutputsDiff(old Rule) Outputs {
	outputs := []string{}
	if r.ConfigureSingle != "" {
		outputs = append(outputs, r.ConfigureSingle)
	}
	outputs = append(outputs, r.ConfigureRow...)
	outputs = append(outputs, r.ConfigureColumn...)
	for k, v := range outputs {
		outputs[k] = strings.SplitN(v, "@", 2)[0]
	}

	outputsOld := []string{}
	if old.ConfigureSingle != "" {
		outputsOld = append(outputsOld, old.ConfigureSingle)
	}
	outputsOld = append(outputsOld, old.ConfigureRow...)
	outputsOld = append(outputsOld, old.ConfigureColumn...)
	for k, v := range outputsOld {
		outputsOld[k] = strings.SplitN(v, "@", 2)[0]
	}

	var diff Outputs
	for _, old := range outputsOld {
		foundInOutputs := false
		for _, v := range outputs {
			if old == v {
				foundInOutputs = true
			}
		}

		if !foundInOutputs {
			diff = append(diff, Output{Name: old})
		}
	}

	return diff
}

// Match returns true iff the rule matches for the given list of outputs.
func (r Rule) Match(outputs Outputs) bool {
	for _, name := range r.OutputsAbsent {
		if outputs.Present(name) {
			return false
		}
	}

	for _, name := range r.OutputsDisconnected {
		if outputs.Connected(name) {
			return false
		}
	}

	for _, name := range r.OutputsPresent {
		if !outputs.Present(name) {
			return false
		}
	}

	for _, name := range r.OutputsConnected {
		if !outputs.Connected(name) {
			return false
		}
	}

	return true
}
