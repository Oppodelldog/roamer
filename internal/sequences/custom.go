package sequences

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
	"log"
	"strconv"
	"strings"
	"time"
)

var keywordMappings = map[string]func() interface{}{
	"W":  func() interface{} { return sequencer.Wait{} },
	"L":  func() interface{} { return sequencer.Loop{} },
	"KD": func() interface{} { return general.KeyDown{} },
	"KU": func() interface{} { return general.KeyUp{} },
	"LD": func() interface{} { return general.LeftMouseButtonDown{} },
	"LU": func() interface{} { return general.LeftMouseButtonUp{} },
	"RD": func() interface{} { return general.RightMouseButtonDown{} },
	"RU": func() interface{} { return general.RightMouseButtonUp{} },
	"MM": func() interface{} { return general.MouseMove{} },
	"SM": func() interface{} { return general.SetMousePos{} },
}

func NewCustomSequenceFunc(script string) func() []sequencer.Elem {
	return func() []sequencer.Elem {
		var seq []sequencer.Elem
		var elements = strings.Split(script, ";")

		for _, elem := range elements {
			elem = strings.Trim(elem, " ")
			var parts = strings.Split(elem, " ")
			keyword := parts[0]
			var t, ok = keywordMappings[keyword]
			if !ok {
				panic(fmt.Sprintf("keyword '%s' not found", keyword))
			}

			var (
				arguments = parts[1:]
				err       error
			)
			elem := t().(sequencer.Elem)
			elem, err = parameterize(elem, arguments)
			if err != nil {
				log.Printf("error while parsing custom sequence: %v", err)
			}

			seq = append(seq, elem)
		}

		return seq
	}
}

func parameterize(elem sequencer.Elem, arguments []string) (sequencer.Elem, error) {
	var err error

	switch v := elem.(type) {
	case sequencer.Wait:
		var duration time.Duration
		duration, err = time.ParseDuration(arguments[0])
		v.Duration = duration
		return v, err
	case general.KeyDown:
		var key int
		key, err = getKeyValue(arguments[0])
		v.Key = key
		return v, err
	case general.KeyUp:
		var key int
		key, err = getKeyValue(arguments[0])
		v.Key = key
		return v, err
	case general.LeftMouseButtonDown:
		return v, err
	case general.LeftMouseButtonUp:
		return v, err
	case general.RightMouseButtonDown:
		return v, err
	case general.RightMouseButtonUp:
		return v, err
	case general.MouseMove:
		v.X, v.Y, err = argInt32Coordinates(arguments)
		return v, err
	case general.SetMousePos:
		var x, y int32
		x, y, err = argInt32Coordinates(arguments)
		v.Pos = mouse.Pos{
			X: x,
			Y: y,
		}
		return v, err
	case sequencer.Loop:
		return v, err
	}

	return elem, fmt.Errorf("unkown keyword %t", elem)
}

func argInt32Coordinates(arguments []string) (int32, int32, error) {
	x, err := strconv.ParseInt(arguments[0], 0, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("arg1 must be int: %w", err)
	}

	y, err := strconv.Atoi(arguments[0])
	if err != nil {
		return 0, 0, fmt.Errorf("arg2 must be int: %w", err)
	}

	return int32(x), int32(y), nil
}

