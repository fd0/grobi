package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// GlobalOptions contains all global options.
type GlobalOptions struct {
	Verbose bool   `short:"v" long:"verbose"     default:"false" description:"Be verbose"`
	Config  string `short:"C" long:"config"                      description:"Read config from this file"`

	cfg *Config
}

func (gopts *GlobalOptions) ReadConfigfile() {
	if gopts.cfg != nil {
		return
	}

	cfg, err := readConfig(gopts.Config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config file: %v\n", err)
		os.Exit(1)
	}

	gopts.cfg = &cfg
}

var globalOpts = GlobalOptions{}
var parser = flags.NewParser(&globalOpts, flags.Default)

func verbosePrintf(format string, args ...interface{}) {
	if !globalOpts.Verbose {
		return
	}

	fmt.Printf(format, args...)
}

func main() {
	_, err := parser.Parse()
	if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
		os.Exit(0)
	}

	if err != nil {
		os.Exit(1)
	}
}
