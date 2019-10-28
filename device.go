// +build linux

package evdev

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

// A Linux input device from which events can be read.
type InputDevice struct {
	Fn string // path to input device (devnode)

	Name string   // device name
	Phys string   // physical topology of device
	File *os.File // an open file handle to the input device

	Bustype uint16 // bus type identifier
	Vendor  uint16 // vendor identifier
	Product uint16 // product identifier
	Version uint16 // version identifier

	EvdevVersion int // evdev protocol version

	Capabilities     map[CapabilityType][]CapabilityCode // supported event types and codes.
	CapabilitiesFlat map[int][]int
}

// Open an evdev input device.
func Open(devnode string) (*InputDevice, error) {
	f, err := os.Open(devnode)
	if err != nil {
		return nil, err
	}

	dev := InputDevice{}
	dev.Fn = devnode
	dev.File = f

	err = dev.setDeviceInfo()
	if err != nil {
		log.Printf("setDeviceInfo error: %v", err)
		return nil, err
	}
	err = dev.setDeviceCapabilities()
	if err != nil {
		log.Printf("setDeviceCapabilities error: %v", err)
		return nil, err
	}

	return &dev, nil
}

// Read and return a slice of input events from device.
func (dev *InputDevice) Read() ([]InputEvent, error) {
	events := make([]InputEvent, 16)
	buffer := make([]byte, eventsize*16)

	_, err := dev.File.Read(buffer)
	if err != nil {
		log.Printf("Read().dev.File.Read error: %v", err)
		return events, err
	}

	b := bytes.NewBuffer(buffer)
	err = binary.Read(b, binary.LittleEndian, &events)
	if err != nil {
		log.Printf("Read().binary.Read error: %v", err)
		return events, err
	}

	// remove trailing structures
	for i := range events {
		if events[i].Time.Sec == 0 {
			events = append(events[:i])
			break
		}
	}

	return events, err
}

// Read and return a single input event.
func (dev *InputDevice) ReadOne() (*InputEvent, error) {
	event := InputEvent{}
	buffer := make([]byte, eventsize)

	_, err := dev.File.Read(buffer)
	if err != nil {
		log.Printf("ReadOne().dev.File.Read error: %v", err)
		return &event, err
	}

	b := bytes.NewBuffer(buffer)
	err = binary.Read(b, binary.LittleEndian, &event)
	if err != nil {
		log.Printf("ReadOne().binary.Read error: %v", err)
		return &event, err
	}

	return &event, err
}

// Get a useful description for an input device. Example:
//   InputDevice /dev/input/event3 (fd 3)
//     name Logitech USB Laser Mouse
//     phys usb-0000:00:12.0-2/input0
//     bus 0x3, vendor 0x46d, product 0xc069, version 0x110
//     events EV_KEY 1, EV_SYN 0, EV_REL 2, EV_MSC 4
func (dev *InputDevice) String() string {
	evtypes := make([]string, 0)

	for ev := range dev.Capabilities {
		evtypes = append(evtypes, fmt.Sprintf("%s %d", ev.Name, ev.Type))
	}
	evtypes_s := strings.Join(evtypes, ", ")

	return fmt.Sprintf(
		"InputDevice %s (fd %d)\n"+
			"  name %s\n"+
			"  phys %s\n"+
			"  bus 0x%04x, vendor 0x%04x, product 0x%04x, version 0x%04x\n"+
			"  events %s",
		dev.Fn, dev.File.Fd(), dev.Name, dev.Phys, dev.Bustype,
		dev.Vendor, dev.Product, dev.Version, evtypes_s)
}

// Gets the event types and event codes that the input device supports.
func (dev *InputDevice) setDeviceCapabilities() error {
	// Capabilities is a map of supported event types to lists of
	// events e.g: {1: [272, 273, 274, 275], 2: [0, 1, 6, 8]}
	// capabilities := make(map[int][]int)
	capabilities := make(map[CapabilityType][]CapabilityCode)

	evbits := new([(EV_MAX + 1) / 8]byte)
	codebits := new([(KEY_MAX + 1) / 8]byte)
	// absbits  := new([6]byte)

	err := ioctl(dev.File.Fd(), uintptr(EVIOCGBIT(0, EV_MAX)), unsafe.Pointer(evbits))
	if err != 0 {
		log.Printf("setDeviceCapabilities().ioctl(EVIOCGBIT(0, EV_MAX)) error: %v", err)
		return err
	}

	// Build a map of the device's capabilities
	for evtype := 0; evtype < EV_MAX; evtype++ {
		if evbits[evtype/8]&(1<<uint(evtype%8)) != 0 {
			eventcodes := make([]CapabilityCode, 0)

			err = ioctl(dev.File.Fd(), uintptr(EVIOCGBIT(evtype, KEY_MAX)), unsafe.Pointer(codebits))
			if err != 0 {
				log.Printf("setDeviceCapabilities().ioctl(EVIOCGBIT(evtype, KEY_MAX)) evtype=%v error: %v", evtype, err)
				//return err
				continue
			}
			for evcode := 0; evcode < KEY_MAX; evcode++ {
				if codebits[evcode/8]&(1<<uint(evcode%8)) != 0 {
					c := CapabilityCode{evcode, ByEventType[evtype][evcode]}
					eventcodes = append(eventcodes, c)
				}
			}

			// capabilities[EV_KEY] = [KEY_A, KEY_B, KEY_C, ...]
			key := CapabilityType{evtype, EV[evtype]}
			capabilities[key] = eventcodes
		}
	}

	dev.Capabilities = capabilities
	return nil
}

