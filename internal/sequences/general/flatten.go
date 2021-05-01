package general

import (
	sequencer2 "github.com/Oppodelldog/roamer/internal/sequencer"
)

func Flatten(seqList [][]sequencer2.Elem) []sequencer2.Elem {
	var s []sequencer2.Elem

	for _, elements := range seqList {
		s = append(s, elements...)
	}

	return s
}
