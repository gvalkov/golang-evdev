// +build linux

// Package evdev is an interface to the Linux input event system. The
// evdev interface serves the purpose of passing events generated in
// the kernel directly to userspace through character devices thata
// are typically located in /dev/input/.
//
// This package also comes with bindings to uinput, the userspace
// input subsystem. Uinput allows userspace programs to create and
// handle input devices from input events can be directly injected
// into the input subsystem.
package evdev

import (
	"os"
	"fmt"
	"unsafe"
	"strings"
	"encoding/binary"
	"bytes"
)

const _MAX_NAME_SIZE = 256

// A Linux input device from which events can be read.
type InputDevice struct {
	Fn   string      // path to input device (devnode)

	Name string      // device name
	Phys string      // physical topology of device
	File *os.File    // an open file handle to the input device

	Bustype uint16   // bus type identifier
	Vendor  uint16   // vendor identifier
	Product uint16   // product identifier
	Version uint16   // version identifier

	Capabilities map[CapabilityType][]CapabilityCode  // supported event types and codes.
}

func Open(devnode string) *InputDevice {
	file, err := os.Open(devnode)

	if err != nil {
		panic(err)
	}

	dev := InputDevice{}
	dev.Fn = devnode
	dev.File = file

	dev.set_device_info()
	dev.set_device_capabilities()

	return &dev
}

// Read and return a slice of input events from device.
func (dev *InputDevice) Read() []InputEvent {
	events := make([]InputEvent, 16)
	buffer := make([]byte, eventsize*16)

	dev.File.Read(buffer)

	b := bytes.NewBuffer(buffer)
	binary.Read(b, binary.LittleEndian, &events)

	// remove trailing structures
	for i := range events {
		if events[i].Time.Sec == 0 {
			events = append(events[:i])
			break
		}
	}

	return events
}

// Read and return a single input event.
func (dev *InputDevice) ReadOne() *InputEvent {
	event := InputEvent{}
	buffer := make([]byte, eventsize)

	dev.File.Read(buffer)

	b := bytes.NewBuffer(buffer)
	binary.Read(b, binary.LittleEndian, &event)

	return &event
}

// Get a useful description for an input device. Example:
//   InputDevice /dev/input/event3 (fd 3)
//     name Logitech USB Laser Mouse
//     phys usb-0000:00:12.0-2/input0
//     bus 0x3, vendor 0x46d, product 0xc069, version 0x110
//     events EV_KEY 1, EV_SYN 0, EV_REL 2, EV_MSC 4
func (dev *InputDevice) String() string {
	capkeys := keys(&dev.Capabilities)
	evtypes := make([]string, 0)

	for k := range capkeys {
		ev := capkeys[k]
		evtypes = append(evtypes, fmt.Sprintf("%s %d", EV[ev], ev))
	}
	evtypes_s := strings.Join(evtypes, ", ")

	return fmt.Sprintf(
		"InputDevice %s (fd %d)\n" +
		"  name %s\n" +
		"  phys %s\n" +
		"  bus 0x%04x, vendor 0x%04x, product 0x%04x, version 0x%04x\n" +
		"  events %s",
		dev.Fn, dev.File.Fd(), dev.Name, dev.Phys, dev.Bustype,
		dev.Vendor, dev.Product, dev.Version, evtypes_s)
}

// Gets the event types and event codes that the input device supports
func (dev *InputDevice) set_device_capabilities() {
    // Capabilities is a map of supported event types to lists of
    // events e.g: {1: [272, 273, 274, 275], 2: [0, 1, 6, 8]}
	// capabilities := make(map[int][]int)
	capabilities := make(map[CapabilityType][]CapabilityCode)

	evbits   := new([ (EV_MAX+1)/8]byte)
	codebits := new([(KEY_MAX+1)/8]byte)

	ioctl(dev.File.Fd(), _EVIOCGBIT(0, EV_MAX), unsafe.Pointer(evbits))

    // Build a map of the device's capabilities
	for evtype := 0; evtype < EV_MAX; evtype++ {
		if evbits[evtype/8] & (1 << uint(evtype % 8)) != 0 {
			eventcodes := make([]CapabilityCode, 0)

			ioctl(dev.File.Fd(), _EVIOCGBIT(evtype, KEY_MAX), unsafe.Pointer(codebits))
			for evcode := 0; evcode < KEY_MAX; evcode++ {
				if codebits[evcode/8] & (1 << uint(evcode % 8)) != 0 {
					c := CapabilityCode{evcode, ByEventType[evtype][evcode]}
					eventcodes = append(eventcodes, c)
				}
			}

            // capabilities[EV_KEY] = [KEY_A, KEY_B, KEY_C, ...]
			key = CapabilityType{evtype, EV[evtype]}
			capabilities[key] = eventcodes
		}
	}

	dev.Capabilities = capabilities
}

func (dev *InputDevice) resolve_capabilities

// An all-in-one function for describing an input device
func (dev *InputDevice) set_device_info() {
	info := device_info{}

	name := new([_MAX_NAME_SIZE]byte)
	phys := new([_MAX_NAME_SIZE]byte)

	ioctl(dev.File.Fd(), _EVIOCGID, unsafe.Pointer(&info))
	ioctl(dev.File.Fd(), _EVIOCGNAME, unsafe.Pointer(name))
	ioctl(dev.File.Fd(), _EVIOCGPHYS, unsafe.Pointer(phys))

	dev.Name = string(name[:])
	dev.Phys = string(phys[:])

	dev.Vendor  = info.vendor
	dev.Bustype = info.bustype
	dev.Product = info.product
	dev.Version = info.version
}

type CapabilityType struct {
	Type uint8
	Type string
}

type CapabilityCode struct {
	Code uint
	Name string
}

// Corresponds to the input_id struct
type device_info struct {
	bustype, vendor, product, version uint16
}

// Return the keys of a map as a slice (dict.keys())
func keys (cap *map[int][]int) []int {
	slice := make([]int, 0)

	for key := range *cap {
		slice = append(slice, key)
	}

	return slice
}
