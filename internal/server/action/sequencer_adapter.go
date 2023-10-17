package action

import (
	"context"
	"github.com/Oppodelldog/roamer/internal/key"
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

func (d *DefaultSequencerAdapter) HasSequence() bool {
	return d.seq.HasSequence()
}

func (d *DefaultSequencerAdapter) Pause() error {
	return d.seq.Pause()
}

func (d *DefaultSequencerAdapter) Abort() {
	d.seq.Abort()
}

func NewSequencerAdapter(afterSequence func(s *sequencer.Sequencer)) *DefaultSequencerAdapter {
	var seq = sequencer.New(context.Background(), 1)

	seq.BeforeSequence(func(s *sequencer.Sequencer) {
		key.ResetPressed()
	})
	seq.AfterSequence(func(s *sequencer.Sequencer) {
		afterSequence(s)
	})

	return &DefaultSequencerAdapter{seq: seq}
}
