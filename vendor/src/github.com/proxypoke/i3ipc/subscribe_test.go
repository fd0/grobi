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

func TestInit(t *testing.T) {
	for _, s := range eventSockets {
		if !s.open {
			t.Error("Init failed: closed event socket found.")
		}
	}
	if len(eventSockets) != int(eventmax) {
		t.Errorf("Too much or not enough event sockets. Got %d, expected %d.\n",
			len(eventSockets), int(eventmax))
	}

	_, err := Subscribe(I3WorkspaceEvent)
	if err != nil {
		t.Errorf("Failed to subscribe: %f\n")
	}
	// TODO: A test to ensure that subscriptions work as intended.
}
