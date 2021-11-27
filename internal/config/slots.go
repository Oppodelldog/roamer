package config

import "fmt"

type Pos struct {
	X int
	Y int
}

type Slots struct {
	SlotsPerRow int
	Slots       []Pos
}

func (s Slots) Get(col, row int) Pos {
	var idx = col + row*s.SlotsPerRow
	if idx > len(s.Slots)-1 {
		panic(fmt.Sprintf("cannot get slots at (%v,%v = idx %v), number of slots is %v", col, row, idx, len(s.Slots)))
	}

	return s.Slots[idx]
}

func (s Slots) At(idx int) Pos {
	idx--
	if idx >= len(s.Slots) {
		panic(fmt.Sprintf("cannot get slots at %v, len of slots is %v", idx, len(s.Slots)))
	}

	if idx < 0 {
		panic("idx must be > 0")
	}

	return s.Slots[idx]
}
