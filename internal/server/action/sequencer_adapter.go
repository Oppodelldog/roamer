package action

import (
	"github.com/Oppodelldog/roamer/internal/key"
	"github.com/Oppodelldog/roamer/internal/sequencer"
)

type defaultSequencerAdapter struct {
	seq *sequencer.Sequencer
}

func (d *defaultSequencerAdapter) EnqueueSequence(f func() []sequencer.Elem) {
	d.seq.EnqueueSequence(f)
}

func (d *defaultSequencerAdapter) IsPaused() bool {
	return d.seq.IsPaused()
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

func NewSequencerAdapter() *defaultSequencerAdapter {
	var seq = sequencer.New(1)

	seq.BeforeSequence(func() {
		key.ResetPressed()
	})

	return &defaultSequencerAdapter{seq: seq}
}
