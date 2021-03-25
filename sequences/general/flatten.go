package general

import "rust-roamer/sequencer"

func Flatten(seqList [][]sequencer.Elem) []sequencer.Elem {
	var s []sequencer.Elem

	for _, elements := range seqList {
		s = append(s, elements...)
	}

	return s
}
