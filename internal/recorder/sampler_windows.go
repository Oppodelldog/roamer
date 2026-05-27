package recorder

import (
	"syscall"

	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
)

var user32 = syscall.NewLazyDLL("user32.dll")
var procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")

type WindowsSampler struct {
	keys []key.RecordableKey
}

func NewWindowsSampler() WindowsSampler {
	return WindowsSampler{keys: key.RecordableKeys()}
}

func (s WindowsSampler) Sample() InputState {
	keys := map[int]bool{}
	for _, recordable := range s.keys {
		keys[recordable.Code] = asyncKeyDown(recordable.VirtualKey)
	}

	pos, _ := mouse.GetCursorPos()

	return InputState{
		Keys:  keys,
		Left:  asyncKeyDown(0x01),
		Right: asyncKeyDown(0x02),
		Pos:   pos,
	}
}

func asyncKeyDown(vk int) bool {
	state, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return state&0x8000 != 0
}
