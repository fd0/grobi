package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jessevdk/go-flags"
)

// GlobalOptions contains all global options.
type GlobalOptions struct {
	Verbose bool   `short:"v" long:"verbose"     default:"false" description:"Be verbose"`
	Config  string `short:"C" long:"config"                      description:"Read config from this file"`
	DryRun  bool   `short:"n" long:"dry-run"                     description:"Only print what commands would be executed without actually runnig them"`

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

// RunCommand runs the given command or prints the arguments to stdout if
// globalOpts.DryRun is true.
func RunCommand(cmd *exec.Cmd) error {
	if globalOpts.DryRun {
		s := fmt.Sprintf("%s", cmd.Args)
		fmt.Printf("%s\n", s[1:len(s)-1])
		return nil
	}

	verbosePrintf("running command %v %v\n", cmd.Path, strings.Join(cmd.Args, " "))
	cmd.Stderr = os.Stderr
	if globalOpts.Verbose {
		cmd.Stdout = os.Stdout
	}
	return cmd.Run()
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