// An all-in-one function for describing an input device.
func (dev *InputDevice) setDeviceInfo() error {
	info := device_info{}

	name := new([MAX_NAME_SIZE]byte)
	phys := new([MAX_NAME_SIZE]byte)

	err := ioctl(dev.File.Fd(), uintptr(EVIOCGID), unsafe.Pointer(&info))
	if err != 0 {
		log.Printf("setDeviceInfo().ioctl(EVIOCGID) error: %v", err)
		return err
	}

	err = ioctl(dev.File.Fd(), uintptr(EVIOCGNAME), unsafe.Pointer(name))
	if err != 0 {
		log.Printf("setDeviceInfo().ioctl(EVIOCGNAME) error: %v", err)
		return err
	}

	// it's ok if the topology info is not available
	err = ioctl(dev.File.Fd(), uintptr(EVIOCGPHYS), unsafe.Pointer(phys))
	if err != 0 {
		log.Printf("setDeviceInfo().ioctl(EVIOCGPHYS) error: %v", err)
	}

	dev.Name = bytesToString(name)
	dev.Phys = bytesToString(phys)

	dev.Vendor = info.vendor
	dev.Bustype = info.bustype
	dev.Product = info.product
	dev.Version = info.version

	evVersion := new(int)
	err = ioctl(dev.File.Fd(), uintptr(EVIOCGVERSION), unsafe.Pointer(evVersion))
	if err != 0 {
		log.Printf("setDeviceInfo().ioctl(EVIOCGVERSION) error: %v", err)
		return err
	}
	dev.EvdevVersion = *evVersion

	return nil
}

// Get repeat rate as a two element array.
//   [0] repeat rate in characters per second
//   [1] amount of time that a key must be depressed before it will start
//       to repeat (in milliseconds)
func (dev *InputDevice) GetRepeatRate() *[2]uint {
	repeat_delay := new([2]uint)
	err := ioctl(dev.File.Fd(), uintptr(EVIOCGREP), unsafe.Pointer(repeat_delay))
	if err != 0 {
		log.Printf("GetRepeatRate().ioctl error: %v", err)
	}
	return repeat_delay
}

// Set repeat rate and delay.
func (dev *InputDevice) SetRepeatRate(repeat, delay uint) {
	repeat_delay := new([2]uint)
	repeat_delay[0], repeat_delay[1] = repeat, delay
	err := ioctl(dev.File.Fd(), uintptr(EVIOCSREP), unsafe.Pointer(repeat_delay))
	if err != 0 {
		log.Printf("SetRepeatRate().ioctl error: %v", err)
	}
}

// Grab the input device exclusively.
func (dev *InputDevice) Grab() error {
	grab := int(1)
	if err := ioctl(dev.File.Fd(), uintptr(EVIOCGRAB), unsafe.Pointer(&grab)); err != 0 {
		log.Printf("Grab().ioctl error: %v", err)
		return err
	}

	return nil
}

// Release a grabbed input device.
func (dev *InputDevice) Release() error {
	if err := ioctl(dev.File.Fd(), uintptr(EVIOCGRAB), unsafe.Pointer(nil)); err != 0 {
		log.Printf("Release().ioctl error: %v", err)
		return err
	}

	return nil
}

type CapabilityType struct {
	Type int
	Name string
}

type CapabilityCode struct {
	Code int
	Name string
}

type AbsInfo struct {
	value      int32
	minimum    int32
	maximum    int32
	fuzz       int32
	flat       int32
	resolution int32
}

// Corresponds to the input_id struct.
type device_info struct {
	bustype, vendor, product, version uint16
}

// Return the keys of a map as a slice (dict.keys())
func keys(cap *map[int][]int) []int {
	slice := make([]int, 0)

	for key := range *cap {
		slice = append(slice, key)
	}

	return slice
}

// Determine if a path exist and is a character input device.
func IsInputDevice(path string) bool {
	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	m := fi.Mode()
	if m&os.ModeCharDevice == 0 {
		return false
	}

	return true
}

// Return a list of accessible input device names matched by
// deviceglob (default '/dev/input/event*').
func ListInputDevicePaths(device_glob string) ([]string, error) {
	paths, err := filepath.Glob(device_glob)

	if err != nil {
		return nil, err
	}

	devices := make([]string, 0)
	for _, path := range paths {
		if IsInputDevice(path) {
			devices = append(devices, path)
		}
	}

	return devices, nil
}

// Return a list of accessible input devices matched by deviceglob
// (default '/dev/input/event/*').
func ListInputDevices(device_glob_arg ...string) ([]*InputDevice, error) {
	device_glob := "/dev/input/event*"
	if len(device_glob_arg) > 0 {
		device_glob = device_glob_arg[0]
	}

	fns, _ := ListInputDevicePaths(device_glob)
	devices := make([]*InputDevice, 0)

	for i := range fns {
		dev, err := Open(fns[i])
		if err == nil {
			devices = append(devices, dev)
		}
	}

	return devices, nil
}

func bytesToString(b *[MAX_NAME_SIZE]byte) string {
	idx := bytes.IndexByte(b[:], 0)
	return string(b[:idx])
}
