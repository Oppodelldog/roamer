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
	SequenceSave  struct {
		PageId        string
		SequenceIndex int
		Sequence      string
		Response      chan<- []byte
	}
	SequenceSaveResult struct {
		PageId        string
		SequenceIndex int
		Sequence      string
		Success       bool
	}
	SoundSession struct {
		Id    string
		Name  string
		Icon  string
		Value float32
	}
	SoundSettings struct {
		Sessions    []SoundSession
		MainSession SoundSession
	}
	LoadSoundSettings struct {
		Response chan<- []byte
	}
	SetSoundVolume struct {
		Id    string
		Value float32
	}
	SetMainSoundVolume struct {
		Value float32
	}
)

const (
	roamerConfig         = "CONFIG"
	seqState             = "SEQUENCE_STATE"
	seqSetConfigSequence = "SEQUENCE_SETCONFIGSEQUENCE"
	seqClearSequence     = "SEQUENCE_CLEARSEQUENCE"
	seqPause             = "SEQUENCE_PAUSE"
	seqAbort             = "SEQUENCE_ABORT"
	seqSave              = "SEQUENCE_SAVE"
	seqSaveResult        = "SEQUENCE_SAVE_RESULT"
	soundSettings        = "SOUND_SETTINGS"
	loadSoundSettings    = "LOAD_SOUND_SETTINGS"
	setSoundVolume       = "SET_SOUND_VOLUME"
	setMainSoundVolume   = "SET_MAIN_SOUND_VOLUME"
)

func msgState(s SequenceState) []byte {
	return jsonEnvelope(seqState, s)
}

func msgConfig(config config.Config) []byte {
	return jsonEnvelope(roamerConfig, config)
}

func msgSoundSettings(settings SoundSettings) []byte {
	return jsonEnvelope(soundSettings, settings)
}

func msgSequenceSaveResult(pageId string, sequenceIndex int, sequence string, success bool) []byte {
	var result = SequenceSaveResult{
		PageId:        pageId,
		SequenceIndex: sequenceIndex,
		Sequence:      sequence,
		Success:       success,
	}

	return jsonEnvelope(seqSaveResult, result)
}
