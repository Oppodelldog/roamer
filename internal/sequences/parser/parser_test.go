package parser_test

import (
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
	"github.com/Oppodelldog/roamer/internal/sequences/parser"
	"reflect"
	"testing"
	"time"
)

func TestNewCustomSequenceFunc(t *testing.T) {
	var (
		script = "KD D;R 1 [LD;W 60ms;LU;R 3 [W 800ms]];;;KU D;"
		got    = parser.NewCustomSequenceFunc(script)()
		want   = []sequencer.Elem{
			general.KeyDown{Key: key.VK_D},
			general.LeftMouseButtonDown{},
			sequencer.Wait{Duration: time.Millisecond * 60},
			general.LeftMouseButtonUp{},
			sequencer.Wait{Duration: time.Millisecond * 800},
			sequencer.Wait{Duration: time.Millisecond * 800},
			sequencer.Wait{Duration: time.Millisecond * 800},
			general.KeyUp{Key: key.VK_D},
		}
	)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("sequences did not match:\ngot : %#v\nwant: %#v\n", got, want)
	}
}
