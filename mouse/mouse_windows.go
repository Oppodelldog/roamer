package mouse

import (
	"syscall"
	"unsafe"
)

var dll = syscall.NewLazyDLL("user32.dll")

// https://docs.microsoft.com/de-de/windows/win32/api/winuser/nf-winuser-mouse_event?redirectedfrom=MSDN
var procMouseBd = dll.NewProc("mouse_event")

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getcursorpos
var procGetCursorPos = dll.NewProc("GetCursorPos")

//https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setcursorpos
var procSetCursorPos = dll.NewProc("SetCursorPos")

const (
	flagLeftDown  = 0x0002
	flagLeftUp    = 0x0004
	flagRightDown = 0x0008
	flagRightUp   = 0x0010
)

func LeftDown() error {
	_, _, err := procMouseBd.Call(uintptr(flagLeftDown), uintptr(0), uintptr(0), uintptr(0), 0)
	return err
}

func LeftUp() error {
	_, _, err := procMouseBd.Call(uintptr(flagLeftUp), uintptr(0), uintptr(0), uintptr(0), 0)
	return err
}

func RightDown() error {
	_, _, err := procMouseBd.Call(uintptr(flagRightDown), uintptr(0), uintptr(0), uintptr(0), 0)
	return err
}

func RightUp() error {
	_, _, err := procMouseBd.Call(uintptr(flagRightUp), uintptr(0), uintptr(0), uintptr(0), 0)
	return err
}

func SetPosition(pos Pos) error {
	_, _, err := procSetCursorPos.Call(uintptr(pos.X), uintptr(pos.Y))
	return err
}

func GetCursorPos() Pos {
	var pos Pos
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos)))

	return pos
}

type Pos struct {
	X int32
	Y int32
}
