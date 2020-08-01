package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
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
	Primary   bool
	MonitorId string
}

func (o Output) String() string {
	var con string
	switch {
	case o.Connected && o.Primary:
		con = " (connected, primary)"
	case o.Connected:
		con = " (connected)"
	case o.Primary:
		con = " (primary)"
	}
	str := fmt.Sprintf("%s%s", o.Name, con)

	if len(o.Modes) > 0 {
		for _, m := range o.Modes {
			if m.Active || m.Default {
				str += fmt.Sprintf(" %v", m)
			}
		}
	}

	if o.MonitorId != "" {
		str += fmt.Sprintf(" [%v]", o.MonitorId)
	}
	return str
}

// Equals checks whether the two Outputs are equal.
func (o Output) Equals(other Output) bool {
	if o.Name != other.Name || o.Connected != other.Connected {
		return false
	}

	if len(o.Modes) != len(other.Modes) {
		return false
	}

	for i := range o.Modes {
		m1 := o.Modes[i]
		m2 := other.Modes[i]

		if m1 != m2 {
			return false
		}
	}

	if o.MonitorId != other.MonitorId {
		return false
	}

	return true
}

// Active returns true if an output has an active mode.
func (o Output) Active() bool {
	for _, mode := range o.Modes {
		if mode.Active {
			return true
		}
	}

	return false
}

// Outputs is a list of outputs.
type Outputs []Output

// Present returns true iff the list of outputs contains the named output.
func (os Outputs) Present(name string) bool {
	for _, o := range os {
		// Check legacy name
		m, err := path.Match(name, o.Name)
		if err != nil {
			return false
		}
		if m {
			return true
		}

		// Check extended name
		m, err = path.Match(name, o.Name+"-"+o.MonitorId)
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
		if !o.Connected {
			continue
		}

		// Check legacy name
		m, err := path.Match(name, o.Name)
		if err != nil {
			return false
		}
		if m {
			return true
		}

		// Check extended name
		m, err = path.Match(name, o.Name+"-"+o.MonitorId)
		if err != nil {
			return false
		}
		if m {
			return true
		}
	}
	return false
}

