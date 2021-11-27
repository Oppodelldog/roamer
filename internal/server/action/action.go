package action

import "fmt"

type (
	SequenceState struct {
		Sequence    string
		IsPaused    bool
		HasSequence bool
	}
	SequenceSetSequence struct {
		Sequence string
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
	seqSetSequence       = "SEQUENCE_SETSEQUENCE"
	seqSetConfigSequence = "SEQUENCE_SETCONFIGSEQUENCE"
	seqPause             = "SEQUENCE_PAUSE"
	seqAbort             = "SEQUENCE_ABORT"
)

func msgState(s SequenceState) []byte {
	return jsonEnvelope(seqState, s)
}

func msgText(format string, a ...interface{}) []byte {
	return jsonEnvelope(serverText, ServerText{Text: fmt.Sprintf(format, a...)})
}
