package evdev

import (
	"os"
	"fmt"
	"unsafe"
	"strings"
	"encoding/binary"
	"bytes"
)

const MAX_NAME_SIZE = 256

// An InputDevice is a Linux input device from which events can be
// read.
type InputDevice struct {
	Fn   string      // path to input device (devnode)

	Name string      // device name
	Phys string      // physical topology of device
	File *os.File    // an open file handle to the input device

	Bustype uint16   // bus type identifier
	Vendor  uint16   // vendor identifier
	Product uint16   // product identifier
	Version uint16   // version identifier

	Capabilities map[int][]int  // supported event types and codes.
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

// Read multiple input events from device
func (dev *InputDevice) Read() []InputEvent {
	events := make([]InputEvent, 16)
	buffer := make([]byte, eventsize*16)

	dev.File.Read(buffer)

	b := bytes.NewBuffer(buffer)
	binary.Read(b, binary.LittleEndian, &events)

	for i := range events {
		if events[i].Time.Sec == 0 {
			events = append(events[:i])
			break
		}
	}

	return events
}

// Read and return a single input event
func (dev *InputDevice) ReadOne() *InputEvent {
	event := InputEvent{}
	buffer := make([]byte, eventsize)

	dev.File.Read(buffer)

	b := bytes.NewBuffer(buffer)
	binary.Read(b, binary.LittleEndian, &event)

	return &event
}

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
		"  bus 0x%x, vendor 0x%x, product 0x%x, version 0x%x\n" +
		"  events %s",
		dev.Fn, dev.File.Fd(), dev.Name, dev.Phys, dev.Bustype,
		dev.Vendor, dev.Product, dev.Version, evtypes_s)
}

// Gets the event types and event codes that the input device supports
func (dev *InputDevice) set_device_capabilities() {
    // Capabilities is a map of supported event types to lists of
    // events e.g: {1: [272, 273, 274, 275], 2: [0, 1, 6, 8]}
	capabilities := make(map[int][]int)

	evbits   := new([ (EV_MAX+1)/8]byte)
	codebits := new([(KEY_MAX+1)/8]byte)

	ioctl(dev.File.Fd(), EVIOCGBIT(0, EV_MAX), unsafe.Pointer(evbits))

    // Build a map of the device's capabilities
	for evtype := 0; evtype < EV_MAX; evtype++ {
		if evbits[evtype/8] & (1 << uint(evtype % 8)) != 0 {
			eventcodes := make([]int, 0)

			ioctl(dev.File.Fd(), EVIOCGBIT(evtype, KEY_MAX), unsafe.Pointer(codebits))
			for evcode := 0; evcode < KEY_MAX; evcode++ {
				if codebits[evcode/8] & (1 << uint(evcode % 8)) != 0 {
					eventcodes = append(eventcodes, evcode)
				}
			}

            // capabilities[EV_KEY] = [KEY_A, KEY_B, KEY_C, ...]
			capabilities[evtype] = eventcodes
		}
	}

	dev.Capabilities = capabilities
}

// An all-in-one function for describing an input device
func (dev *InputDevice) set_device_info() {
	info := device_info{}

	name := new([MAX_NAME_SIZE]byte)
	phys := new([MAX_NAME_SIZE]byte)

	ioctl(dev.File.Fd(), EVIOCGID, unsafe.Pointer(&info))
	ioctl(dev.File.Fd(), EVIOCGNAME, unsafe.Pointer(name))
	ioctl(dev.File.Fd(), EVIOCGPHYS, unsafe.Pointer(phys))

	dev.Name = string(name[:])
	dev.Phys = string(phys[:])

	dev.Vendor  = info.vendor
	dev.Bustype = info.bustype
	dev.Product = info.product
	dev.Version = info.version
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
