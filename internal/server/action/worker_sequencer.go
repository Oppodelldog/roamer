package action

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/config"
	"github.com/Oppodelldog/roamer/internal/script"
	"github.com/Oppodelldog/roamer/internal/sequencer"
)

type SequencerAdapter interface {
	Pause() error
	Abort()
	IsPaused() bool
	HasSequence() bool
	EnqueueSequence(func() []sequencer.Elem)
}

func StartSequencerWorker(actions <-chan Action, broadcast chan<- []byte) {
	var (
		seq             SequencerAdapter = NewSequencerAdapter()
		sequenceCaption                  = ""
		pageTitle                        = ""
	)

	go func() {
		for action := range actions {
			switch v := action.(type) {
			case SequenceSetConfigSequence:
				page, ok := config.RoamerPage(v.PageId)
				if !ok {
					fmt.Printf("roamer-page '%s' not found", v.PageId)
					continue
				}

				if len(page.Actions) <= v.SequenceIndex {
					fmt.Printf("roamer-page '%s' not found", v.PageId)
					continue
				}

				var action = page.Actions[v.SequenceIndex]

				elems, err := script.Parse(action.Sequence)
				if err != nil {
					fmt.Println(err)
					continue
				}

				fmt.Println(script.Write(elems))
				seq.EnqueueSequence(func() []sequencer.Elem {
					return elems
				})

				sequenceCaption = action.Caption
				pageTitle = page.Title
			case SequenceClearSequence:
				sequenceCaption = ""
				pageTitle = ""

			case SequencePause:
				err := seq.Pause()
				if err != nil {
					fmt.Println(err)
					continue
				}
			case SequenceAbort:
				seq.Abort()

				sequenceCaption = ""
				pageTitle = ""
			default:
				fmt.Printf("unknown action for sequence worker: %T\n", action)
			}

			broadcast <- msgState(SequenceState{
				PageTitle:   pageTitle,
				Caption:     sequenceCaption,
				IsPaused:    seq.IsPaused(),
				HasSequence: seq.HasSequence()})
		}
	}()
}
