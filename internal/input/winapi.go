package input

import (
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
)

type WinAPIExecutor struct{}

func (WinAPIExecutor) KeyDown(k int) error {
	return key.Down(k)
}

func (WinAPIExecutor) KeyUp(k int) error {
	return key.Up(k)
}

func (WinAPIExecutor) LeftDown() error {
	return mouse.LeftDown()
}

func (WinAPIExecutor) LeftUp() error {
	return mouse.LeftUp()
}

func (WinAPIExecutor) RightDown() error {
	return mouse.RightDown()
}

func (WinAPIExecutor) RightUp() error {
	return mouse.RightUp()
}

func (WinAPIExecutor) SetPosition(pos mouse.Pos) error {
	return mouse.SetPosition(pos)
}

func (WinAPIExecutor) GetCursorPos() (mouse.Pos, error) {
	return mouse.GetCursorPos()
}

func (WinAPIExecutor) Move(x, y int32) error {
	return mouse.Move(x, y)
}

func (WinAPIExecutor) ResetPressed() {
	key.ResetPressed()
	mouse.ResetPressed()
}
