package main

import (
	"errors"
	"fmt"
	"path"
	"sync"
	"sync/atomic"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

// XInput protocol constants, see X11/extensions/XI.h, XIproto.h and XI2proto.h.
// xgb has no XInput package, the few requests grobi needs are built by hand.
const (
	xiExtensionName = "XInputExtension"

	xiSelectExtensionEvent = 6  // XI 1.x request
	xiQueryVersion         = 47 // XI 2.x request
	xiQueryDevice          = 48 // XI 2.x request

	// XI version grobi announces, XI2 requests need at least 2.0
	xiMajorVersion = 2
	xiMinorVersion = 0

	xiAllDevices     = 0
	xiMasterPointer  = 1
	xiMasterKeyboard = 2

	// Wire format. Sizes and offsets are in bytes, lengths on the wire are
	// in units of xUnit bytes.
	xUnit      = 4
	card16Size = 2

	// request header: major opcode, minor opcode, length
	reqMajorOpcode = 0
	reqMinorOpcode = 1
	reqLength      = 2

	// xXIQueryVersionReq
	queryVersionReqSize  = 8
	queryVersionReqMajor = 4
	queryVersionReqMinor = 6
	// xXIQueryVersionReply
	queryVersionReplyMajor = 8

	// xXIQueryDeviceReq
	queryDeviceReqSize     = 8
	queryDeviceReqDeviceID = 4
	// xXIQueryDeviceReply, followed by num_devices xXIDeviceInfo
	queryDeviceReplySize       = 32
	queryDeviceReplyNumDevices = 8
	// xXIDeviceInfo, followed by the padded name and num_classes classes
	deviceInfoSize       = 12
	deviceInfoUse        = 2
	deviceInfoNumClasses = 6
	deviceInfoNameLen    = 8
	// xXIAnyInfo, the common header of all device classes
	classHeaderSize = 4
	classLength     = 2 // in xUnit, including the header

	// xSelectExtensionEventReq with one event class
	selectEventReqSize       = 16
	selectEventReqWindow     = 4
	selectEventReqNumClasses = 8
	selectEventReqClasses    = 12

	// xiDevicePresenceNotify is the offset of DevicePresenceNotify from the
	// extension's first event. It is a plain 32 byte event (unlike the XI2
	// generic events, which xgb can't read), sent when an input device is
	// added, removed, enabled or disabled.
	xiDevicePresenceNotify = 15

	// xiDevicePresenceClass is the event class for DevicePresenceNotify:
	// device id 256 ("no device") and _devicePresence (0).
	xiDevicePresenceClass = 0x10000
)

// InputDevices is the list of names of the attached input devices (XInput
// slave and floating devices, master devices are left out). A device which
// has several functions (e.g. a dock with audio buttons) may be listed more
// than once.
type InputDevices []string

// Present returns true iff an input device name matches the pattern.
func (devs InputDevices) Present(pattern string) bool {
	for _, name := range devs {
		if m, _ := path.Match(pattern, name); m {
			return true
		}
	}
	return false
}

// xinputInit checks for the XInput extension and announces XI 2.0 support,
// which the server requires before XI2 requests. It returns the extension's
// major opcode and first event number.
func xinputInit(X *xgb.Conn) (opcode, firstEvent byte, err error) {
	ext, err := xproto.QueryExtension(X, uint16(len(xiExtensionName)), xiExtensionName).Reply()
	if err != nil {
		return 0, 0, err
	}
	if !ext.Present {
		return 0, 0, errors.New("no XInput extension on the X server")
	}

	buf := newXIRequest(ext.MajorOpcode, xiQueryVersion, queryVersionReqSize)
	xgb.Put16(buf[queryVersionReqMajor:], xiMajorVersion)
	xgb.Put16(buf[queryVersionReqMinor:], xiMinorVersion)

	reply, err := request(X, buf, true)
	if err != nil {
		return 0, 0, fmt.Errorf("XIQueryVersion: %w", err)
	}
	if len(reply) < queryVersionReplyMajor+card16Size || xgb.Get16(reply[queryVersionReplyMajor:]) < xiMajorVersion {
		return 0, 0, errors.New("no XInput 2 support on the X server")
	}

	return ext.MajorOpcode, ext.FirstEvent, nil
}

// newXIRequest returns a request of size bytes with the header filled in.
func newXIRequest(opcode, minor byte, size int) []byte {
	buf := make([]byte, size)
	buf[reqMajorOpcode] = opcode
	buf[reqMinorOpcode] = minor
	xgb.Put16(buf[reqLength:], uint16(size/xUnit))
	return buf
}

// request sends a raw request and returns the reply, or the error for a
// request without reply.
func request(X *xgb.Conn, buf []byte, hasReply bool) ([]byte, error) {
	cookie := X.NewCookie(true, hasReply)
	X.NewRequest(buf, cookie)
	if !hasReply {
		return nil, cookie.Check()
	}
	return cookie.Reply()
}

// queryInputDevices returns the attached input devices.
func queryInputDevices(X *xgb.Conn, opcode byte) (InputDevices, error) {
	buf := newXIRequest(opcode, xiQueryDevice, queryDeviceReqSize)
	xgb.Put16(buf[queryDeviceReqDeviceID:], xiAllDevices)

	reply, err := request(X, buf, true)
	if err != nil {
		return nil, fmt.Errorf("XIQueryDevice: %w", err)
	}
	return parseQueryDeviceReply(reply)
}

// parseQueryDeviceReply parses an XIQueryDevice reply.
func parseQueryDeviceReply(buf []byte) (InputDevices, error) {
	errShort := errors.New("XIQueryDevice reply too short")
	if len(buf) < queryDeviceReplySize {
		return nil, errShort
	}

	n := int(xgb.Get16(buf[queryDeviceReplyNumDevices:]))
	devs := InputDevices{}
	b := queryDeviceReplySize
	for range n {
		if len(buf) < b+deviceInfoSize {
			return nil, errShort
		}
		use := xgb.Get16(buf[b+deviceInfoUse:])
		numClasses := int(xgb.Get16(buf[b+deviceInfoNumClasses:]))
		nameLen := int(xgb.Get16(buf[b+deviceInfoNameLen:]))
		b += deviceInfoSize

		if len(buf) < b+xgb.Pad(nameLen) {
			return nil, errShort
		}
		name := string(buf[b : b+nameLen])
		b += xgb.Pad(nameLen)

		// skip the device classes
		for range numClasses {
			if len(buf) < b+classHeaderSize {
				return nil, errShort
			}
			b += xUnit * int(xgb.Get16(buf[b+classLength:]))
			if len(buf) < b {
				return nil, errShort
			}
		}

		if use == xiMasterPointer || use == xiMasterKeyboard {
			continue
		}
		devs = append(devs, name)
	}

	return devs, nil
}

// inputsFailed is set once reading the input devices failed, so the error
// is only logged once.
var inputsFailed atomic.Bool

// ReadInputDevices returns the attached input devices. It returns nil if they
// can't be read.
func ReadInputDevices() InputDevices {
	devs, err := readInputDevices()
	if err != nil {
		if !inputsFailed.Swap(true) {
			V("input devices: %v\n", err)
		}
		return nil
	}
	inputsFailed.Store(false)
	return devs
}

// inputConn is the X connection for reading the input devices. It is kept
// open, the input devices are read on every poll.
var inputConn struct {
	sync.Mutex
	X       *xgb.Conn
	opcode  byte
	initErr error // XInput missing, permanent for this connection
}

func readInputDevices() (InputDevices, error) {
	c := &inputConn
	c.Lock()
	defer c.Unlock()

	if c.X == nil {
		X, err := xgb.NewConn()
		if err != nil {
			return nil, err
		}
		c.X = X
		c.opcode, _, c.initErr = xinputInit(X)
	}
	if c.initErr != nil {
		return nil, c.initErr
	}

	devs, err := queryInputDevices(c.X, c.opcode)
	var xerr xgb.Error
	if err != nil && !errors.As(err, &xerr) {
		// the connection is broken, reconnect on the next call
		c.X.Close()
		c.X = nil
	}
	return devs, err
}

// devicePresenceEvent is an XInput DevicePresenceNotify event.
type devicePresenceEvent []byte

func (ev devicePresenceEvent) Bytes() []byte  { return ev }
func (ev devicePresenceEvent) String() string { return "DevicePresenceNotify" }

// registerInputEvents checks for XInput on X and registers the
// DevicePresenceNotify event with xgb. It writes the global
// xgb.NewEventFuncs, so it must run before events are selected on any
// connection, like randr.Init. It returns the XInput major opcode.
func registerInputEvents(X *xgb.Conn) (byte, error) {
	opcode, firstEvent, err := xinputInit(X)
	if err != nil {
		return 0, err
	}

	evNum := int(firstEvent) + xiDevicePresenceNotify
	xgb.NewEventFuncs[evNum] = func(buf []byte) xgb.Event { return devicePresenceEvent(buf) }
	return opcode, nil
}

// selectInputEvents asks for a DevicePresenceNotify event on root whenever an
// input device is added, removed, enabled or disabled.
func selectInputEvents(X *xgb.Conn, opcode byte, root xproto.Window) error {
	buf := newXIRequest(opcode, xiSelectExtensionEvent, selectEventReqSize)
	xgb.Put32(buf[selectEventReqWindow:], uint32(root))
	xgb.Put16(buf[selectEventReqNumClasses:], 1)
	xgb.Put32(buf[selectEventReqClasses:], xiDevicePresenceClass)
	if _, err := request(X, buf, false); err != nil {
		return fmt.Errorf("SelectExtensionEvent: %w", err)
	}
	return nil
}
