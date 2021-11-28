package action

import (
	"fmt"

	"github.com/Oppodelldog/roamer/internal/config"
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/script"
	"github.com/Oppodelldog/roamer/internal/sequencer"
)

func Worker(actions <-chan Action, broadcast chan<- []byte) {
	var (
		sequence = ""
		seq      = sequencer.New(1)
	)

	seq.BeforeSequence(func() {
		key.ResetPressed()
	})

	go func() {
		for action := range actions {
			switch v := action.(type) {
			case SequenceClearSequence:
				sequence = ""
			case SequenceSetConfigSequence:
				page, ok := config.RoamerPage(v.PageId)
				if !ok {
					fmt.Printf("roamer-page '%s' not found", v.PageId)
					return
				}

				if len(page.Actions) <= v.SequenceIndex {
					fmt.Printf("roamer-page '%s' not found", v.PageId)
					return
				}

				var action = page.Actions[v.SequenceIndex]

				elems, err := script.Parse(action.Sequence)
				if err != nil {
					panic(err)
				}

				fmt.Println(script.Write(elems))
				seq.EnqueueSequence(func() []sequencer.Elem {
					return elems
				})

				sequence = action.Sequence
			case SequencePause:
				err := seq.Pause()
				if err != nil {
					panic(err)
				}
			case SequenceAbort:
				seq.Abort()

				sequence = ""
			case LoadSoundSettings:
				v.Response <- msgSoundSettings(getSoundSettings())
				continue
			case SetSoundVolume:
				setSoundVol(v.Id, v.Value)
				continue
			default:
				fmt.Printf("unknown action: %T\n", action)
			}

			broadcast <- msgState(SequenceState{
				Sequence:    sequence,
				IsPaused:    seq.IsPaused(),
				HasSequence: seq.HasSequence()})
		}
	}()
}
