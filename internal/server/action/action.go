package action

import (
	"github.com/Oppodelldog/roamer/internal/config"
)

const (
	// server -> client
	seqState      = "SEQUENCE_STATE"
	seqSaveResult = "SEQUENCE_SAVE_RESULT"
	soundSettings = "SOUND_SETTINGS"
	roamerConfig  = "CONFIG"

	// client -> server
	seqSetConfigSequence = "SEQUENCE_SETCONFIGSEQUENCE"
	seqClearSequence     = "SEQUENCE_CLEARSEQUENCE"
	seqPause             = "SEQUENCE_PAUSE"
	seqAbort             = "SEQUENCE_ABORT"
	macroNew             = "SEQUENCE_NEW"
	macroDelete          = "SEQUENCE_DELETE"
	seqSave              = "SEQUENCE_SAVE"
	loadSoundSettings    = "LOAD_SOUND_SETTINGS"
	setSoundVolume       = "SET_SOUND_VOLUME"
	setMainSoundVolume   = "SET_MAIN_SOUND_VOLUME"
	pageNew              = "PAGE_NEW"
	pageDelete           = "PAGE_DELETE"
	pagesSave            = "PAGES_SAVE"
)

type (
	SequenceState struct {
		PageTitle   string
		Caption     string
		IsPlaying   bool
		HasSequence bool
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
)

type (
	Responder interface {
		SetRespondChannel(chan<- []byte)
	}
	BaseAction struct {
		Response chan<- []byte `json:"-"`
	}
	SequenceClearSequence struct {
		BaseAction
	}
	SequenceSetConfigSequence struct {
		BaseAction
		PageId        string
		SequenceIndex int
	}
	SequencePause struct {
		BaseAction
	}
	SequenceAbort struct {
		BaseAction
	}
	SequenceNew struct {
		BaseAction
		PageId string
	}
	SequenceDelete struct {
		BaseAction
		PageId        string
		SequenceIndex int
	}
	SequenceSave struct {
		BaseAction
		PageId        string
		SequenceIndex int
		Caption       string
		Sequence      string
	}
	SequenceSaveResult struct {
		BaseAction
		PageId        string
		SequenceIndex int
		Sequence      string
		Success       bool
	}
	LoadSoundSettings struct {
		BaseAction
	}
	SetSoundVolume struct {
		BaseAction
		Id    string
		Value float32
	}
	SetMainSoundVolume struct {
		BaseAction
		Value float32
	}
	PageNew struct {
		BaseAction
	}
	PageDelete struct {
		BaseAction
		PageId string
	}
	PagesSave struct {
		BaseAction
		Pages config.Pages
	}
)

func (r *BaseAction) SetRespondChannel(c chan<- []byte) {
	r.Response = c
}

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
