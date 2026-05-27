package action

import (
	"github.com/Oppodelldog/roamer/internal/config"
	"github.com/Oppodelldog/roamer/internal/inputmode"
	"github.com/Oppodelldog/roamer/internal/script"
)

const (
	// server -> client
	seqState          = "SEQUENCE_STATE"
	seqSaveResult     = "SEQUENCE_SAVE_RESULT"
	seqFormatResult   = "SEQUENCE_FORMAT_RESULT"
	seqValidateResult = "SEQUENCE_VALIDATE_RESULT"
	seqElementEvent   = "SEQUENCE_ELEMENT_EVENT"
	recorderState     = "RECORDER_STATE"
	recorderInput     = "RECORDER_INPUT_EVENT"
	remoteInfo        = "REMOTE_INFO"
	soundSettings     = "SOUND_SETTINGS"
	inputModeState    = "INPUT_MODE_STATE"
	logMessage        = "LOG_MESSAGE"
	roamerConfig      = "CONFIG"

	// client -> server
	seqSetConfigSequence = "SEQUENCE_SETCONFIGSEQUENCE"
	seqClearSequence     = "SEQUENCE_CLEARSEQUENCE"
	seqPause             = "SEQUENCE_PAUSE"
	seqAbort             = "SEQUENCE_ABORT"
	seqReleaseInputs     = "SEQUENCE_RELEASE_INPUTS"
	macroNew             = "SEQUENCE_NEW"
	macroDelete          = "SEQUENCE_DELETE"
	macroDuplicate       = "SEQUENCE_DUPLICATE"
	macroMove            = "SEQUENCE_MOVE"
	seqSave              = "SEQUENCE_SAVE"
	seqFormat            = "SEQUENCE_FORMAT"
	seqValidate          = "SEQUENCE_VALIDATE"
	loadSoundSettings    = "LOAD_SOUND_SETTINGS"
	setSoundVolume       = "SET_SOUND_VOLUME"
	setMainSoundVolume   = "SET_MAIN_SOUND_VOLUME"
	setInputMode         = "SET_INPUT_MODE"
	recorderStart        = "RECORDER_START"
	recorderStop         = "RECORDER_STOP"
	remoteMacroSave      = "REMOTE_MACRO_SAVE"
	pageNew              = "PAGE_NEW"
	pageDelete           = "PAGE_DELETE"
	pagesSave            = "PAGES_SAVE"
)

type (
	SequenceState struct {
		PageId              string
		SequenceIndex       int
		PageTitle           string
		Caption             string
		State               string
		Error               string
		IsPlaying           bool
		HasSequence         bool
		QueuedPageId        string
		QueuedSequenceIndex int
		QueuedPageTitle     string
		QueuedCaption       string
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
	InputModeState struct {
		Mode   string
		DryRun bool
	}
	SequenceElementEvent struct {
		PageId        string
		SequenceIndex int
		PageTitle     string
		Caption       string
		Label         string
		DurationMs    int64
		StepIndex     int
		TotalSteps    int
		NextLabel     string
		RemainingMs   int64
		Progress      int
		Cycle         int
		IsLoop        bool
	}
	RecorderState struct {
		Active   bool
		Sequence string
		Count    int
		Error    string
	}
	RecorderInputEvent struct {
		Label string
		Kind  string
	}
	RemoteInfo struct {
		Urls    []string
		Targets []RemoteTarget
	}
	RemoteTarget struct {
		Url    string
		QrCode string
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
	SequenceReleaseInputs struct {
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
	SequenceDuplicate struct {
		BaseAction
		PageId        string
		SequenceIndex int
	}
	SequenceMove struct {
		BaseAction
		PageId        string
		SequenceIndex int
		Offset        int
	}
	SequenceSave struct {
		BaseAction
		PageId        string
		SequenceIndex int
		Caption       string
		Icon          string
		Sequence      string
	}
	SequenceFormat struct {
		BaseAction
		PageId        string
		SequenceIndex int
		Sequence      string
	}
	SequenceValidate struct {
		BaseAction
		PageId        string
		SequenceIndex int
		Sequence      string
	}
	SequenceSaveResult struct {
		BaseAction
		PageId        string
		SequenceIndex int
		Sequence      string
		Meta          script.Metadata
		Success       bool
		Error         string
	}
	SequenceFormatResult struct {
		BaseAction
		PageId        string
		SequenceIndex int
		Sequence      string
		Meta          script.Metadata
		Success       bool
		Error         string
	}
	SequenceValidateResult struct {
		BaseAction
		PageId        string
		SequenceIndex int
		Sequence      string
		Meta          script.Metadata
		Success       bool
		Error         string
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
	SetInputMode struct {
		BaseAction
		DryRun bool
	}
	RecorderStart struct {
		BaseAction
	}
	RecorderStop struct {
		BaseAction
	}
	RemoteMacroSave struct {
		BaseAction
		PageId   string
		Caption  string
		Sequence string
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

func msgSequenceElementEvent(event SequenceElementEvent) []byte {
	return jsonEnvelope(seqElementEvent, event)
}

func msgRecorderState(state RecorderState) []byte {
	return jsonEnvelope(recorderState, state)
}

func msgRecorderInputEvent(event RecorderInputEvent) []byte {
	return jsonEnvelope(recorderInput, event)
}

func msgRemoteInfo(info RemoteInfo) []byte {
	return jsonEnvelope(remoteInfo, info)
}

func msgConfig(config config.Config) []byte {
	return jsonEnvelope(roamerConfig, configView(config))
}

func msgSoundSettings(settings SoundSettings) []byte {
	return jsonEnvelope(soundSettings, settings)
}

func msgInputModeState() []byte {
	return jsonEnvelope(inputModeState, InputModeState{
		Mode:   inputmode.Name(),
		DryRun: inputmode.IsDryRun(),
	})
}

func msgLogMessage(message string) []byte {
	return jsonEnvelope(logMessage, message)
}

func msgSequenceSaveResult(pageId string, sequenceIndex int, sequence string, meta script.Metadata, success bool, errMsg string) []byte {
	var result = SequenceSaveResult{
		PageId:        pageId,
		SequenceIndex: sequenceIndex,
		Sequence:      sequence,
		Meta:          meta,
		Success:       success,
		Error:         errMsg,
	}

	return jsonEnvelope(seqSaveResult, result)
}

func msgSequenceFormatResult(pageId string, sequenceIndex int, sequence string, meta script.Metadata, success bool, errMsg string) []byte {
	var result = SequenceFormatResult{
		PageId:        pageId,
		SequenceIndex: sequenceIndex,
		Sequence:      sequence,
		Meta:          meta,
		Success:       success,
		Error:         errMsg,
	}

	return jsonEnvelope(seqFormatResult, result)
}

func msgSequenceValidateResult(pageId string, sequenceIndex int, sequence string, meta script.Metadata, success bool, errMsg string) []byte {
	var result = SequenceValidateResult{
		PageId:        pageId,
		SequenceIndex: sequenceIndex,
		Sequence:      sequence,
		Meta:          meta,
		Success:       success,
		Error:         errMsg,
	}

	return jsonEnvelope(seqValidateResult, result)
}
