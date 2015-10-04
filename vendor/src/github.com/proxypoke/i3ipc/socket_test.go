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
	"testing"
)

func TestGetIPCSocket(t *testing.T) {
	ipc, err := GetIPCSocket()
	if err != nil {
		t.Errorf("Failed to acquire the IPC socket: %v", err)
	}
	ipc.Close()
	if ipc.open {
		t.Error("IPC socket appears open after closing.")
	}
}

func TestRaw(t *testing.T) {
	ipc, _ := GetIPCSocket()

	_, err := ipc.Raw(I3GetWorkspaces, "")
	if err != nil {
		t.Errorf("Raw message sending failed: %v", err)
	}
}
