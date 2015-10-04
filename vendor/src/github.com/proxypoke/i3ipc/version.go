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

// Struct representing the version of i3. For documentation of the
// fields, refer to http://i3wm.org/docs/ipc.html#_version_reply.
type I3Version struct {
	Major          int32
	Minor          int32
	Patch          int32
	Human_Readable string
}

// GetVersion fetches the version of i3.
func (self *IPCSocket) GetVersion() (version I3Version, err error) {
	json_reply, err := self.Raw(I3GetVersion, "")
	if err != nil {
		return
	}

	err = json.Unmarshal(json_reply, &version)
	return
}
