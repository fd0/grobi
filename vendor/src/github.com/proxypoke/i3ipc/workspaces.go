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

// Struct representing a workspace. For documentation of the fields,
// refer to http://i3wm.org/docs/ipc.html#_workspaces_reply.
type Workspace struct {
	Num     int32
	Name    string
	Visible bool
	Focused bool
	Rect    Rect
	Output  string
	Urgent  bool
}

// The geometry of a window, output or workspace.
type Rect struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

// GetWorkspaces fetches a list of all current workspaces.
func (self *IPCSocket) GetWorkspaces() (workspaces []Workspace, err error) {
	json_reply, err := self.Raw(I3GetWorkspaces, "")
	if err != nil {
		return
	}

	err = json.Unmarshal(json_reply, &workspaces)
	return
}
