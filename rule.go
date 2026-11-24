package main

// Rule is a rule to configure outputs.
type Rule struct {
	Name string

	OutputsConnected    []string `yaml:"outputs_connected"`
	OutputsDisconnected []string `yaml:"outputs_disconnected"`
	OutputsPresent      []string `yaml:"outputs_present"`
	OutputsAbsent       []string `yaml:"outputs_absent"`

	// Lid is "open", "closed" or empty (lid ignored).
	Lid string `yaml:"lid"`

	// InputsPresent and InputsAbsent match XInput device names, see "grobi show".
	InputsPresent []string `yaml:"inputs_present"`
	InputsAbsent  []string `yaml:"inputs_absent"`

	ConfigureRow     []string `yaml:"configure_row"`
	ConfigureColumn  []string `yaml:"configure_column"`
	ConfigureSingle  string   `yaml:"configure_single"`
	ConfigureCommand string   `yaml:"configure_command"`

	Primary string `yaml:"primary"`

	DisableOrder []string `yaml:"disable_order"`

	Atomic bool `yaml:"atomic"`

	ExecuteAfter []string `yaml:"execute_after"`
}

// Match returns true iff the rule matches the given state. A rule with a lid
// condition never matches while the lid state is unknown, a rule with an
// input condition never matches while the input devices are unknown.
func (r Rule) Match(state State) bool {
	if r.Lid != "" && r.Lid != state.Lid.String() {
		return false
	}

	if len(r.InputsPresent) > 0 || len(r.InputsAbsent) > 0 {
		if state.Inputs == nil {
			return false
		}
		for _, pat := range r.InputsPresent {
			if !state.Inputs.Present(pat) {
				return false
			}
		}
		for _, pat := range r.InputsAbsent {
			if state.Inputs.Present(pat) {
				return false
			}
		}
	}

	outputs := state.Outputs
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
