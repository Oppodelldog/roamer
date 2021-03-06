package sequencer

import (
	"rust-roamer/mouse"
	"time"
)

type Elem interface {
}

type Wait struct {
	Elem
	Duration time.Duration
}

type KeyDown struct {
	Elem
	Key int
}

type KeyUp struct {
	Elem
	Key int
}

type LeftMouseButtonDown struct {
	Elem
}

type RightMouseButtonDown struct {
	Elem
}

type SetMousePos struct {
	Elem
	Pos mouse.Pos
}

type LeftMouseButtonUp struct {
	Elem
}

type RightMouseButtonUp struct {
	Elem
}

type LookupMousePos struct {
	Elem
}

type Loop struct {
	Elem
}
