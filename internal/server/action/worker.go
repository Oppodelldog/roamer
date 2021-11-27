package action

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/config"
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/script"
	"github.com/Oppodelldog/roamer/internal/sequencer"
)

func Worker(actions <-chan Action, broadcast chan<- []byte) {
	var seq = sequencer.New(1)
	seq.BeforeSequence(func() {
		key.ResetPressed()
	})
	var sequence = ""

	go func() {
		for {
			select {
			case action := <-actions:
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

					s, err := script.Write(elems)
					if err != nil {
						panic(err)
					}
					fmt.Println(s)
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
				default:
					fmt.Printf("unknown action: %T\n", action)
				}

				broadcast <- msgState(SequenceState{
					Sequence:    sequence,
					IsPaused:    seq.IsPaused(),
					HasSequence: seq.HasSequence()})
			}
		}
	}()
}
