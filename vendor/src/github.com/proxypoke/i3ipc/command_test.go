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

func TestCommand(t *testing.T) {
	ipc, _ := GetIPCSocket()
	defer ipc.Close()

	// `exec /bin/true` is a good NOP operation for testing
	success, err := ipc.Command("exec /bin/true")
	if !success {
		t.Error("Unsuccessful command.")
	}
	if err != nil {
		t.Errorf("An error occurred during command: %v", err)
	}
}
