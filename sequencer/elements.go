package sequencer

import "time"

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
}

type LeftMouseButtonUp struct {
}

type Loop struct {
	Elem
}
