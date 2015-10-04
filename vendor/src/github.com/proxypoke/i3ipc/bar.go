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

// Struct representing the configuration of a bar. For documentation of the
// fields, refer to http://i3wm.org/docs/ipc.html#_bar_config_reply.
type I3Bar struct {
	Id                string
	Mode              string
	Position          string
	Status_Command    string
	Font              string
	Workspace_Buttons bool
	Verbose           bool
	Colors            Colors
}

// Struct representing colors as used in I3Bar.
type Colors struct {
	Focused_Workspace_Border string
	Focused_Workspace_Bg     string
	Focused_Workspace_Text   string
}

// GetBarIds fetches a list of IDs for all configured bars.
func (self *IPCSocket) GetBarIds() (ids []string, err error) {
	json_reply, err := self.Raw(I3GetBarConfig, "")
	if err != nil {
		return
	}

	err = json.Unmarshal(json_reply, &ids)
	return
}

// GetBarConfig returns the configuration of the bar with the given ID.
func (self *IPCSocket) GetBarConfig(id string) (bar I3Bar, err error) {
	json_reply, err := self.Raw(I3GetBarConfig, id)
	if err != nil {
		return
	}

	err = json.Unmarshal(json_reply, &bar)
	return
}
