package action

import "fmt"

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
	ServerText    struct {
		Text string
	}
)

const (
	serverText           = "SERVER_TEXT"
	seqState             = "SEQUENCE_STATE"
	seqSetConfigSequence = "SEQUENCE_SETCONFIGSEQUENCE"
	seqClearSequence     = "SEQUENCE_CLEARSEQUENCE"
	seqPause             = "SEQUENCE_PAUSE"
	seqAbort             = "SEQUENCE_ABORT"
)

func msgState(s SequenceState) []byte {
	return jsonEnvelope(seqState, s)
}

func msgText(format string, a ...interface{}) []byte {
	return jsonEnvelope(serverText, ServerText{Text: fmt.Sprintf(format, a...)})
}
