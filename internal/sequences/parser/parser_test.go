package parser_test

import (
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
	"github.com/Oppodelldog/roamer/internal/sequences/parser"
	"reflect"
	"testing"
	"time"
)

func TestNewCustomSequenceFunc(t *testing.T) {
	type testData struct {
		script string
		want   []sequencer.Elem
	}

	var tests = map[string]testData{
		"all commands": {
			script: "W 3ms;L;R 3 [W 4s];KD A;KU B;LD;LU;RD;RU;MM 10 20;SM 30 40",
			want: []sequencer.Elem{
				sequencer.Wait{Duration: time.Millisecond * 3},
				sequencer.Loop{},
				sequencer.Wait{Duration: 4 * time.Second},
				sequencer.Wait{Duration: 4 * time.Second},
				sequencer.Wait{Duration: 4 * time.Second},
				general.KeyDown{Key: key.VK_A},
				general.KeyUp{Key: key.VK_B},
				general.LeftMouseButtonDown{},
				general.LeftMouseButtonUp{},
				general.RightMouseButtonDown{},
				general.RightMouseButtonUp{},
				general.MouseMove{X: 10, Y: 20},
				general.SetMousePos{Pos: mouse.Pos{X: 30, Y: 40}},
			},
		},
		"nested repeats": {
			script: "KD D;R 1 [LD;W 60ms;LU;R 3 [W 800ms]];;;KU D;",
			want: []sequencer.Elem{
				general.KeyDown{Key: key.VK_D},
				general.LeftMouseButtonDown{},
				sequencer.Wait{Duration: time.Millisecond * 60},
				general.LeftMouseButtonUp{},
				sequencer.Wait{Duration: time.Millisecond * 800},
				sequencer.Wait{Duration: time.Millisecond * 800},
				sequencer.Wait{Duration: time.Millisecond * 800},
				general.KeyUp{Key: key.VK_D},
			},
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			var got = parser.NewCustomSequenceFunc(data.script)()
			if !reflect.DeepEqual(data.want, got) {
				t.Fatalf("sequences did not match:\ngot : %#v\nwant: %#v\n", got, data.want)
			}
		})
	}

}
