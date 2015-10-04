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
	"log"
)

// Type for subscribable events.
type EventType int32

// Enumeration of currently available event types.
const (
	I3WorkspaceEvent EventType = iota
	I3OutputEvent
	I3ModeEvent
	// private value used for setting up internal stuff in init()
	// The idea is that if there's a new type of event added to i3, it only
	// needs to be added here and in the payloads slice below, and the rest of
	// the code won't need to change.
	eventmax
)

// This slice is used to map event types to their string representation.
var payloads []string = []string{"workspace", "output", "mode"}

// Dynamically add an event type by defining a name for it. Just in case i3 adds
// a new one and this library hasn't been updated yet. Returns the EventType
// which gets assigned to it.
//
// XXX: If you use this to add more than one new event type, add them in the
// RIGHT ORDER. I hope this case never pops up (because that would mean that
// this library is severely outdated), but I thought I'd warn you anyways.
func AddEventType(name string) (type_ EventType) {
	payloads = append(payloads, name)
	return EventType(len(payloads) - 1)
}

// Event describes an event reply from i3.
type Event struct {
	Type EventType
	// "change" is the name of the single field of the JSON map that i3 sends
	// when an event occurs, describing what happened.
	Change string
}

// Struct for replies from subscribe messages.
type subscribeReply struct {
	Success bool
}

// Subscription related errors.
type SubscribeError string

func (self SubscribeError) Error() string {
	return string(self)
}

// Private subscribe function. Sets up the socket.
func (self *IPCSocket) subscribe(type_ EventType) (err error) {
	json_reply, err := self.Raw(I3Subscribe, "[\""+payloads[type_]+"\"]")
	if err != nil {
		return
	}

	var subs_reply subscribeReply
	err = json.Unmarshal(json_reply, &subs_reply)
	if err != nil {
		return
	}

	if !subs_reply.Success {
		// TODO: Better error description.
		err = SubscribeError("Could not subscribe.")
	}
	return
}

// Subscribe to an event type. Returns a channel from which events can be read.
func Subscribe(type_ EventType) (subs chan Event, err error) {
	if type_ >= eventmax || type_ < 0 {
		err = SubscribeError("No such event type.")
		return
	}
	subs = make(chan Event)
	eventSockets[type_].subscribers = append(
		eventSockets[type_].subscribers, subs)
	return
}

// Listen for events on this socket, multiplexing them to all subscribers.
//
// XXX: This will cause all messages which are not events to be DROPPED.
func (self *IPCSocket) listen() {
	for {
		if !self.open {
			break
		}
		msg, err := self.recv()
		// XXX: This ignores all errors. Maybe a FIXME, maybe not.
		if err != nil {
			continue
		}
		// Drop non-event messages.
		if !msg.IsEvent {
			continue
		}

		var event Event
		event.Type = EventType(msg.Type)
		err = json.Unmarshal(msg.Payload, &event)

		// Send each subscriber the event in a nonblocking manner.
		for _, subscriber := range self.subscribers {
			select {
			case subscriber <- event: // NOP
			default:
				// If the event can't be written, just ignore this
				// subscriber.
			}
		}
	}
}

var eventSockets []*IPCSocket

func init() {
	// Check whether we have as much payloads as we have event types. You know,
	// just in case I'm coding on my third Club-Mate at 0400 in the morning when
	// updating this lib.
	if len(payloads) != int(eventmax) {
		log.Fatalf("Too much or not enough payloads: got %d, expected %d.\n",
			len(payloads), int(eventmax))
	}

	// Set up an IPCSocket to receive events for every type of event.
	var ev EventType
	for ; ev < eventmax; ev++ {
		sock, err := GetIPCSocket()
		if err != nil {
			log.Fatalf("Can't get i3 socket. Please make sure i3 is running. %v.", err)
		}
		err = sock.subscribe(ev)
		if err != nil {
			log.Fatalf("Can't subscribe: %v", err)
		}
		go sock.listen()
		if err != nil {
			log.Fatalf("Can't set up event sockets: %v", err)
		}

		eventSockets = append(eventSockets, sock)
	}
}
