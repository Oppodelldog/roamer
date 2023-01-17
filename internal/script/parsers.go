package script

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
)

var ErrUnknownKey = errors.New("unknown key")
var ErrExpectedLiteral = errors.New("expected literal")
var ErrExpectedArgumentSeparator = errors.New("expected argument separator")
var ErrExpectedBlockOpen = errors.New("repeat expects arg2 to be a sequence that is introduces with '" + string(valBlockOpen) + "'")

func parseRepeat(v sequencer.Repeat, t *TokenStream) (sequencer.Elem, error) {
	var (
		err      error
		lit      Token
		sep      Token
		open     Token
		times    int
		sequence []sequencer.Elem
	)

	lit, err = parseLiteral(t)
	if err != nil {
		return nil, err
	}

	times, err = strconv.Atoi(lit.Value)
	if err != nil {
		return nil, err
	}

	if sepErr := consumeSeparators(t); sepErr != nil {
		return nil, fmt.Errorf("%w after arg1 in repeat, but was: %s(='%s')", ErrExpectedArgumentSeparator, sep.Type, sep.Value)
	}

	open = t.Consume()
	if open.Type != blockOpen {
		return nil, fmt.Errorf(
			"%w, but was: %s(='%s')",
			ErrExpectedBlockOpen,
			open.Type,
			open.Value,
		)
	}

	sequence, err = parse(t)
	if err != nil {
		return nil, err
	}

	v.Times = times
	v.Sequence = sequence

	return v, nil
}

func parseLiteral(t *TokenStream) (Token, error) {

	var errSeparators = consumeSeparators(t)
	var lit = t.Consume()

	if errSeparators != nil {
		return Token{}, errSeparators
	}

	if lit.Type != literal {
		return Token{}, fmt.Errorf("W %w, but was '%s'", ErrExpectedLiteral, lit.Type)
	}

	return lit, nil
}

func consumeSeparators(t *TokenStream) error {
	var sep = t.Consume()

	if sep.Type != argumentSeparator {
		return fmt.Errorf("%w ' ', but got '%s'", ErrExpectedArgumentSeparator, sep.Type)

	}

	for !t.isEOF() && t.Peek().Type == argumentSeparator {
		t.Consume()
	}

	return nil
}

func parseWait(v sequencer.Wait, t *TokenStream) (sequencer.Elem, error) {
	var (
		lit      Token
		err      error
		duration time.Duration
	)

	lit, err = parseLiteral(t)
	if err != nil {
		return nil, err
	}

	duration, err = time.ParseDuration(lit.Value)
	if err != nil {
		return nil, err
	}

	v.Duration = duration

	return v, nil
}

func parseKey(elem sequencer.Elem, t *TokenStream) (sequencer.Elem, error) {
	var (
		lit     Token
		err     error
		keycode int
	)

	lit, err = parseLiteral(t)
	if err != nil {
		return nil, err
	}

	keycode, err = keyCodeFromString(lit.Value)
	if err != nil {
		return nil, err
	}

	switch v := elem.(type) {
	case general.KeyDown:
		v.Key = keycode

		return v, nil
	case general.KeyUp:
		v.Key = keycode

		return v, nil
	}

	panic(fmt.Sprintf("invalid element %T", elem))
}

func parseMouseMove(v general.MouseMove, t *TokenStream) (sequencer.Elem, error) {
	var err error

	v.X, v.Y, err = argInt32Coordinates(t)

	return v, err
}

func parseSetMousePos(v general.SetMousePos, t *TokenStream) (elem sequencer.Elem, err error) {
	var pos mouse.Pos

	pos.X, pos.Y, err = argInt32Coordinates(t)
	v.Pos = pos

	return v, err
}

func argInt32Coordinates(t *TokenStream) (int32, int32, error) {
	var (
		err  error
		lit1 Token
		lit2 Token
		x    int64
		y    int64
	)

	lit1, err = parseLiteral(t)
	if err != nil {
		return 0, 0, err
	}

	lit2, err = parseLiteral(t)
	if err != nil {
		return 0, 0, err
	}

	if lit1.Type != literal {
		return 0, 0, fmt.Errorf("arg1 %w, but was '%s'", ErrExpectedLiteral, lit1.Type)
	}

	if lit2.Type != literal {
		return 0, 0, fmt.Errorf("arg2 %w, but was '%s'", ErrExpectedLiteral, lit2.Type)
	}

	x, err = strconv.ParseInt(lit1.Value, 0, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("arg1 must be int: %w", err)
	}

	y, err = strconv.ParseInt(lit2.Value, 0, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("arg2 must be int: %w", err)
	}

	return int32(x), int32(y), nil
}

