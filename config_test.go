package main

import (
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestReadConfigSample(t *testing.T) {
	cfg, err := readConfig("doc/grobi.conf")
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Rules) != 6 {
		t.Fatalf("want 6 rules, got %d", len(cfg.Rules))
	}

	r := cfg.Rules[0]
	if r.Name != "Docking Station" || !r.Atomic || r.Primary != "HDMI3" {
		t.Errorf("unexpected first rule: %+v", r)
	}

	if len(cfg.ExecuteAfter) == 0 || len(cfg.OnFailure) == 0 {
		t.Errorf("execute_after/on_failure not parsed: %+v", cfg)
	}
}

// YAML 1.1 style booleans were accepted by yaml.v2, make sure they still are.
func TestRuleAtomicYAML11Bools(t *testing.T) {
	for _, v := range []string{"true", "yes", "on", "True", "YES"} {
		var r Rule
		if err := yaml.Unmarshal([]byte("atomic: "+v), &r); err != nil {
			t.Errorf("atomic: %s: %v", v, err)
			continue
		}
		if !r.Atomic {
			t.Errorf("atomic: %s: parsed as false", v)
		}
	}
}
