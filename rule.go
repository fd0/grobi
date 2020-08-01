package main

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
