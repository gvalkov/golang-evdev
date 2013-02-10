// +build ignore

// Integer constants defined in linux/input.h can be accessed
// directly or through the provided reverse and forward mappings:
//
//   evdev.KEY_A  // 30
//   evdev.ecodes["KEY_A"]  // 30
//   evdev.KEY[30] // "KEY_A"
//   evdev.REL[0]  // "REL_X"
//   evdev.EV[evdev.EV_KEY]  // "EV_KEY"
//   evdev.ByEventType[EV_REL][0]  // "REL_X"
//
// Generated on: %(uname)s

package evdev

import "strings"

const (
%(codes)s
)

var ecodes = map[string] int {
%(codemap)s
}


// TODO figure out how to write on one line

var KEY = map[int] string {}
var ABS = map[int] string {}
var REL = map[int] string {}
var SW	= map[int] string {}
var MSC = map[int] string {}
var LED = map[int] string {}
var BTN = map[int] string {}
var REP = map[int] string {}
var SND = map[int] string {}
var ID	= map[int] string {}
var EV	= map[int] string {}
var BUS = map[int] string {}
var SYN = map[int] string {}

var ByEventType = map[int] map[int]string {
    EV_KEY: KEY,
    EV_ABS: ABS,
    EV_REL: REL,
    EV_SW:  SW,
    EV_MSC: MSC,
    EV_LED: LED,
    EV_REP: REP,
    EV_SND: SND,
    EV_SYN: SYN,
}

func init() {
	for code, value := range ecodes {
		switch {
		case strings.HasPrefix(code, "KEY"): KEY[value] = code
		case strings.HasPrefix(code, "ABS"): ABS[value] = code
		case strings.HasPrefix(code, "REL"): REL[value] = code
		case strings.HasPrefix(code, "SW"):   SW[value] = code
		case strings.HasPrefix(code, "MSC"): MSC[value] = code
		case strings.HasPrefix(code, "LED"): LED[value] = code
		case strings.HasPrefix(code, "BTN"): BTN[value] = code
		case strings.HasPrefix(code, "SND"): SND[value] = code
		case strings.HasPrefix(code, "ID"):   ID[value] = code
		case strings.HasPrefix(code, "EV"):   EV[value] = code
		case strings.HasPrefix(code, "BUS"): BUS[value] = code
		case strings.HasPrefix(code, "SYN"): SYN[value] = code
		}
	}
}
