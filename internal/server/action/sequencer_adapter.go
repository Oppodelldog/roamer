package action

import (
	"context"
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/sequencer"
)

type defaultSequencerAdapter struct {
	seq *sequencer.Sequencer
}

func (d *defaultSequencerAdapter) EnqueueSequence(f func() []sequencer.Elem) {
	d.seq.EnqueueSequence(f)
}

func (d *defaultSequencerAdapter) IsPlaying() bool {
	return d.seq.IsPlaying()
}

func (d *defaultSequencerAdapter) HasSequence() bool {
	return d.seq.HasSequence()
}

func (d *defaultSequencerAdapter) Pause() error {
	return d.seq.Pause()
}

func (d *defaultSequencerAdapter) Abort() {
	d.seq.Abort()
}

func NewSequencerAdapter(afterSequence func(s *sequencer.Sequencer)) *defaultSequencerAdapter {
	var seq = sequencer.New(context.Background(), 1)

	seq.BeforeSequence(func(s *sequencer.Sequencer) {
		key.ResetPressed()
	})
	seq.AfterSequence(func(s *sequencer.Sequencer) {
		afterSequence(s)
	})

	return &defaultSequencerAdapter{seq: seq}
}
