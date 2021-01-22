package main

import (
	"fmt"
	"os"
	"os/exec"
)

func RunCommandsOnFailure(err *error, commands []string) func() {
	return func() {
		r := recover()
		if r != nil {
			fmt.Fprintf(os.Stderr, "recovered error: %v\n", r)
		}

		if *err == nil && r == nil {
			fmt.Fprintf(os.Stderr, "no error found, exiting\n")
			return
		}

		fmt.Fprintf(os.Stderr, "encountered error: %v\n", *err)

		for _, cmd := range commands {
			fmt.Fprintf(os.Stderr, "running on_failure command: %v\n", cmd)
			err := RunCommand(exec.Command("sh", "-c", cmd))
			if err != nil {
				fmt.Fprintf(os.Stderr, "command failed: %v\n", err)
			}
		}
	}
}
