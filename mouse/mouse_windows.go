package mouse

import "syscall"

var dll = syscall.NewLazyDLL("user32.dll")

// https://docs.microsoft.com/de-de/windows/win32/api/winuser/nf-winuser-mouse_event?redirectedfrom=MSDN
var procMouseBd = dll.NewProc("mouse_event")

const (
	flagLeftDown = 0x0002
	flagLeftUp   = 0x0004
)

func LeftDown() error {
	_, _, err := procMouseBd.Call(uintptr(flagLeftDown), uintptr(0), uintptr(0), uintptr(0), 0)
	return err
}

func LeftUp() error {
	_, _, err := procMouseBd.Call(uintptr(flagLeftUp), uintptr(0), uintptr(0), uintptr(0), 0)
	return err
}
