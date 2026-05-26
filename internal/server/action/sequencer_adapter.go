package action

import (
	"context"
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/mouse"
	"github.com/Oppodelldog/roamer/internal/sequencer"
)

type DefaultSequencerAdapter struct {
	seq *sequencer.Sequencer
}

func (d *DefaultSequencerAdapter) EnqueueSequence(f func() []sequencer.Elem) {
	d.seq.EnqueueSequence(f)
}

func (d *DefaultSequencerAdapter) IsPlaying() bool {
	return d.seq.IsPlaying()
}

func (d *DefaultSequencerAdapter) State() sequencer.State {
	return d.seq.State()
}

func (d *DefaultSequencerAdapter) HasSequence() bool {
	return d.seq.HasSequence()
}

func (d *DefaultSequencerAdapter) Pause() error {
	return d.seq.Pause()
}

func (d *DefaultSequencerAdapter) Abort() {
	d.seq.Abort()
	d.ReleaseInputs()
}

func (d *DefaultSequencerAdapter) ReleaseInputs() {
	key.ResetPressed()
	mouse.ResetPressed()
}

func NewSequencerAdapter(beforeSequence, afterSequence func(s *sequencer.Sequencer), elementError func(s *sequencer.Sequencer, elem sequencer.Elem, err error)) *DefaultSequencerAdapter {
	var seq = sequencer.New(context.Background(), 1)

	seq.BeforeSequence(func(s *sequencer.Sequencer) {
		key.ResetPressed()
		mouse.ResetPressed()
		beforeSequence(s)
	})
	seq.AfterSequence(func(s *sequencer.Sequencer) {
		afterSequence(s)
	})
	seq.ElementError(func(s *sequencer.Sequencer, elem sequencer.Elem, err error) {
		elementError(s, elem, err)
	})

	return &DefaultSequencerAdapter{seq: seq}
}
