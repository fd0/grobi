package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Output encapsulates a physical output with detected modes.
type Output struct {
	Name      string
	Modes     Modes
	Connected bool
}

func (o Output) String() string {
	var con string
	if o.Connected {
		con = " (connected)"
	}
	str := fmt.Sprintf("%s%s", o.Name, con)
	if len(o.Modes) > 0 {
		str += fmt.Sprintf(" %v", o.Modes)
	}
	return str
}

// Outputs is a list of outputs.
type Outputs []Output

// Present returns true iff the list of outputs contains the named output.
func (os Outputs) Present(name string) bool {
	for _, o := range os {
		m, err := path.Match(name, o.Name)
		if err != nil {
			return false
		}

		if m {
			return true
		}
	}
	return false
}

// Connected returns true iff the list of outputs contains the named output and
// it is connected.
func (os Outputs) Connected(name string) bool {
	for _, o := range os {
		m, err := path.Match(name, o.Name)
		if err != nil {
			return false
		}

		if m && o.Connected {
			return true
		}
	}
	return false
}

// Mode is an output mode that may be active or default.
type Mode struct {
	Name    string
	Default bool
	Active  bool
}

func (m Mode) String() string {
	var suffix string

	if m.Active {
		suffix += "*"
	}

	if m.Default {
		suffix += "+"
	}

	return m.Name + suffix
}

// Modes is a list of Mode.
type Modes []Mode

func (m Modes) String() string {
	var str []string
	for _, mode := range m {
		str = append(str, mode.String())
	}
	return strings.Join(str, " ")
}

// errNotModeLine is returned by parseModeLine when the line doesn't match
// the format for a mode line.
var errNotModeLine = errors.New("not a mode line")

// parseOutputLine returns the output parsed from the string.
func parseOutputLine(line string) (Output, error) {
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
		output.Connected = true
	case "disconnected":
		output.Connected = false
	default:
		return Output{}, fmt.Errorf("unknown state %q", ws.Text())
	}

	return output, nil
}

// parseModeLine returns the mode parsed from the string.
func parseModeLine(line string) (mode Mode, err error) {
	if !strings.HasPrefix(line, "  ") {
		return Mode{}, errNotModeLine
	}

	ws := bufio.NewScanner(bytes.NewReader([]byte(line)))
	ws.Split(bufio.ScanWords)

	if !ws.Scan() {
		return Mode{}, fmt.Errorf("line too short, mode name not found: %s", line)
	}
	mode.Name = ws.Text()

	if !ws.Scan() {
		return Mode{}, fmt.Errorf("line too short, no refresh rate found: %s", line)
	}
	rate := ws.Text()

	if rate[len(rate)-1] == '+' {
		mode.Default = true
	}

	if rate[len(rate)-2] == '*' {
		mode.Active = true
	}

	// handle single-word "+", which happens when a mode is default but not active
	if ws.Scan() && ws.Text() == "+" {
		mode.Default = true
	}

	return mode, nil
}

// RandrParse returns the list of outputs parsed from the reader.
func RandrParse(rd io.Reader) (outputs Outputs, err error) {
	ls := bufio.NewScanner(rd)

	const (
		StateStart = iota
		StateOutput
		StateMode
	)

	var (
		state  = StateStart
		output Output
	)

nextLine:
	for ls.Scan() {
		line := ls.Text()

		for {
			switch state {
			case StateStart:
				if strings.HasPrefix(line, "Screen ") {
					state = StateOutput
					continue nextLine
				}
				return nil, fmt.Errorf(`first line should start with "Screen", found: %v`, line)

			case StateOutput:
				output, err = parseOutputLine(line)
				if err != nil {
					return nil, err
				}
				state = StateMode
				continue nextLine

			case StateMode:
				mode, err := parseModeLine(line)
				if err == errNotModeLine {
					outputs = append(outputs, output)
					output = Output{}
					state = StateOutput
					continue
				}

				if err != nil {
					return nil, err
				}

				output.Modes = append(output.Modes, mode)
				continue nextLine
			}
		}
	}

	if output.Name != "" {
		outputs = append(outputs, output)
	}

	return outputs, nil
}

// GetOutputs runs `xrandr` and returns the parsed output.
func GetOutputs() (Outputs, error) {
	cmd := exec.Command("xrandr", "--query")
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return RandrParse(bytes.NewReader(output))
}