func keyCodeFromString(s string) (int, error) {
	value, ok := keyCodeStringMap[s]
	if !ok {
		return 0, fmt.Errorf("%w '%s'", ErrUnknownKey, s)
	}

	return value, nil
}

func stringFromKeyCode(v int) string {
	for s, code := range keyCodeStringMap {
		if code == v {
			return s
		}
	}

	return ""
}

var keyCodeStringMap = map[string]int{
	"ESC": key.VkEsc,
	"1":   key.Vk1,
	"2":   key.Vk2,
	"3":   key.Vk3,
	"4":   key.Vk4,
	"5":   key.Vk5,
	"6":   key.Vk6,
	"7":   key.Vk7,
	"8":   key.Vk8,
	"9":   key.Vk9,
	"0":   key.Vk0,
	"Q":   key.VkQ,
	"W":   key.VkW,
	"E":   key.VkE,
	"R":   key.VkR,
	"T":   key.VkT,
	"Y":   key.VkY,
	"U":   key.VkU,
	"I":   key.VkI,
	"O":   key.VkO,
	"P":   key.VkP,
	"A":   key.VkA,
	"S":   key.VkS,
	"D":   key.VkD,
	"F":   key.VkF,
	"G":   key.VkG,
	"H":   key.VkH,
	"J":   key.VkJ,
	"K":   key.VkK,
	"L":   key.VkL,
	"Z":   key.VkZ,
	"X":   key.VkX,
	"C":   key.VkC,
	"V":   key.VkV,
	"B":   key.VkB,
	"N":   key.VkN,
	"M":   key.VkM,
	"F1":  key.VkF1,
	"F2":  key.VkF2,
	"F3":  key.VkF3,
	"F4":  key.VkF4,
	"F5":  key.VkF5,
	"F6":  key.VkF6,
	"F7":  key.VkF7,
	"F8":  key.VkF8,
	"F9":  key.VkF9,
	"F10": key.VkF10,
	"F11": key.VkF11,
	"F12": key.VkF12,

	"NUMLOCK":    key.VkNumlock,
	"SCROLLLOCK": key.VkScrolllock,
	"RESERVED":   key.VkReserved,
	"MINUS":      key.VkMinus,
	"EQUAL":      key.VkEqual,
	"BACKSPACE":  key.VkBackspace,
	"TAB":        key.VkTab,
	"LEFTBRACE":  key.VkLeftbrace,
	"RIGHTBRACE": key.VkRightbrace,
	"ENTER":      key.VkEnter,
	"SEMICOLON":  key.VkSemicolon,
	"APOSTROPHE": key.VkApostrophe,
	"GRAVE":      key.VkGrave,
	"BACKSLASH":  key.VkBackslash,
	"COMMA":      key.VkComma,
	"DOT":        key.VkDot,
	"SLASH":      key.VkSlash,
	"KPASTERISK": key.VkKpasterisk,
	"SPACE":      key.VkSpace,
	"CAPSLOCK":   key.VkCapslock,

	"LBUTTON":    key.VkLButton,
	"RBUTTON":    key.VkRButton,
	"CANCEL":     key.VkCancel,
	"MBUTTON":    key.VkMButton,
	"XBUTTON1":   key.VkXButton1,
	"XBUTTON2":   key.VkXButton2,
	"BACK":       key.VkBack,
	"CLEAR":      key.VkClear,
	"PAUSE":      key.VkPause,
	"CAPITAL":    key.VkCapital,
	"KANA":       key.VkKana,
	"HANGUEL":    key.VkHanguel,
	"HANGUL":     key.VkHangul,
	"JUNJA":      key.VkJunja,
	"FINAL":      key.VkFinal,
	"HANJA":      key.VkHanja,
	"KANJI":      key.VkKanji,
	"CONVERT":    key.VkConvert,
	"NONCONVERT": key.VkNonconvert,
	"ACCEPT":     key.VkAccept,
	"MODECHANGE": key.VkModechange,
	"PAGEUP":     key.VkPageup,
	"PAGEDOWN":   key.VkPagedown,
	"END":        key.VkEnd,
	"HOME":       key.VkHome,
	"LEFT":       key.VkLeft,
	"UP":         key.VkUp,
	"RIGHT":      key.VkRight,
	"DOWN":       key.VkDown,
	"SELECT":     key.VkSelect,
	"PRINT":      key.VkPrint,
	"EXECUTE":    key.VkExecute,
	"SNAPSHOT":   key.VkSnapshot,
	"INSERT":     key.VkInsert,
	"DELETE":     key.VkDelete,
	"HELP":       key.VkHelp,

	"SCROLL":              key.VkScroll,
	"LMENU":               key.VkLMenu,
	"RMENU":               key.VkRMenu,
	"BROWSER_BACK":        key.VkBrowserBack,
	"BROWSER_FORWARD":     key.VkBrowserForward,
	"BROWSER_REFRESH":     key.VkBrowserRefresh,
	"BROWSER_STOP":        key.VkBrowserStop,
	"BROWSER_SEARCH":      key.VkBrowserSearch,
	"BROWSER_FAVORITES":   key.VkBrowserFavorites,
	"BROWSER_HOME":        key.VkBrowserHome,
	"VOLUME_MUTE":         key.VkVolumeMute,
	"VOLUME_DOWN":         key.VkVolumeDown,
	"VOLUME_UP":           key.VkVolumeUp,
	"MEDIA_NEXT_TRACK":    key.VkMediaNextTrack,
	"MEDIA_PREV_TRACK":    key.VkMediaPrevTrack,
	"MEDIA_STOP":          key.VkMediaStop,
	"MEDIA_PLAY_PAUSE":    key.VkMediaPlayPause,
	"LAUNCH_MAIL":         key.VkLaunchMail,
	"LAUNCH_MEDIA_SELECT": key.VkLaunchMediaSelect,
	"LAUNCH_APP1":         key.VkLaunchApp1,
	"LAUNCH_APP2":         key.VkLaunchApp2,
	"OEM_1":               key.VkOem1,
	"OEM_PLUS":            key.VkOemPlus,
	"OEM_COMMA":           key.VkOemComma,
	"OEM_MINUS":           key.VkOemMinus,
	"OEM_PERIOD":          key.VkOemPeriod,
	"OEM_2":               key.VkOem2,
	"OEM_3":               key.VkOem3,
	"OEM_4":               key.VkOem4,
	"OEM_5":               key.VkOem5,
	"OEM_6":               key.VkOem6,
	"OEM_7":               key.VkOem7,
	"OEM_8":               key.VkOem8,
	"OEM_102":             key.VkOem102,
	"PROCESSKEY":          key.VkProcessKey,
	"PACKET":              key.VkPacket,
	"ATTN":                key.VkAttn,
	"CRSEL":               key.VkCrSel,
	"EXSEL":               key.VkExSel,
	"EREOF":               key.VkErEof,
	"PLAY":                key.VkPlay,
	"ZOOM":                key.VkZoom,
	"NONAME":              key.VkNoname,
	"PA1":                 key.VkPa1,
	"OEM_CLEAR":           key.VkOemClear,

	"KP0":     key.VkKp0,
	"KP1":     key.VkKp1,
	"KP2":     key.VkKp2,
	"KP3":     key.VkKp3,
	"KP4":     key.VkKp4,
	"KP5":     key.VkKp5,
	"KP6":     key.VkKp6,
	"KP7":     key.VkKp7,
	"KP8":     key.VkKp8,
	"KP9":     key.VkKp9,
	"KPMINUS": key.VkKpminus,
	"KPPLUS":  key.VkKpplus,
	"KPDOT":   key.VkKpdot,

	"SP1":  key.VkSp1,
	"SP2":  key.VkSp2,
	"SP3":  key.VkSp3,
	"SP4":  key.VkSp4,
	"SP5":  key.VkSp5,
	"SP6":  key.VkSp6,
	"SP7":  key.VkSp7,
	"SP8":  key.VkSp8,
	"SP9":  key.VkSp9,
	"SP10": key.VkSp10,
	"SP11": key.VkSp11,
	"SP12": key.VkSp12,

	"SHIFT":    key.VkShift,
	"CTRL":     key.VkCtrl,
	"ALT":      key.VkAlt,
	"LSHIFT":   key.VkLshift,
	"RSHIFT":   key.VkRshift,
	"LCONTROL": key.VkLcontrol,
	"CONTROL":  key.KRcontrol,
	"LWIN":     key.VkLwin,
	"RWIN":     key.VkRwin,
}
