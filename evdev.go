package evdev

import (
	"os"
	"fmt"
	"unsafe"
	// "strings"
)

const MAX_NAME_SIZE = 256

// input.h ioctls
var EVIOCGID   = _IOR('E', 0x02, 8)  // 8 <- sizeof(struct input_id)
var EVIOCGNAME = _IOC(_IOC_READ, 'E', 0x06, MAX_NAME_SIZE)
var EVIOCGPHYS = _IOC(_IOC_READ, 'E', 0x07, MAX_NAME_SIZE)

func EVIOCGBIT(ev, len int) int {
	return _IOC(_IOC_READ, 'E', 0x20 + ev, len)
}


type device_info struct {
	bustype, vendor, product, version uint16
}

type InputDevice struct {
	Fn   string
	Name string
	Phys string
	Fd   uintptr

	Bustype uint16
	Vendor  uint16
	Product uint16
	Version uint16

	Capabilities map[int][]int
}

func Open(devnode string) *InputDevice {
	file, err := os.Open(devnode)

	if err != nil {
		panic(err)
	}

	dev := InputDevice{}
	dev.Fn = devnode
	dev.Fd = file.Fd()

	dev.get_device_info()
	dev.get_device_capabilities()

	return &dev
}

func (dev *InputDevice) get_device_info() {
	info := device_info{}

	name := new([MAX_NAME_SIZE]byte)
	phys := new([MAX_NAME_SIZE]byte)

	// fmt.Println("size", unsafe.Sizeof(info))

	ioctl(dev.Fd, EVIOCGID, unsafe.Pointer(&info))
	ioctl(dev.Fd, EVIOCGNAME, unsafe.Pointer(name))
	ioctl(dev.Fd, EVIOCGPHYS, unsafe.Pointer(phys))

	dev.Name = string(name[:])
	dev.Phys = string(phys[:])

	dev.Vendor  = info.vendor
	dev.Bustype = info.bustype
	dev.Product = info.product
	dev.Version = info.version
}

func keys (cap *map[int][]int) []int {
	slice := make([]int, len(*cap))

	for key := range *cap {
		slice = append(slice, key)
	}

	return slice
}

func (dev *InputDevice) String() string {
	evtypes := keys(&dev.Capabilities)
	evtypes_str := make([]string, len(evtypes))

	for k := range evtypes {
		evtypes_str = append(evtypes_str, fmt.Sprintf("%s %d", EV[k], k))
		fmt.Println(evtypes_str[0])
	}
	// evtypes_str = strings.Join(evtypes_str, ", ")

	return fmt.Sprintf(
		"InputDevice %s (fd %d)\n" +
		"  name %s\n" +
		"  phys %s\n" +
		"  bus 0x%x, vendor 0x%x, product 0x%x, version 0x%x\n" +
		"  events %s",
		dev.Fn, dev.Fd, dev.Name, dev.Phys, dev.Bustype,
		dev.Vendor, dev.Product, dev.Version, evtypes_str)
}

func (dev *InputDevice) get_device_capabilities() {
	capabilities := make(map[int][]int)

	evbits   := new([ (EV_MAX+1)/8]byte)
	codebits := new([(KEY_MAX+1)/8]byte)

	ioctl(dev.Fd, EVIOCGBIT(0, EV_MAX), unsafe.Pointer(evbits))

	for evtype := 0; evtype < EV_MAX; evtype++ {
		if evbits[evtype/8] & (1 << uint(evtype % 8)) != 0 {

			eventcodes := make([]int, 0)

			ioctl(dev.Fd, EVIOCGBIT(evtype, KEY_MAX), unsafe.Pointer(codebits))
			for evcode := 0; evcode < KEY_MAX; evcode++ {
				if codebits[evcode/8] & (1 << uint(evcode % 8)) != 0 {
					eventcodes = append(eventcodes, evcode)
				}
			}

			capabilities[evtype] = eventcodes
		}
	}

	dev.Capabilities = capabilities
u
