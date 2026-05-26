package action

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/config"
	"github.com/Oppodelldog/roamer/internal/logger"
	"github.com/Oppodelldog/roamer/internal/script"
	"github.com/Oppodelldog/roamer/internal/sequencer"
)

type SequencerAdapter interface {
	Pause() error
	Abort()
	State() sequencer.State
	IsPlaying() bool
	HasSequence() bool
	EnqueueSequence(func() []sequencer.Elem)
	ReleaseInputs()
}

type sequenceInfo struct {
	pageId        string
	pageTitle     string
	sequenceIndex int
	caption       string
	loops         bool
}

func StartSequencerWorker(actions <-chan Action, broadcast chan<- []byte) {
	var (
		active = sequenceInfo{sequenceIndex: -1}
		queued = sequenceInfo{sequenceIndex: -1}
	)

	broadcastSequenceState := func(s *sequencer.Sequencer) {
		state := s.State().String()
		if state == "playing" && active.loops {
			state = "looping"
		}

		if state == "idle" && active.sequenceIndex >= 0 {
			state = "stopped"
		}

		broadcast <- msgState(SequenceState{
			PageId:              active.pageId,
			SequenceIndex:       active.sequenceIndex,
			PageTitle:           active.pageTitle,
			Caption:             active.caption,
			State:               state,
			IsPlaying:           s.IsPlaying(),
			HasSequence:         s.HasSequence(),
			QueuedPageId:        queued.pageId,
			QueuedSequenceIndex: queued.sequenceIndex,
			QueuedPageTitle:     queued.pageTitle,
			QueuedCaption:       queued.caption})
	}

	broadcastSequenceError := func(s *sequencer.Sequencer, elem sequencer.Elem, err error) {
		broadcast <- msgState(SequenceState{
			PageId:              active.pageId,
			SequenceIndex:       active.sequenceIndex,
			PageTitle:           active.pageTitle,
			Caption:             active.caption,
			State:               "error",
			Error:               fmt.Sprintf("%T failed: %v", elem, err),
			IsPlaying:           s.IsPlaying(),
			HasSequence:         s.HasSequence(),
			QueuedPageId:        queued.pageId,
			QueuedSequenceIndex: queued.sequenceIndex,
			QueuedPageTitle:     queued.pageTitle,
			QueuedCaption:       queued.caption})
	}

	beforeSequence := func(s *sequencer.Sequencer) {
		if queued.sequenceIndex >= 0 {
			active = queued
			queued = sequenceInfo{sequenceIndex: -1}
		}

		broadcastSequenceState(s)
	}

	var seq SequencerAdapter = NewSequencerAdapter(beforeSequence, broadcastSequenceState, broadcastSequenceError)

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

				next := sequenceInfo{
					pageId:        v.PageId,
					pageTitle:     page.Title,
					sequenceIndex: v.SequenceIndex,
					caption:       action.Caption,
					loops:         script.AnalyzeElems(elems).HasLoop,
				}

				queued = next
				if !seq.HasSequence() {
					active = next
				}

				fmt.Println(script.Write(elems))
				seq.EnqueueSequence(func() []sequencer.Elem {
					return elems
				})
			case SequenceClearSequence:
				active = sequenceInfo{sequenceIndex: -1}
				queued = sequenceInfo{sequenceIndex: -1}

			case SequencePause:
				err := seq.Pause()
				if err != nil {
					fmt.Println(err)
					continue
				}
			case SequenceAbort:
				seq.Abort()
				queued = sequenceInfo{sequenceIndex: -1}
			case SequenceReleaseInputs:
				seq.ReleaseInputs()
				logger.Println("released pressed inputs")
			default:
				fmt.Printf("unknown action for sequence worker: %T\n", action)
			}

			state := seq.State().String()
			if state == "playing" && active.loops {
				state = "looping"
			}
			if state == "idle" && active.sequenceIndex >= 0 {
				state = "stopped"
			}

			broadcast <- msgState(SequenceState{
				PageId:              active.pageId,
				SequenceIndex:       active.sequenceIndex,
				PageTitle:           active.pageTitle,
				Caption:             active.caption,
				State:               state,
				IsPlaying:           seq.IsPlaying(),
				HasSequence:         seq.HasSequence(),
				QueuedPageId:        queued.pageId,
				QueuedSequenceIndex: queued.sequenceIndex,
				QueuedPageTitle:     queued.pageTitle,
				QueuedCaption:       queued.caption})
		}
	}()
}
