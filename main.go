package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jessevdk/go-flags"
)

// GlobalOptions contains all global options.
type GlobalOptions struct {
	Verbose      bool   `short:"v" long:"verbose"                     description:"Be verbose"`
	Config       string `short:"C" long:"config"                      description:"Read config from this file"`
	DryRun       bool   `short:"n" long:"dry-run"                     description:"Only print what commands would be executed without actually runnig them"`
	PollInterval uint   `short:"i" long:"interval"    default:"2"     description:"Number of seconds between polls, set to zero to disable polling"`
	ActivePoll   bool   `short:"a" long:"active-poll"                 description:"Force xrandr to re-detect outputs during polling"`
	Pause        uint   `short:"p" long:"pause"       default:"0"     description:"Number of seconds to pause after a change was executed"`
	Logfile      string `short:"l" long:"logfile"                     description:"Write log to file"`

	cfg     *Config
	log     *log.Logger
	logfile *log.Logger
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

	V("running command %v %v\n", cmd.Path, strings.Join(cmd.Args, " "))
	cmd.Stderr = os.Stderr
	if globalOpts.Verbose {
		cmd.Stdout = os.Stdout
	}
	return cmd.Run()
}

var globalOpts = GlobalOptions{}
var parser = flags.NewParser(&globalOpts, flags.Default)

func V(s string, data ...interface{}) {
	if globalOpts.Verbose && globalOpts.log == nil {
		globalOpts.log = log.New(os.Stdout, "grobi: ", log.Lmicroseconds|log.Ltime)
	}

	if globalOpts.log != nil {
		globalOpts.log.Printf(s, data...)
	}

	if globalOpts.Logfile != "" && globalOpts.logfile == nil {
		f, err := os.OpenFile(globalOpts.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to open logfile: %v\n", err)
			os.Exit(23)
		}
		globalOpts.logfile = log.New(f, "", log.Lmicroseconds|log.Ltime)
	}

	if globalOpts.logfile != nil {
		globalOpts.logfile.Printf(s, data...)
	}
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
