package evdev

import "testing"

func TestAccess(t *testing.T) {
	if KEY_A != ecodes["KEY_A"] {
		t.Error()
	}

	if KEY[ecodes["KEY_A"]] != "KEY_A" {
		t.Error()
	}

	if REL[0] != "REL_X" {
		t.Error()
	}

	if ByEventType[EV_KEY][ecodes["KEY_A"]] != "KEY_A" {
		t.Error()
	}
}