// Equals checks whether the two Outputs are equal.
func (os Outputs) Equals(other Outputs) bool {
	if len(os) != len(other) {
		return false
	}

	for i := range os {
		out1 := os[i]
		out2 := other[i]

		if !out1.Equals(out2) {
			return false
		}
	}

	return true
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

// Generates the monitor id from the edid
func GenerateMonitorId(s string) (string, error) {
	var errEdidCorrupted = errors.New("corrupt EDID: " + s)
	if len(s) < 32 || s[:16] != "00ffffffffffff00" {
		return "", errEdidCorrupted
	}

	edid, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}

	// we only parse EDID 1.3 and 1.4
	if edid[18] != 1 || (edid[19] < 3 || edid[19] > 4) {
		return "", fmt.Errorf("unknown EDID version %d.%d", edid[18], edid[19])
	}

	manuf := binary.BigEndian.Uint16(edid[8:10])

	// The first bit is resevered and needs to be zero
	if manuf&0x8000 != 0x0000 {
		return "", errEdidCorrupted
	}

	// Decode the manufacturer 'A' = 0b00001, 'B' = 0b00010, ..., 'Z' = 0b11010
	var manufacturer string
	mask := uint16(0x7C00) // 0b0111110000000000
	for i := uint(0); i <= 10; i += 5 {
		number := ((manuf & (mask >> i)) >> (10 - i))
		manufacturer += string(number + 'A' - 1)
	}

	// Decode the product and serial number
	product := binary.LittleEndian.Uint16(edid[10:12])
	serial := binary.LittleEndian.Uint32(edid[12:16])

	// Decode four descriptor blocks
	var displayName, displaySerialNumber string
	for i := 0; i < 4; i++ {
		d := edid[54+i*18 : 54+18+i*18]

		// interesting descriptors start with three zeroes
		if d[0] != 0 || d[1] != 0 || d[2] != 0 {
			continue
		}

		switch d[3] {
		case 0xff: // display serial number
			displaySerialNumber = strings.TrimSpace(string(d[5:]))
		case 0xfc: // display name
			displayName = strings.TrimSpace(string(d[5:]))
		}
	}

	str := fmt.Sprintf("%s-%d-%d-%v-%v", manufacturer, product, serial, displayName, displaySerialNumber)
	return str, nil
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

	if !ws.Scan() {
		return output, nil
	}

	if ws.Text() == "primary" {
		output.Primary = true
		ws.Scan()
	}

	if output.Connected {
		return output, nil
	}

	// handle special case when output is disconnected but still active
	arg := strings.Split(ws.Text(), "+")
	if len(arg) != 3 {
		return output, nil
	}

	mode := arg[0]
	output.Modes = append(output.Modes, Mode{Name: mode, Active: true})

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

var errNotEdidLine = errors.New("not an edid line")

// parseEdidLine returns the partial EDID on that line
func parseEdidLine(line string) (edid string, err error) {
	if !strings.HasPrefix(line, "		") {
		return "", errNotEdidLine
	}

	ws := bufio.NewScanner(bytes.NewReader([]byte(line)))
	ws.Split(bufio.ScanWords)

	if !ws.Scan() {
		return "", fmt.Errorf("line too short, no edid part found: %s", line)
	}
	edid = ws.Text()

	if ws.Scan() {
		return "", fmt.Errorf("line too long, expected only one edid part: %s", line)
	}

	return edid, nil
}

// RandrParse returns the list of outputs parsed from the reader.
func RandrParse(rd io.Reader) (outputs Outputs, err error) {
	ls := bufio.NewScanner(rd)

	const (
		StateStart = iota
		StateOutput
		StateAdditionalProperties
		StateEdid
		StateMode
	)

	var (
		state       = StateStart
		output      Output
		currentEdid string
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
				state = StateAdditionalProperties
				continue nextLine

			case StateAdditionalProperties:
				if strings.HasPrefix(line, "	EDID:") {
					state = StateEdid
					currentEdid = ""
					continue nextLine
				}
				if !strings.HasPrefix(line, "	") {
					state = StateMode
					continue
				}
				continue nextLine

			case StateEdid:
				edid_part, err := parseEdidLine(line)
				if err == errNotEdidLine {
					monitorId, err := GenerateMonitorId(currentEdid)
					if err != nil {
						return nil, err
					}
					output.MonitorId = monitorId
					state = StateAdditionalProperties
					continue
				}
				if err != nil {
					return nil, err
				}
				currentEdid += edid_part
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

func runXrandr(extraArgs ...string) *exec.Cmd {
	args := []string{"--query", "--props"}
	args = append(args, extraArgs...)
	cmd := exec.Command("xrandr", args...)
	cmd.Stderr = os.Stderr
	return cmd
}

// GetOutputs runs `xrandr` and returns the parsed output.
func GetOutputs() (Outputs, error) {
	cmd := runXrandr("--current")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return RandrParse(bytes.NewReader(output))
}

// DetectOutputs runs `xrandr`, rescans the outputs and returns the parsed outputs.
func DetectOutputs() (Outputs, error) {
	cmd := runXrandr()
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return RandrParse(bytes.NewReader(output))
}

// BuildCommandOutputRow return a sequence of calls to `xrandr` to configure
// all named outputs in a row, left to right, given the currently active
// Outputs and a list of output names, optionally followed by "@" and the
// desired mode, e.g. LVDS1@1377x768.
func BuildCommandOutputRow(rule Rule, current Outputs) ([]*exec.Cmd, error) {
	var outputs []string
	var row bool

	switch {
	case rule.ConfigureSingle != "":
		outputs = []string{rule.ConfigureSingle}
	case len(rule.ConfigureRow) > 0:
		outputs = rule.ConfigureRow
		row = true
	case len(rule.ConfigureColumn) > 0:
		outputs = rule.ConfigureColumn
		row = false
	default:
		return nil, errors.New("empty monitor row configuration")
	}

	V("enable outputs: %v\n", outputs)

	command := "xrandr"
	enableOutputArgs := [][]string{}

	active := make(map[string]struct{})
	var lastOutput = ""
	for i, output := range outputs {
		data := strings.SplitN(output, "@", 2)
		name := data[0]
		mode := ""
		if len(data) > 1 {
			mode = data[1]
		}

		active[name] = struct{}{}

		args := []string{}
		args = append(args, "--output", name)
		if mode == "" {
			args = append(args, "--auto")
		} else {
			args = append(args, "--mode", mode)
		}

		if i > 0 {
			if row {
				args = append(args, "--right-of", lastOutput)
			} else {
				args = append(args, "--below", lastOutput)
			}
		}

		if rule.Primary == name {
			args = append(args, "--primary")
		}

		lastOutput = name
		enableOutputArgs = append(enableOutputArgs, args)
	}

	disableOutputs := make(map[string]struct{})
	for _, output := range current {
		if !output.Connected && len(output.Modes) == 0 {
			continue
		}

		// disable unneeded outputs that are still active
		if _, ok := active[output.Name]; !ok {
			disableOutputs[output.Name] = struct{}{}
		}
	}

	disableOutputArgs := [][]string{}

	// honour disable_order if present
	for _, name := range rule.DisableOrder {
		if _, ok := disableOutputs[name]; ok {
			args := []string{"--output", name, "--off"}
			disableOutputArgs = append(disableOutputArgs, args)

			delete(disableOutputs, name)
		}
	}

	// collect remaining outputs to be disabled
	for name := range disableOutputs {
		args := []string{"--output", name, "--off"}
		disableOutputArgs = append(disableOutputArgs, args)
	}

	// enable/disable all monitors in one call to xrandr
	if rule.Atomic {
		V("using one atomic call to xrandr\n")
		args := []string{}
		for _, disableArgs := range disableOutputArgs {
			args = append(args, disableArgs...)
		}
		for _, enableArgs := range enableOutputArgs {
			args = append(args, enableArgs...)
		}
		cmd := exec.Command(command, args...)
		return []*exec.Cmd{cmd}, nil
	}

	V("splitting the configuration into several calls to xrandr\n")

	// otherwise return several calls to xrandr
	cmds := []*exec.Cmd{}

	// disable an output
	if len(disableOutputArgs) > 0 {
		cmds = append(cmds, exec.Command(command, disableOutputArgs[0]...))
		disableOutputArgs = disableOutputArgs[1:]
	}

	// now for each newly enabled output, also disable another output
	for len(disableOutputArgs) > 0 || len(enableOutputArgs) > 0 {
		args := []string{}
		if len(disableOutputArgs) > 0 {
			args = append(args, disableOutputArgs[0]...)
			disableOutputArgs = disableOutputArgs[1:]
		}
		if len(enableOutputArgs) > 0 {
			args = append(args, enableOutputArgs[0]...)
			enableOutputArgs = enableOutputArgs[1:]
		}

		cmds = append(cmds, exec.Command(command, args...))
	}

	return cmds, nil
}

// DisableOutputs returns a call to `xrandr` to switch off the specified outputs.
func DisableOutputs(off Outputs) (*exec.Cmd, error) {
	if len(off) == 0 {
		return nil, nil
	}

	command := "xrandr"
	args := []string{}

	var outputs []string
	for _, output := range off {
		outputs = append(outputs, output.Name)
		args = append(args, "--output", output.Name, "--off")
	}

	V("disable outputs: %v\n", outputs)

	return exec.Command(command, args...), nil
}
