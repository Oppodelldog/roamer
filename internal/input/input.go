package input

import "github.com/Oppodelldog/roamer/internal/mouse"

type Executor interface {
	KeyDown(key int) error
	KeyUp(key int) error
	LeftDown() error
	LeftUp() error
	RightDown() error
	RightUp() error
	SetPosition(pos mouse.Pos) error
	GetCursorPos() (mouse.Pos, error)
	Move(x, y int32) error
	ResetPressed()
}
