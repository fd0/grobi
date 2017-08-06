package main

import (
	"fmt"
	"runtime"
)

var version = "compiled manually"
var compiledAt = "unknown time"

type CmdVersion struct{}

func init() {
	_, err := parser.AddCommand("version",
		"display version",
		"The version command displays detailed information about the version",
		&CmdVersion{})
	if err != nil {
		panic(err)
	}
}

func (cmd CmdVersion) Execute(args []string) error {
	fmt.Printf("grobi %s\ncompiled with %v on %v\n",
		version, runtime.Version(), runtime.GOOS)

	return nil
}
