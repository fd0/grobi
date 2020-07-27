package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config holds all configuration for grobi.
type Config struct {
	Rules []Rule

	ExecuteAfter []string `yaml:"execute_after"`
	OnFailure    []string `yaml:"on_failure"`
}

// xdgConfigDir returns the config directory according to the xdg standard, see
// http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html.
func xdgConfigDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir
	}

	return filepath.Join(os.Getenv("HOME"), ".config")
}

// openConfigFile returns a reader for the config file.
func openConfigFile(name string) (io.ReadCloser, error) {
	for _, filename := range []string{
		name,
		os.Getenv("GROBI_CONFIG"),
		filepath.Join(xdgConfigDir(), "grobi.conf"),
		filepath.Join(os.Getenv("HOME"), ".grobi.conf"),
		"/etc/xdg/grobi.conf"} {
		if filename != "" {
			if f, err := os.Open(filename); err == nil {
				V("reading config from %v\n", filename)
				return f, nil
			}
		}
	}

	return nil, errors.New("could not find config file")
}

// readConfig returns a configuration struct read from a configuration file.
func readConfig(name string) (Config, error) {
	rd, err := openConfigFile(name)
	if err != nil {
		return Config{}, err
	}

	buf, err := ioutil.ReadAll(rd)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return Config{}, err
	}

	err = rd.Close()
	if err != nil {
		return Config{}, err
	}

	if err = cfg.Valid(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Valid returns an error if the config is invalid, ie a pattern is malformed.
func (cfg Config) Valid() error {

	for _, rule := range cfg.Rules {
		for _, list := range [][]string{rule.OutputsPresent, rule.OutputsAbsent, rule.OutputsConnected, rule.OutputsDisconnected} {
			for _, pat := range list {
				if _, err := path.Match(pat, ""); err != nil {
					return fmt.Errorf("pattern %q malformed: %v", pat, err)
				}
			}
		}
	}

	return nil
}
