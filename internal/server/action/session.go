package action

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Oppodelldog/roamer/internal/config"

	"github.com/Oppodelldog/roamer/internal/server/ws"
)

type Action interface{}

var actionTypesById = map[string]func() Responder{
	macroNew:    func() Responder { return new(SequenceNew) },
	seqSave:     func() Responder { return new(SequenceSave) },
	macroDelete: func() Responder { return new(SequenceDelete) },
	pageNew:     func() Responder { return new(PageNew) },
	pagesSave:   func() Responder { return new(PagesSave) },
	pageDelete:  func() Responder { return new(PageDelete) },
}

var sequencerActionTypesById = map[string]func() Responder{
	seqSetConfigSequence: func() Responder { return new(SequenceSetConfigSequence) },
	seqClearSequence:     func() Responder { return new(SequenceClearSequence) },
	seqPause:             func() Responder { return new(SequencePause) },
	seqAbort:             func() Responder { return new(SequenceAbort) },
}

var soundSettingActions = map[string]func() Responder{
	loadSoundSettings:  func() Responder { return new(LoadSoundSettings) },
	setSoundVolume:     func() Responder { return new(SetSoundVolume) },
	setMainSoundVolume: func() Responder { return new(SetMainSoundVolume) },
}

func ClientSession(c *ws.Client, actions, sequencerActions, soundActions, loggerActions chan Action) {
	var ctx, cancel = context.WithCancel(context.Background())
	c.CancelSession = cancel

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic in search session: %v\n", r)
			}
		}()

		c.ToClient() <- msgConfig(config.Roamer())

		for {
			select {
			case <-ctx.Done():
				fmt.Print("client session closed\n")

				return
			case msg := <-c.ToServer():
				var (
					envelope Envelope
					err      = json.NewDecoder(bytes.NewBuffer(msg)).Decode(&envelope)
				)

				if err != nil {
					fmt.Printf("error decoding client message envelope: %v", err)
					break
				}

				var (
					targetChannel chan Action
					newResponder  func() Responder
					exists        bool
				)

				if newResponder, exists = actionTypesById[envelope.Type]; exists {
					targetChannel = actions
				} else if newResponder, exists = sequencerActionTypesById[envelope.Type]; exists {
					targetChannel = sequencerActions
				} else if newResponder, exists = soundSettingActions[envelope.Type]; exists {
					targetChannel = soundActions
				} else {
					fmt.Printf("unknown message type: %v", envelope.Type)
					break
				}

				message, errDecode := decode(newResponder(), envelope.Payload, c.ToClient())
				if errDecode != nil {
					fmt.Printf("error decoding client message '%T': %v\n", message, errDecode)
					break
				}

				targetChannel <- message
			}
		}
	}()
}

func decode(t interface{}, payload json.RawMessage, toClient chan<- []byte) (interface{}, error) {
	err := json.NewDecoder(bytes.NewBuffer(payload)).Decode(&t)

	if r, ok := t.(Responder); ok {
		r.SetRespondChannel(toClient)
	} else {
		panic(fmt.Sprintf("target type of decoding must be a Responder, but it was %T", t))
	}

	t = reflect.ValueOf(t).Elem().Interface()

	return t, err
}
