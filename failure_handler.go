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
			return
		}

		for _, cmd := range commands {
			err := RunCommand(exec.Command("sh", "-c", cmd))
			if err != nil {
				fmt.Fprintf(os.Stderr, "running command %q failed: %v\n", cmd, err)
			}
		}
	}
}
