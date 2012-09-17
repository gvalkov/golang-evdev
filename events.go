package evdev

import (
	"syscall"
	"unsafe"
	"fmt"
)

type InputEvent struct {
	Time syscall.Timeval  // time in seconds since epoch at which event occurred
	Type  uint16  // event type - one of ecodes.EV_*
	Code  uint16  // event code related to the event type
	Value int32   // event value related to the event type
}

func (ev *InputEvent) String() string {
	return fmt.Sprintf("event at %d.%d, code %02d, type %02d, val %02d",
		   ev.Time.Sec, ev.Time.Usec, ev.Code, ev.Type, ev.Value)
}

var eventsize = int(unsafe.Sizeof(InputEvent{}))


type KeyEventState uint8
const (
	KeyUp KeyEventState = 0x0
	KeyDown KeyEventState = 0x1
	KeyHold KeyEventState = 0x2
)

// KeyEvents are used to describe state changes of keyboards, buttons,
// or other key-like devices.
type KeyEvent struct {
	Event *InputEvent
	Scancode uint16
	Keycode  uint16
	State KeyEventState
}

func (kev *KeyEvent) New(ev *InputEvent) {
	kev.Event = ev
	kev.Keycode = 0 // :todo
	kev.Scancode = ev.Code

	switch ev.Value {
	case 0: kev.State = KeyUp
	case 2: kev.State = KeyHold
	case 1: kev.State = KeyDown
	}
}

func (kev *KeyEvent) String() string {
	state := "unknown"

	switch kev.State {
	case KeyUp: state = "up"
	case KeyHold: state = "hold"
	case KeyDown: state = "down"
	}

	return fmt.Sprintf("key event at %d.%d, %d (%d), ()",
		               kev.Event.Time.Sec, kev.Event.Time.Usec,
		               kev.Scancode, kev.Event.Code, state)
}
