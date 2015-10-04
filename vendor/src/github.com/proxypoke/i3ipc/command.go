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
	"encoding/json"
)

// Error for replies from a command to i3.
type CommandError string

func (self CommandError) Error() string {
	return string(self)
}

// Struct for replies from command messages.
type commandReply struct {
	Success bool
	Error   string
}

// Send a command to i3.
// FIXME: Doesn't support chained commands yet.
func (self *IPCSocket) Command(action string) (success bool, err error) {
	json_reply, err := self.Raw(I3Command, action)
	if err != nil {
		return
	}

	var cmd_reply []commandReply
	err = json.Unmarshal(json_reply, &cmd_reply)
	if err != nil {
		return
	}

	success = cmd_reply[0].Success
	if cmd_reply[0].Error == "" {
		err = nil
	} else {
		err = CommandError(cmd_reply[0].Error)
	}

	return
}
