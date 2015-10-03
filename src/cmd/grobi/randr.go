package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Output struct {
	Name       string
	Modes      []Mode
	Active     bool
}

type Mode struct {
	Name  string
	Rate string
	Default bool
	Active bool
}

func ParseOutputLine(line string) (Output, error) {
	output := Output{}

	ws := bufio.NewScanner(bytes.NewReader([]byte(line)))
	ws.Split(bufio.ScanWords)

	if !ws.Scan() {
		return Output{}, fmt.Errorf("line too short, name not found: %s", line)
	}
	output.Name = ws.Text()

	if !ws.Scan() {
		return Output{}, fmt.Errorf("line too short, state not found: %s", line)
	}

	switch ws.Text() {
	case "connected":
		output.Active = true
	case "disconnected":
		output.Active = false
	default:
		return Output{}, fmt.Errorf("unknown state %q", ws.Text())
	}

	return output, nil
}

// func ParseResolutionLine(line) Mode

func RandrParse(rd io.Reader) (outputs []Output, err error) {
	ls := bufio.NewScanner(rd)

	for ls.Scan() {
		line := ls.Text()
		fmt.Printf("line: %v\n", line)

		if strings.HasPrefix(line, "Screen ") {
			continue
		}

		output, err := ParseOutputLine(line)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, output)
	}

	return nil, nil
}
