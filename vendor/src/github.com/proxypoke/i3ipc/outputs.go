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

// Struct representing an output. For documentation of the fields,
// refer to http://i3wm.org/docs/ipc.html#_outputs_reply.
type Output struct {
	Name              string
	Active            bool
	Rect              Rect
	Current_Workspace string
	//Primary           bool
}

// GetOutputs fetches the list of current outputs.
func (self *IPCSocket) GetOutputs() (outputs []Output, err error) {
	json_reply, err := self.Raw(I3GetOutputs, "")
	if err != nil {
		return
	}

	err = json.Unmarshal(json_reply, &outputs)
	if err == nil {
		return
	}
	// Outputs which aren't displaying any workspace will have JSON-null set as
	// their value for current_workspace. Since Go's equivalent, nil, can't be
	// assigned to strings, it will cause Unmarshall to return with an
	// UnmarshalTypeError, but otherwise correctly unmarshal the JSON input. We
	// simply ignore this error due to this reason.
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		err = nil
	}
	return
}
