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

// GetMarks returns a list of marks which are currently set.
func (self *IPCSocket) GetMarks() (marks []string, err error) {
	json_reply, err := self.Raw(I3GetMarks, "")
	if err != nil {
		return
	}

	err = json.Unmarshal(json_reply, &marks)
	return
}
