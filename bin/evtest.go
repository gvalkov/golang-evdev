// +build linux
// evdev example - input device event monitor
package main

import (
	"errors"
	"evdev"
	"fmt"
	"os"
	"strings"
)

const (
	usage       = "usage: evtest <device> [<type> <value>]"
	device_glob = "/dev/input/event*"
)

// Select a device from a list of accessible input devices.
func select_device() (*evdev.InputDevice, error) {
	fns, _ := evdev.ListInputDevices(device_glob)
	devices := make([]*evdev.InputDevice, 0)

	for i := range fns {
		dev, err := evdev.Open(fns[i])
		if err == nil {
			devices = append(devices, dev)
		}
	}

	lines := make([]string, 0)
	max := 0
	if len(devices) > 0 {
		for i := range devices {
			dev := devices[i]
			str := fmt.Sprintf("%-3d %-20s %-35s %s", i, dev.Fn, dev.Name, dev.Phys)
			if len(str) > max {
				max = len(str)
			}
			lines = append(lines, str)
		}
		fmt.Printf("%-3s %-20s %-35s %s\n", "ID", "Device", "Name", "Phys")
		fmt.Printf(strings.Repeat("-", max) + "\n")
		fmt.Printf(strings.Join(lines, "\n") + "\n")

		var choice int
		choice_max := len(lines) - 1

	ReadChoice:
		fmt.Printf("Select device [0-%d]: ", choice_max)
		_, err := fmt.Scan(&choice)
		if err != nil || choice > choice_max || choice < 0 {
			goto ReadChoice
		}

		return devices[choice], nil
	}

	errmsg := fmt.Sprintf("no accessible input devices matched by %s", device_glob)
	return nil, errors.New(errmsg)
}

func format_event(ev *evdev.InputEvent) string {
	var res, f, code_name string

	code := int(ev.Code)
	etype := int(ev.Type)

	switch ev.Type {
	case evdev.EV_SYN:
		if ev.Code == evdev.SYN_MT_REPORT {
			f = "time %d.%-8d +++++++++ %s ++++++++"
		} else {
			f = "time %d.%-8d --------- %s --------"
		}
		return fmt.Sprintf(f, ev.Time.Sec, ev.Time.Usec, evdev.SYN[code])
	case evdev.EV_KEY:
		val, haskey := evdev.KEY[code]
		if haskey {
			code_name = val
		} else {
			val, haskey := evdev.BTN[code]
			if haskey {
				code_name = val
			} else {
				code_name = "?"
			}
		}
	default:
		m, haskey := evdev.ByEventType[etype]
		if haskey {
			code_name = m[code]
		} else {
			code_name = "?"
		}
	}

	evfmt := "time %d.%-8d type %d (%s), code %-3d (%s), value %d"
	res = fmt.Sprintf(evfmt, ev.Time.Sec, ev.Time.Usec, etype,
		evdev.EV[int(ev.Type)], ev.Code, code_name, ev.Value)

	return res
}

func main() {
	var dev *evdev.InputDevice
	var events []evdev.InputEvent
	var err error

	if len(os.Args) == 1 {
		dev, err = select_device()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if len(os.Args) == 2 {
		fmt.Printf("Selecting Device 2\n")
	} else if len(os.Args) == 4 {
		fmt.Printf("Selecting Device 4\n")
	} else {
		fmt.Printf(usage + "\n")
	}

	info := fmt.Sprintf("bus 0x%04x, vendor 0x%04x, product 0x%04x, version 0x%04x",
		dev.Bustype, dev.Vendor, dev.Product, dev.Version)

	fmt.Printf("Device name: %s\n", dev.Name)
	fmt.Printf("Device info: %s\n", info)
	fmt.Printf("Repeat settings: \n")
	fmt.Printf("Device capabilities: \n")

	fmt.Printf("Listening for events ...\n")

	for {
		events = dev.Read()
		for i := range events {
			str := format_event(&events[i])
			fmt.Println(str)
		}
	}
}



// def print_event(e):
//     if e.type == ecodes.EV_SYN:
//         if e.code == ecodes.SYN_MT_REPORT:
//             print('time {:<16} +++++++++ {} ++++++++'.format(e.timestamp(), ecodes.SYN[e.code]))
//         else:
//             print('time {:<16} --------- {} --------'.format(e.timestamp(), ecodes.SYN[e.code]))
//     else:
//         if e.type in ecodes.bytype:
//             codename = ecodes.bytype[e.type][e.code]
//         else:
//             codename = '?'
// 
//         print(evfmt.format(e.timestamp(), e.type, ecodes.EV[e.type], e.code, codename, e.value))
// 
// 
// print('Device name: {.name}'.format(device))
// print('Device info: {.info}'.format(device))
// print('Repeat settings: {}'.format(device.repeat))
// 
// print('Device capabilities:')
// for type, codes in device.capabilities(verbose=True).items():
//     print('  Type {} {}:'.format(*type))
//     for i in codes:
//         if isinstance(i[1], AbsInfo):
//             print('    Code {:<4} {}:'.format(*i[0]))
//             print('      {}'.format(i[1]))
//         else:
//             print('    Code {:<4} {}'.format(*i))
//     print('')
// 
// 
