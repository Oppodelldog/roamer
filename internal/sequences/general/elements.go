package general

import (
	"fmt"
	key2 "github.com/Oppodelldog/roamer/internal/key"
	mouse2 "github.com/Oppodelldog/roamer/internal/mouse"
)

type KeyDown struct {
	Key int
}

func (e KeyDown) Do() error {
	fmt.Println("down ", e.Key)
	return key2.Down(e.Key)
}

type KeyUp struct {
	Key int
}

func (e KeyUp) Do() error {
	fmt.Println("up ", e.Key)
	return key2.Up(e.Key)
}

type LeftMouseButtonDown struct {
}

func (e LeftMouseButtonDown) Do() error {
	fmt.Println("lmb-down")
	return mouse2.LeftDown()
}

type RightMouseButtonDown struct {
}

func (e RightMouseButtonDown) Do() error {
	fmt.Println("rmb-down")
	return mouse2.RightDown()
}

type SetMousePos struct {
	Pos mouse2.Pos
}

func (e SetMousePos) Do() error {
	err := mouse2.SetPosition(e.Pos)
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
	return mouse2.LeftUp()
}

type RightMouseButtonUp struct {
}

func (e RightMouseButtonUp) Do() error {
	fmt.Println("rmb-up")
	return mouse2.RightUp()
}

type LookupMousePos struct {
}

func (e LookupMousePos) Do() error {
	pos := mouse2.GetCursorPos()
	fmt.Printf("current mouse pos: %#v\n", pos)

	return nil
}

type MouseMove struct {
	X int32
	Y int32
}

func (e MouseMove) Do() error {
	fmt.Println("move ", e.X, e.Y)
	err := mouse2.Move(e.X, e.Y)

	return err
}
