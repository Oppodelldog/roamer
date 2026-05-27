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
	steps         []sequencer.ElementEvent
	totalMs       int64
}

func StartSequencerWorker(actions <-chan Action, broadcast chan<- []byte) {
	var (
		active          = sequenceInfo{sequenceIndex: -1}
		queued          = sequenceInfo{sequenceIndex: -1}
		activeStepIndex = 0
		activeCycle     = 0
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
			activeCycle = 0
		}

		activeStepIndex = 0
		activeCycle++
		broadcastSequenceState(s)
	}

	broadcastElementEvent := func(s *sequencer.Sequencer, elem sequencer.Elem, event sequencer.ElementEvent) {
		if _, isLoop := elem.(sequencer.Loop); isLoop {
			broadcast <- msgSequenceElementEvent(SequenceElementEvent{
				PageId:        active.pageId,
				SequenceIndex: active.sequenceIndex,
				PageTitle:     active.pageTitle,
				Caption:       active.caption,
				Label:         event.Label,
				DurationMs:    event.DurationMs,
				StepIndex:     len(active.steps),
				TotalSteps:    len(active.steps),
				NextLabel:     firstStepLabel(active.steps),
				RemainingMs:   0,
				Progress:      100,
				Cycle:         activeCycle,
				IsLoop:        true,
			})
			activeStepIndex = 0
			return
		}

		stepIndex := activeStepIndex
		activeStepIndex++
		if len(active.steps) > 0 && activeStepIndex > len(active.steps) {
			activeStepIndex = len(active.steps)
		}

		broadcast <- msgSequenceElementEvent(SequenceElementEvent{
			PageId:        active.pageId,
			SequenceIndex: active.sequenceIndex,
			PageTitle:     active.pageTitle,
			Caption:       active.caption,
			Label:         event.Label,
			DurationMs:    event.DurationMs,
			StepIndex:     stepIndex + 1,
			TotalSteps:    len(active.steps),
			NextLabel:     nextStepLabel(active.steps, stepIndex+1),
			RemainingMs:   remainingMs(active.steps, stepIndex),
			Progress:      progressPercent(stepIndex+1, len(active.steps)),
			Cycle:         activeCycle,
		})
	}

	var seq SequencerAdapter = NewSequencerAdapter(beforeSequence, broadcastSequenceState, broadcastElementEvent, broadcastSequenceError)

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
				next.steps, next.totalMs = timelineSteps(elems)

				queued = next
				if !seq.HasSequence() {
					active = next
					activeCycle = 0
				}

				fmt.Println(script.Write(elems))
				seq.EnqueueSequence(func() []sequencer.Elem {
					return elems
				})
			case SequenceClearSequence:
				active = sequenceInfo{sequenceIndex: -1}
				queued = sequenceInfo{sequenceIndex: -1}
				activeStepIndex = 0
				activeCycle = 0

			case SequencePause:
				err := seq.Pause()
				if err != nil {
					fmt.Println(err)
					continue
				}
			case SequenceAbort:
				seq.Abort()
				queued = sequenceInfo{sequenceIndex: -1}
				activeStepIndex = 0
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

func timelineSteps(elems []sequencer.Elem) ([]sequencer.ElementEvent, int64) {
	steps := []sequencer.ElementEvent{}
	var totalMs int64

	for _, elem := range elems {
		if _, isLoop := elem.(sequencer.Loop); isLoop {
			continue
		}

		event := sequencer.DescribeElement(elem)
		steps = append(steps, event)
		totalMs += event.DurationMs
	}

	return steps, totalMs
}

func firstStepLabel(steps []sequencer.ElementEvent) string {
	if len(steps) == 0 {
		return ""
	}

	return steps[0].Label
}

func nextStepLabel(steps []sequencer.ElementEvent, nextIndex int) string {
	if nextIndex < 0 || nextIndex >= len(steps) {
		return ""
	}

	return steps[nextIndex].Label
}

func remainingMs(steps []sequencer.ElementEvent, nextIndex int) int64 {
	var remaining int64
	for i := nextIndex; i < len(steps); i++ {
		remaining += steps[i].DurationMs
	}

	return remaining
}

func progressPercent(stepIndex, totalSteps int) int {
	if totalSteps <= 0 {
		return 0
	}

	progress := stepIndex * 100 / totalSteps
	if progress > 100 {
		return 100
	}

	return progress
}
