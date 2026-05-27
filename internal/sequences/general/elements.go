package general

import (
	"fmt"
	"strings"

	"github.com/Oppodelldog/roamer/internal/input"
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
)

type KeyDown struct {
	Key int
}

func (e KeyDown) String() string {
	return fmt.Sprintf("KD %s", key.Name(e.Key))
}

func (e KeyDown) Do() error {
	fmt.Println("down ", e.Key)
	return input.Current().KeyDown(e.Key)
}

type KeyUp struct {
	Key int
}

func (e KeyUp) String() string {
	return fmt.Sprintf("KU %s", key.Name(e.Key))
}

func (e KeyUp) Do() error {
	fmt.Println("up ", e.Key)
	return input.Current().KeyUp(e.Key)
}

type LeftMouseButtonDown struct {
}

func (e LeftMouseButtonDown) String() string {
	return "LD"
}

func (e LeftMouseButtonDown) Do() error {
	fmt.Println("lmb-down")
	return input.Current().LeftDown()
}

type RightMouseButtonDown struct {
}

func (e RightMouseButtonDown) String() string {
	return "RD"
}

func (e RightMouseButtonDown) Do() error {
	fmt.Println("rmb-down")
	return input.Current().RightDown()
}

type SetMousePos struct {
	Pos mouse.Pos
}

func (e SetMousePos) String() string {
	return fmt.Sprintf("SM %d %d", e.Pos.X, e.Pos.Y)
}

func (e SetMousePos) Do() error {
	err := input.Current().SetPosition(e.Pos)
	if err != nil {
		fmt.Printf("error SetMousePos: %v\n", err)
	}

	fmt.Printf("set mouse pos to: %#v\n", e.Pos)

	return nil
}

type LeftMouseButtonUp struct {
}

func (e LeftMouseButtonUp) String() string {
	return "LU"
}

func (e LeftMouseButtonUp) Do() error {
	fmt.Println("lmb-up")
	return input.Current().LeftUp()
}

type RightMouseButtonUp struct {
}

func (e RightMouseButtonUp) String() string {
	return "RU"
}

func (e RightMouseButtonUp) Do() error {
	fmt.Println("rmb-up")
	return input.Current().RightUp()
}

type LookupMousePos struct {
}

func (e LookupMousePos) String() string {
	return "MP"
}

func (e LookupMousePos) Do() error {
	pos, err := input.Current().GetCursorPos()
	if err != nil && isNotSuccess(err) {
		return err
	}

	fmt.Printf("current mouse pos: %#v\n", pos)

	return nil
}

type MouseMove struct {
	X int32
	Y int32
}

func (e MouseMove) String() string {
	return fmt.Sprintf("MM %d %d", e.X, e.Y)
}

func (e MouseMove) Do() error {
	fmt.Println("move ", e.X, e.Y)
	err := input.Current().Move(e.X, e.Y)

	return err
}

func isNotSuccess(err error) bool {
	return !strings.Contains(err.Error(), "success")
}
