package script

import (
	"testing"

	"github.com/Oppodelldog/roamer/internal/mouse"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
)

func TestWrite(t *testing.T) {
	var seq = []sequencer.Elem{
		general.MouseMove{X: 1, Y: 2},
		general.LookupMousePos{},
		general.RightMouseButtonDown{},
		general.RightMouseButtonUp{},
		general.LeftMouseButtonDown{},
		general.LeftMouseButtonUp{},
		general.SetMousePos{Pos: mouse.Pos{X: 3, Y: 4}},
		general.KeyUp{Key: keyCodeStringMap["TAB"]},
		general.KeyDown{Key: keyCodeStringMap["ESC"]},
		sequencer.NoOperation{},
	}

	var (
		got  = Write(seq)
		want = "MM 1 2;MP;RD;RU;LD;LU;SM 3 4;KU TAB;KD ESC;NOP"
	)

	if got != want {
		t.Fatalf("wanted: %v, but got %v", want, got)
	}
}
