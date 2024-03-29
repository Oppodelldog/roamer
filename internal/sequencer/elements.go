package sequencer

import "time"

type ElemFunc func() error

func (e ElemFunc) Do() error {
	return e()
}

// Elem is part of a sequence of elements. As processing sequencer calls the Do method.
type Elem interface {
	Do() error
}

// Wait lets the sequence sleep for the given amount of time.
type Wait struct {
	Elem
	Duration time.Duration
}

// Loop placed at the last element of a sequence indicated sequencer may loop the whole sequence.
type Loop struct {
	Elem
}

// Repeat repeats the embedded sequence n times.
type Repeat struct {
	Elem
	Times    int
	Sequence []Elem
}

// NoOperation does nothing.
type NoOperation struct {
	Elem
}
