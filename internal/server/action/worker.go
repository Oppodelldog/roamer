package action

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/config"
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/script"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences"
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
				case SequenceSetSequence:
					seq.EnqueueSequence(sequences.NewBuildInSequenceFunc(v.Sequence))
					sequence = v.Sequence
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

					seq.EnqueueSequence(func() []sequencer.Elem {
						seq, err := script.Parse(page.Actions[v.SequenceIndex].Sequence)
						if err != nil {
							panic(err)
						}

						return seq
					})
					sequence = page.Actions[v.SequenceIndex].Sequence
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
