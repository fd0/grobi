package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	cfg, err := readConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "errror reading config file: %v", err)
		os.Exit(1)
	}

	spew.Dump(cfg)

	// ch, err := i3ipc.Subscribe(i3ipc.I3OutputEvent)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "unable to connect to i3: %v", err)
	// 	os.Exit(1)
	// }

	// for range ch {
	// 	fmt.Printf("received output change event\n")
	// 	cmd := exec.Command("sh", "-c", cfg.DefaultAction)
	// 	cmd.Stderr = os.Stderr
	// 	err = cmd.Run()
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "error running command %q\n", cfg.DefaultAction)
	// 	}
	// }

	outputs, err := GetOutputs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running xrandr: %v\n", err)
	}

	spew.Dump(outputs)
}