func getKeyValue(s string) (int, error) {
	value, ok := map[string]int{
		"ESC": key.VK_ESC,
		"1":   key.VK_1,
		"2":   key.VK_2,
		"3":   key.VK_3,
		"4":   key.VK_4,
		"5":   key.VK_5,
		"6":   key.VK_6,
		"7":   key.VK_7,
		"8":   key.VK_8,
		"9":   key.VK_9,
		"0":   key.VK_0,
		"Q":   key.VK_Q,
		"W":   key.VK_W,
		"E":   key.VK_E,
		"R":   key.VK_R,
		"T":   key.VK_T,
		"Y":   key.VK_Y,
		"U":   key.VK_U,
		"I":   key.VK_I,
		"O":   key.VK_O,
		"P":   key.VK_P,
		"A":   key.VK_A,
		"S":   key.VK_S,
		"D":   key.VK_D,
		"F":   key.VK_F,
		"G":   key.VK_G,
		"H":   key.VK_H,
		"J":   key.VK_J,
		"K":   key.VK_K,
		"L":   key.VK_L,
		"Z":   key.VK_Z,
		"X":   key.VK_X,
		"C":   key.VK_C,
		"V":   key.VK_V,
		"B":   key.VK_B,
		"N":   key.VK_N,
		"M":   key.VK_M,
		"F1":  key.VK_F1,
		"F2":  key.VK_F2,
		"F3":  key.VK_F3,
		"F4":  key.VK_F4,
		"F5":  key.VK_F5,
		"F6":  key.VK_F6,
		"F7":  key.VK_F7,
		"F8":  key.VK_F8,
		"F9":  key.VK_F9,
		"F10": key.VK_F10,
		"F11": key.VK_F11,
		"F12": key.VK_F12,

		"NUMLOCK":    key.VK_NUMLOCK,
		"SCROLLLOCK": key.VK_SCROLLLOCK,
		"RESERVED":   key.VK_RESERVED,
		"MINUS":      key.VK_MINUS,
		"EQUAL":      key.VK_EQUAL,
		"BACKSPACE":  key.VK_BACKSPACE,
		"TAB":        key.VK_TAB,
		"LEFTBRACE":  key.VK_LEFTBRACE,
		"RIGHTBRACE": key.VK_RIGHTBRACE,
		"ENTER":      key.VK_ENTER,
		"SEMICOLON":  key.VK_SEMICOLON,
		"APOSTROPHE": key.VK_APOSTROPHE,
		"GRAVE":      key.VK_GRAVE,
		"BACKSLASH":  key.VK_BACKSLASH,
		"COMMA":      key.VK_COMMA,
		"DOT":        key.VK_DOT,
		"SLASH":      key.VK_SLASH,
		"KPASTERISK": key.VK_KPASTERISK,
		"SPACE":      key.VK_SPACE,
		"CAPSLOCK":   key.VK_CAPSLOCK,

		"LBUTTON":    key.VK_LBUTTON,
		"RBUTTON":    key.VK_RBUTTON,
		"CANCEL":     key.VK_CANCEL,
		"MBUTTON":    key.VK_MBUTTON,
		"XBUTTON1":   key.VK_XBUTTON1,
		"XBUTTON2":   key.VK_XBUTTON2,
		"BACK":       key.VK_BACK,
		"CLEAR":      key.VK_CLEAR,
		"PAUSE":      key.VK_PAUSE,
		"CAPITAL":    key.VK_CAPITAL,
		"KANA":       key.VK_KANA,
		"HANGUEL":    key.VK_HANGUEL,
		"HANGUL":     key.VK_HANGUL,
		"JUNJA":      key.VK_JUNJA,
		"FINAL":      key.VK_FINAL,
		"HANJA":      key.VK_HANJA,
		"KANJI":      key.VK_KANJI,
		"CONVERT":    key.VK_CONVERT,
		"NONCONVERT": key.VK_NONCONVERT,
		"ACCEPT":     key.VK_ACCEPT,
		"MODECHANGE": key.VK_MODECHANGE,
		"PAGEUP":     key.VK_PAGEUP,
		"PAGEDOWN":   key.VK_PAGEDOWN,
		"END":        key.VK_END,
		"HOME":       key.VK_HOME,
		"LEFT":       key.VK_LEFT,
		"UP":         key.VK_UP,
		"RIGHT":      key.VK_RIGHT,
		"DOWN":       key.VK_DOWN,
		"SELECT":     key.VK_SELECT,
		"PRINT":      key.VK_PRINT,
		"EXECUTE":    key.VK_EXECUTE,
		"SNAPSHOT":   key.VK_SNAPSHOT,
		"INSERT":     key.VK_INSERT,
		"DELETE":     key.VK_DELETE,
		"HELP":       key.VK_HELP,

		"SCROLL":              key.VK_SCROLL,
		"LMENU":               key.VK_LMENU,
		"RMENU":               key.VK_RMENU,
		"BROWSER_BACK":        key.VK_BROWSER_BACK,
		"BROWSER_FORWARD":     key.VK_BROWSER_FORWARD,
		"BROWSER_REFRESH":     key.VK_BROWSER_REFRESH,
		"BROWSER_STOP":        key.VK_BROWSER_STOP,
		"BROWSER_SEARCH":      key.VK_BROWSER_SEARCH,
		"BROWSER_FAVORITES":   key.VK_BROWSER_FAVORITES,
		"BROWSER_HOME":        key.VK_BROWSER_HOME,
		"VOLUME_MUTE":         key.VK_VOLUME_MUTE,
		"VOLUME_DOWN":         key.VK_VOLUME_DOWN,
		"VOLUME_UP":           key.VK_VOLUME_UP,
		"MEDIA_NEXT_TRACK":    key.VK_MEDIA_NEXT_TRACK,
		"MEDIA_PREV_TRACK":    key.VK_MEDIA_PREV_TRACK,
		"MEDIA_STOP":          key.VK_MEDIA_STOP,
		"MEDIA_PLAY_PAUSE":    key.VK_MEDIA_PLAY_PAUSE,
		"LAUNCH_MAIL":         key.VK_LAUNCH_MAIL,
		"LAUNCH_MEDIA_SELECT": key.VK_LAUNCH_MEDIA_SELECT,
		"LAUNCH_APP1":         key.VK_LAUNCH_APP1,
		"LAUNCH_APP2":         key.VK_LAUNCH_APP2,
		"OEM_1":               key.VK_OEM_1,
		"OEM_PLUS":            key.VK_OEM_PLUS,
		"OEM_COMMA":           key.VK_OEM_COMMA,
		"OEM_MINUS":           key.VK_OEM_MINUS,
		"OEM_PERIOD":          key.VK_OEM_PERIOD,
		"OEM_2":               key.VK_OEM_2,
		"OEM_3":               key.VK_OEM_3,
		"OEM_4":               key.VK_OEM_4,
		"OEM_5":               key.VK_OEM_5,
		"OEM_6":               key.VK_OEM_6,
		"OEM_7":               key.VK_OEM_7,
		"OEM_8":               key.VK_OEM_8,
		"OEM_102":             key.VK_OEM_102,
		"PROCESSKEY":          key.VK_PROCESSKEY,
		"PACKET":              key.VK_PACKET,
		"ATTN":                key.VK_ATTN,
		"CRSEL":               key.VK_CRSEL,
		"EXSEL":               key.VK_EXSEL,
		"EREOF":               key.VK_EREOF,
		"PLAY":                key.VK_PLAY,
		"ZOOM":                key.VK_ZOOM,
		"NONAME":              key.VK_NONAME,
		"PA1":                 key.VK_PA1,
		"OEM_CLEAR":           key.VK_OEM_CLEAR,

		"KP0":     key.VK_KP0,
		"KP1":     key.VK_KP1,
		"KP2":     key.VK_KP2,
		"KP3":     key.VK_KP3,
		"KP4":     key.VK_KP4,
		"KP5":     key.VK_KP5,
		"KP6":     key.VK_KP6,
		"KP7":     key.VK_KP7,
		"KP8":     key.VK_KP8,
		"KP9":     key.VK_KP9,
		"KPMINUS": key.VK_KPMINUS,
		"KPPLUS":  key.VK_KPPLUS,
		"KPDOT":   key.VK_KPDOT,

		"SP1":  key.VK_SP1,
		"SP2":  key.VK_SP2,
		"SP3":  key.VK_SP3,
		"SP4":  key.VK_SP4,
		"SP5":  key.VK_SP5,
		"SP6":  key.VK_SP6,
		"SP7":  key.VK_SP7,
		"SP8":  key.VK_SP8,
		"SP9":  key.VK_SP9,
		"SP10": key.VK_SP10,
		"SP11": key.VK_SP11,
		"SP12": key.VK_SP12,

		"SHIFT":    key.VK_SHIFT,
		"CTRL":     key.VK_CTRL,
		"ALT":      key.VK_ALT,
		"LSHIFT":   key.VK_LSHIFT,
		"RSHIFT":   key.VK_RSHIFT,
		"LCONTROL": key.VK_LCONTROL,
		"CONTROL":  key.K_RCONTROL,
		"LWIN":     key.VK_LWIN,
		"RWIN":     key.VK_RWIN,
	}[s]

	if !ok {
		return 0, fmt.Errorf("unknown key '%s'", s)
	}

	return value, nil
}
