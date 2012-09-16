package evdev

import (
	"syscall"
)


type InputEvent struct {
	time syscall.Timeval
	etype uint16
	code  uint16
	value int32
}