// Author: slowpoke <mail plus git at slowpoke dot io>
// Repository: https://github.com/proxypoke/i3ipc
//
// This program is free software under the terms of the
// Do What The Fuck You Want To Public License.
// It comes without any warranty, to the extent permitted by
// applicable law. For a copy of the license, see COPYING or
// head to http://sam.zoy.org/wtfpl/COPYING.

package i3ipc

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"unsafe"
)

const (
	// Magic string for the IPC API.
	_MAGIC string = "i3-ipc"
	// The length of the i3 message "header" is 14 bytes: 6 for the _MAGIC
	// string, 4 for the length of the payload (int32 in native byte order) and
	// another 4 for the message type (also int32 in NBO).
	_HEADERLEN = 14
)

// A message from i3. Can either be a Reply or an Event.
type Message struct {
	Payload []byte
	IsEvent bool
	Type    int32
}

// The types of messages that Raw() accepts.
type MessageType int32

const (
	I3Command MessageType = iota
	I3GetWorkspaces
	I3Subscribe
	I3GetOutputs
	I3GetTree
	I3GetMarks
	I3GetBarConfig
	I3GetVersion
)

// Error for unknown message types.
type MessageTypeError string

func (self MessageTypeError) Error() string {
	return string(self)
}

// Error for communiation failures.
type MessageError string

func (self MessageError) Error() string {
	return string(self)
}

// An Unix socket to communicate with i3.
type IPCSocket struct {
	socket      net.Conn
	open        bool
	subscribers []chan Event
}

// Close the connection to the underlying Unix socket.
func (self *IPCSocket) Close() error {
	self.open = false
	return self.socket.Close()
}

// Create a new IPC socket.
func GetIPCSocket() (ipc *IPCSocket, err error) {
	var out bytes.Buffer
	ipc = &IPCSocket{}

	cmd := exec.Command("i3", "--get-socketpath")
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}

	path := strings.TrimSpace(out.String())
	sock, err := net.Dial("unix", path)
	ipc.socket = sock
	ipc.open = true
	return
}

// Receive a raw json bytestring from the socket and return a Message.
func (self *IPCSocket) recv() (msg Message, err error) {
	header := make([]byte, _HEADERLEN)
	n, err := self.socket.Read(header)

	// Check if this is a valid i3 message.
	if n != _HEADERLEN || err != nil {
		return
	}
	magic_string := string(header[:len(_MAGIC)])
	if magic_string != _MAGIC {
		err = MessageError(fmt.Sprintf(
			"Invalid magic string: got %q, expected %q.",
			magic_string, _MAGIC))
		return
	}

	var bytelen [4]byte
	// Copy the byte values from the slice into the byte array. This is
	// necessary because the adress of a slice does not point to the actual
	// values in memory.
	for i, b := range header[len(_MAGIC) : len(_MAGIC)+4] {
		bytelen[i] = b
	}
	length := *(*int32)(unsafe.Pointer(&bytelen))

	msg.Payload = make([]byte, length)
	n, err = self.socket.Read(msg.Payload)
	if n != int(length) || err != nil {
		return
	}

	// Figure out the type of message.
	var bytetype [4]byte
	for i, b := range header[len(_MAGIC)+4 : len(_MAGIC)+8] {
		bytetype[i] = b
	}
	type_ := *(*uint32)(unsafe.Pointer(&bytetype))

	// Reminder: event messages have the highest bit of the type set to 1
	if type_>>31 == 1 {
		msg.IsEvent = true
	}
	// Use the remaining bits
	msg.Type = int32(type_ & 0x7F)

	return
}

// Send raw messages to i3. Returns a json bytestring.
func (self *IPCSocket) Raw(type_ MessageType, args string) (json_reply []byte, err error) {
	// Set up the parts of the message.
	var (
		message  []byte = []byte(_MAGIC)
		payload  []byte = []byte(args)
		length   int32  = int32(len(payload))
		bytelen  [4]byte
		bytetype [4]byte
	)

	// Black Magic™.
	bytelen = *(*[4]byte)(unsafe.Pointer(&length))
	bytetype = *(*[4]byte)(unsafe.Pointer(&type_))

	for _, b := range bytelen {
		message = append(message, b)
	}
	for _, b := range bytetype {
		message = append(message, b)
	}
	for _, b := range payload {
		message = append(message, b)
	}

	_, err = self.socket.Write(message)
	if err != nil {
		return
	}

	msg, err := self.recv()
	if err == nil {
		json_reply = msg.Payload
	}
	if msg.IsEvent {
		err = MessageTypeError("Received an event instead of a reply.")
	}
	return
}
