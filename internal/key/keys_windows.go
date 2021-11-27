package key

import (
	"fmt"
	"strings"
	"syscall"
)

const (
	VkShift           = 0x10 + 0xFFF
	VkCtrl            = 0x11 + 0xFFF
	VkAlt             = 0x12 + 0xFFF
	VkLshift          = 0xA0 + 0xFFF
	VkRshift          = 0xA1 + 0xFFF
	VkLcontrol        = 0xA2 + 0xFFF
	KRcontrol         = 0xA3 + 0xFFF
	VkLwin            = 0x5B + 0xFFF
	VkRwin            = 0x5C + 0xFFF
	keyEventFKeyUp    = 0x0002
	keyeventfScancode = 0x0008
)

const (
	VkSp1  = 41
	VkSp2  = 12
	VkSp3  = 13
	VkSp4  = 26
	VkSp5  = 27
	VkSp6  = 39
	VkSp7  = 40
	VkSp8  = 43
	VkSp9  = 51
	VkSp10 = 52
	VkSp11 = 53
	VkSp12 = 86

	VkEsc = 1
	Vk1   = 2
	Vk2   = 3
	Vk3   = 4
	Vk4   = 5
	Vk5   = 6
	Vk6   = 7
	Vk7   = 8
	Vk8   = 9
	Vk9   = 10
	Vk0   = 11
	VkQ   = 16
	VkW   = 17
	VkE   = 18
	VkR   = 19
	VkT   = 20
	VkY   = 21
	VkU   = 22
	VkI   = 23
	VkO   = 24
	VkP   = 25
	VkA   = 30
	VkS   = 31
	VkD   = 32
	VkF   = 33
	VkG   = 34
	VkH   = 35
	VkJ   = 36
	VkK   = 37
	VkL   = 38
	VkZ   = 44
	VkX   = 45
	VkC   = 46
	VkV   = 47
	VkB   = 48
	VkN   = 49
	VkM   = 50
	VkF1  = 59
	VkF2  = 60
	VkF3  = 61
	VkF4  = 62
	VkF5  = 63
	VkF6  = 64
	VkF7  = 65
	VkF8  = 66
	VkF9  = 67
	VkF10 = 68
	VkF11 = 87
	VkF12 = 88

	VkNumlock    = 69
	VkScrolllock = 70
	VkReserved   = 0
	VkMinus      = 12
	VkEqual      = 13
	VkBackspace  = 14
	VkTab        = 15
	VkLeftbrace  = 26
	VkRightbrace = 27
	VkEnter      = 28
	VkSemicolon  = 39
	VkApostrophe = 40
	VkGrave      = 41
	VkBackslash  = 43
	VkComma      = 51
	VkDot        = 52
	VkSlash      = 53
	VkKpasterisk = 55
	VkSpace      = 57
	VkCapslock   = 58

	VkKp0     = 82
	VkKp1     = 79
	VkKp2     = 80
	VkKp3     = 81
	VkKp4     = 75
	VkKp5     = 76
	VkKp6     = 77
	VkKp7     = 71
	VkKp8     = 72
	VkKp9     = 73
	VkKpminus = 74
	VkKpplus  = 78
	VkKpdot   = 83

	// add 0xFFF for all Virtual key

	VkLButton    = 0x01 + 0xFFF
	VkRButton    = 0x02 + 0xFFF
	VkCancel     = 0x03 + 0xFFF
	VkMButton    = 0x04 + 0xFFF
	VkXButton1   = 0x05 + 0xFFF
	VkXButton2   = 0x06 + 0xFFF
	VkBack       = 0x08 + 0xFFF
	VkClear      = 0x0C + 0xFFF
	VkPause      = 0x13 + 0xFFF
	VkCapital    = 0x14 + 0xFFF
	VkKana       = 0x15 + 0xFFF
	VkHanguel    = 0x15 + 0xFFF
	VkHangul     = 0x15 + 0xFFF
	VkJunja      = 0x17 + 0xFFF
	VkFinal      = 0x18 + 0xFFF
	VkHanja      = 0x19 + 0xFFF
	VkKanji      = 0x19 + 0xFFF
	VkConvert    = 0x1C + 0xFFF
	VkNonconvert = 0x1D + 0xFFF
	VkAccept     = 0x1E + 0xFFF
	VkModechange = 0x1F + 0xFFF
	VkPageup     = 0x21 + 0xFFF
	VkPagedown   = 0x22 + 0xFFF
	VkEnd        = 0x23 + 0xFFF
	VkHome       = 0x24 + 0xFFF
	VkLeft       = 0x25 + 0xFFF
	VkUp         = 0x26 + 0xFFF
	VkRight      = 0x27 + 0xFFF
	VkDown       = 0x28 + 0xFFF
	VkSelect     = 0x29 + 0xFFF
	VkPrint      = 0x2A + 0xFFF
	VkExecute    = 0x2B + 0xFFF
	VkSnapshot   = 0x2C + 0xFFF
	VkInsert     = 0x2D + 0xFFF
	VkDelete     = 0x2E + 0xFFF
	VkHelp       = 0x2F + 0xFFF

	VkScroll            = 0x91 + 0xFFF
	VkLMenu             = 0xA4 + 0xFFF
	VkRMenu             = 0xA5 + 0xFFF
	VkBrowserBack       = 0xA6 + 0xFFF
	VkBrowserForward    = 0xA7 + 0xFFF
	VkBrowserRefresh    = 0xA8 + 0xFFF
	VkBrowserStop       = 0xA9 + 0xFFF
	VkBrowserSearch     = 0xAA + 0xFFF
	VkBrowserFavorites  = 0xAB + 0xFFF
	VkBrowserHome       = 0xAC + 0xFFF
	VkVolumeMute        = 0xAD + 0xFFF
	VkVolumeDown        = 0xAE + 0xFFF
	VkVolumeUp          = 0xAF + 0xFFF
	VkMediaNextTrack    = 0xB0 + 0xFFF
	VkMediaPrevTrack    = 0xB1 + 0xFFF
	VkMediaStop         = 0xB2 + 0xFFF
	VkMediaPlayPause    = 0xB3 + 0xFFF
	VkLaunchMail        = 0xB4 + 0xFFF
	VkLaunchMediaSelect = 0xB5 + 0xFFF
	VkLaunchApp1        = 0xB6 + 0xFFF
	VkLaunchApp2        = 0xB7 + 0xFFF
	VkOem1              = 0xBA + 0xFFF
	VkOemPlus           = 0xBB + 0xFFF
	VkOemComma          = 0xBC + 0xFFF
	VkOemMinus          = 0xBD + 0xFFF
	VkOemPeriod         = 0xBE + 0xFFF
	VkOem2              = 0xBF + 0xFFF
	VkOem3              = 0xC0 + 0xFFF
	VkOem4              = 0xDB + 0xFFF
	VkOem5              = 0xDC + 0xFFF
	VkOem6              = 0xDD + 0xFFF
	VkOem7              = 0xDE + 0xFFF
	VkOem8              = 0xDF + 0xFFF
	VkOem102            = 0xE2 + 0xFFF
	VkProcessKey        = 0xE5 + 0xFFF
	VkPacket            = 0xE7 + 0xFFF
	VkAttn              = 0xF6 + 0xFFF
	VkCrSel             = 0xF7 + 0xFFF
	VkExSel             = 0xF8 + 0xFFF
	VkErEof             = 0xF9 + 0xFFF
	VkPlay              = 0xFA + 0xFFF
	VkZoom              = 0xFB + 0xFFF
	VkNoname            = 0xFC + 0xFFF
	VkPa1               = 0xFD + 0xFFF
	VkOemClear          = 0xFE + 0xFFF
)

