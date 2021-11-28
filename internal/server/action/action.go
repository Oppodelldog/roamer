package action

import (
	"github.com/Oppodelldog/roamer/internal/config"
)

type (
	SequenceState struct {
		Sequence    string
		IsPaused    bool
		HasSequence bool
	}
	SequenceClearSequence struct {
	}
	SequenceSetConfigSequence struct {
		PageId        string
		SequenceIndex int
	}
	SequencePause struct{}
	SequenceAbort struct{}
)

const (
	roamerConfig         = "CONFIG"
	seqState             = "SEQUENCE_STATE"
	seqSetConfigSequence = "SEQUENCE_SETCONFIGSEQUENCE"
	seqClearSequence     = "SEQUENCE_CLEARSEQUENCE"
	seqPause             = "SEQUENCE_PAUSE"
	seqAbort             = "SEQUENCE_ABORT"
)

func msgState(s SequenceState) []byte {
	return jsonEnvelope(seqState, s)
}

func msgConfig(config config.Config) []byte {
	return jsonEnvelope(roamerConfig, config)
}
