package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds all configuration for grobi.
type Config struct {
	DefaultAction string
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
func openConfigFile() (io.ReadCloser, error) {
	filename := filepath.Join(xdgConfigDir(), "grobi.conf")
	if f, err := os.Open(filename); err == nil {
		return f, nil
	}

	filename = filepath.Join(os.Getenv("HOME"), ".grobi.conf")
	if f, err := os.Open(filename); err == nil {
		return f, nil
	}

	return nil, errors.New("could not find config file")
}

// readConfig returns a configuration struct read from a configuration file.
func readConfig() (Config, error) {
	rd, err := openConfigFile()
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	_, err = toml.DecodeReader(rd, &cfg)
	if err != nil {
		return Config{}, err
	}

	err = rd.Close()
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