var dll = syscall.NewLazyDLL("user32.dll")
var procKeyBd = dll.NewProc("keybd_event")
var states = map[int]bool{}

func Down(key int) error {
	var flag = 0

	states[key] = true

	if key < 0xFFF {
		flag |= keyeventfScancode
	} else {
		key -= 0xFFF
	}

	vkey := key + 0x80

	_, _, err := procKeyBd.Call(uintptr(key), uintptr(vkey), uintptr(flag), 0)
	if err != nil && isNotSuccess(err) {
		return fmt.Errorf("error Down: %w", err)
	}

	return nil
}

func Up(key int) error {
	states[key] = false
	flag := keyEventFKeyUp

	if key < 0xFFF {
		flag |= keyeventfScancode
	} else {
		key -= 0xFFF
	}

	vkey := key + 0x80

	_, _, err := procKeyBd.Call(uintptr(key), uintptr(vkey), uintptr(flag), 0)
	if err != nil && isNotSuccess(err) {
		return fmt.Errorf("error Up: %w", err)
	}

	return nil
}

func isNotSuccess(err error) bool {
	return !strings.Contains(err.Error(), "success")
}

func ResetPressed() {
	for key, pressed := range states {
		if pressed {
			err := Up(key)
			if err != nil && isNotSuccess(err) {
				fmt.Printf("error when ressetting pressed state for key %v: %v", key, err)
			}
		}
	}
}
