package script_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
	"github.com/Oppodelldog/roamer/internal/script"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences/general"
)

func TestNewCustomSequenceFunc(t *testing.T) {
	type testData struct {
		script string
		want   []sequencer.Elem
	}

	var tests = map[string]testData{
		"all commands": {
			script: "NOP W 3ms;L;R 3 [W 4s];KD A;KU B;LD;LU;RD;RU;MM 10 20;SM 30 40;MP;W 1.2s;SM -1 -2",
			want: []sequencer.Elem{
				sequencer.NoOperation{},
				sequencer.Wait{Duration: time.Millisecond * 3},
				sequencer.Loop{},
				sequencer.Wait{Duration: 4 * time.Second},
				sequencer.Wait{Duration: 4 * time.Second},
				sequencer.Wait{Duration: 4 * time.Second},
				general.KeyDown{Key: key.VkA},
				general.KeyUp{Key: key.VkB},
				general.LeftMouseButtonDown{},
				general.LeftMouseButtonUp{},
				general.RightMouseButtonDown{},
				general.RightMouseButtonUp{},
				general.MouseMove{X: 10, Y: 20},
				general.SetMousePos{Pos: mouse.Pos{X: 30, Y: 40}},
				general.LookupMousePos{},
				sequencer.Wait{Duration: 1*time.Second + 200*time.Millisecond},
				general.SetMousePos{Pos: mouse.Pos{X: -1, Y: -2}},
			},
		},
		"ignore spaces": {
			script: "W  3ms  ;  L ; R  3  [ W  4s ] ; KD  A ;  KU B ; LD  ;  LU ; RD ; RU  ;  MM  10  20  ; SM  30  40 ; MP ; W  1.2s ; SM  -1  -2 ",
			want: []sequencer.Elem{
				sequencer.Wait{Duration: time.Millisecond * 3},
				sequencer.Loop{},
				sequencer.Wait{Duration: 4 * time.Second},
				sequencer.Wait{Duration: 4 * time.Second},
				sequencer.Wait{Duration: 4 * time.Second},
				general.KeyDown{Key: key.VkA},
				general.KeyUp{Key: key.VkB},
				general.LeftMouseButtonDown{},
				general.LeftMouseButtonUp{},
				general.RightMouseButtonDown{},
				general.RightMouseButtonUp{},
				general.MouseMove{X: 10, Y: 20},
				general.SetMousePos{Pos: mouse.Pos{X: 30, Y: 40}},
				general.LookupMousePos{},
				sequencer.Wait{Duration: 1*time.Second + 200*time.Millisecond},
				general.SetMousePos{Pos: mouse.Pos{X: -1, Y: -2}},
			},
		},
		"nested repeats": {
			script: "KD D;R 1 [LD;W 60ms;LU;R 3 [W 800ms]];;;KU D;",
			want: []sequencer.Elem{
				general.KeyDown{Key: key.VkD},
				general.LeftMouseButtonDown{},
				sequencer.Wait{Duration: time.Millisecond * 60},
				general.LeftMouseButtonUp{},
				sequencer.Wait{Duration: time.Millisecond * 800},
				sequencer.Wait{Duration: time.Millisecond * 800},
				sequencer.Wait{Duration: time.Millisecond * 800},
				general.KeyUp{Key: key.VkD},
			},
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			var got, err = script.Parse(data.script)
			if err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}

			if !reflect.DeepEqual(data.want, got) {
				t.Fatalf("sequences did not match:\ngot : %#v\nwant: %#v\n", got, data.want)
			}
		})
	}
}
