package general

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
)

type KeyDown struct {
	Key int
}

func (e KeyDown) Do() error {
	fmt.Println("down ", e.Key)
	return key.Down(e.Key)
}

type KeyUp struct {
	Key int
}

func (e KeyUp) Do() error {
	fmt.Println("up ", e.Key)
	return key.Up(e.Key)
}

type LeftMouseButtonDown struct {
}

func (e LeftMouseButtonDown) Do() error {
	fmt.Println("lmb-down")
	return mouse.LeftDown()
}

type RightMouseButtonDown struct {
}

func (e RightMouseButtonDown) Do() error {
	fmt.Println("rmb-down")
	return mouse.RightDown()
}

type SetMousePos struct {
	Pos mouse.Pos
}

func (e SetMousePos) Do() error {
	err := mouse.SetPosition(e.Pos)
	if err != nil {
		fmt.Printf("error SetMousePos: %v\n", err)
	}
	fmt.Printf("set mouse pos to: %#v\n", e.Pos)

	return nil
}

type LeftMouseButtonUp struct {
}

func (e LeftMouseButtonUp) Do() error {
	fmt.Println("lmb-up")
	return mouse.LeftUp()
}

type RightMouseButtonUp struct {
}

func (e RightMouseButtonUp) Do() error {
	fmt.Println("rmb-up")
	return mouse.RightUp()
}

type LookupMousePos struct {
}

func (e LookupMousePos) Do() error {
	pos := mouse.GetCursorPos()
	fmt.Printf("current mouse pos: %#v\n", pos)

	return nil
}

type MouseMove struct {
	X int32
	Y int32
}

func (e MouseMove) Do() error {
	fmt.Println("move ", e.X, e.Y)
	err := mouse.Move(e.X, e.Y)

	return err
}
