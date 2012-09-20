package evdev_test

import (
	"fmt"
	"evdev"
)

func ExampleOpen() {
	device := Open("/dev/input/event3")
}

func Example() {
	device := Open("/dev/input/event3")

	fmt.Println(device)
	// InputDevice /dev/input/event3 (fd 3)
	//   name Logitech USB Laser Mouse
	//   phys usb-0000:00:12.0-2/input0
	//   bus 0x0003, vendor 0x046d, product 0xc069, version 0x110
	//   events EV_KEY 1, EV_SYN 0, EV_REL 2, EV_MSC 4

	fmt.Println(device.Capabilities)
	// map[ 4:[4 272 273 274 275 276 277 278 279]
	//      0:[0 1 2 4]
	//      2:[0 1 6 8 272 273 274 275 276 277 278 279]
	//      1:[272 273 274 275 276 277 278 279] ]

	fmt.Println(device.ResolveCapabilities())
	// map[ 4:[4 272 273 274 275 276 277 278 279]
	//      0:[0 1 2 4]
	//      2:[0 1 6 8 272 273 274 275 276 277 278 279]
	//      1:[272 273 274 275 276 277 278 279] ]

}
