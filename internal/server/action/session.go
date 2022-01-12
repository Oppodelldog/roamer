package action

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Oppodelldog/roamer/internal/audioctl"

	"github.com/Oppodelldog/roamer/internal/config"

	"github.com/Oppodelldog/roamer/internal/server/ws"
)

type Action interface{}

var actionTypesById = map[string]func() Responder{
	seqClearSequence:     func() Responder { return new(SequenceClearSequence) },
	seqSetConfigSequence: func() Responder { return new(SequenceSetConfigSequence) },
	seqPause:             func() Responder { return new(SequencePause) },
	seqAbort:             func() Responder { return new(SequenceAbort) },
	seqSave:              func() Responder { return new(SequenceSave) },
	loadSoundSettings:    func() Responder { return new(LoadSoundSettings) },
	setSoundVolume:       func() Responder { return new(SetSoundVolume) },
	setMainSoundVolume:   func() Responder { return new(SetMainSoundVolume) },
	pageNew:              func() Responder { return new(PageNew) },
	pageDelete:           func() Responder { return new(PageDelete) },
	macroNew:             func() Responder { return new(SequenceNew) },
	macroDelete:          func() Responder { return new(SequenceDelete) },
	pagesSave:            func() Responder { return new(PagesSave) },
}

func ClientSession(c *ws.Client, actions chan Action) {
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

				if newType, ok := actionTypesById[envelope.Type]; ok {
					message, errDecode := decode(newType(), envelope.Payload, c.ToClient())
					if errDecode != nil {
						fmt.Printf("error decoding client message '%T': %v\n", message, errDecode)
						break
					}
					actions <- message
				} else {
					fmt.Printf("unknown message type: %v", envelope.Type)
				}
			}
		}
	}()
}

func decode(t interface{}, payload json.RawMessage, toClient chan<- []byte) (interface{}, error) {
	err := json.NewDecoder(bytes.NewBuffer(payload)).Decode(&t)

	if t, ok := t.(Responder); ok {
		t.SetRespondChannel(toClient)
	} else {
		panic("target type of decoding must be a Responder")
	}

	t = reflect.ValueOf(t).Elem().Interface()

	return t, err
}

func getSoundSettings() SoundSettings {
	var (
		settings      = SoundSettings{}
		soundSessions []SoundSession
	)

	device, err := audioctl.DefaultDevice()
	if err != nil {
		fmt.Println(err)
		return settings
	}

	defer device.Release()

	for i, session := range device.Sessions {
		var (
			err         error
			sessionId   string
			displayName string
			iconPath    string
			volume      float32
		)

		sessionId, err = session.GetSessionInstanceIdentifier()
		if err != nil {
			fmt.Printf("error getting sound session (%v) display name: %v", i, err)
			continue
		}

		displayName, err = session.GetDisplayNameEnhanced()
		if err != nil {
			fmt.Printf("error getting sound session (%v) display name: %v", i, err)
			continue
		}

		iconPath, err = session.GetIconPath()
		if err != nil {
			fmt.Printf("error getting sound session (%v) icon path: %v", i, err)
			continue
		}

		volume, err = session.GetMasterVolume()
		if err != nil {
			fmt.Printf("error getting sound session (%v) volume: %v", i, err)
			continue
		}

		soundSessions = append(soundSessions, SoundSession{
			Id:    sessionId,
			Name:  displayName,
			Icon:  iconPath,
			Value: volume,
		})
	}

	settings.Sessions = soundSessions
	settings.MainSession = getMasterSoundSession(device)

	return settings
}

func getMasterSoundSession(dev *audioctl.Device) SoundSession {
	scalar, err := dev.GetMasterVolumeLevelScalar()
	if err != nil {
		fmt.Printf("cannot get main volume: %v\n", err)
		return SoundSession{}
	}

	return SoundSession{
		Name:  "Main Volume",
		Icon:  "",
		Value: scalar,
	}
}

func setMainSoundVol(value float32) {
	device, err := audioctl.DefaultDevice()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer device.Release()

	err = device.SetMasterVolumeLevelScalar(value)
	if err != nil {
		fmt.Println(err)
	}
}

func setSoundVol(id string, value float32) {
	device, err := audioctl.DefaultDevice()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer device.Release()

	for i, session := range device.Sessions {
		sid, err := session.GetSessionInstanceIdentifier()
		if err != nil {
			fmt.Printf("error getting sound session (%v) id: %v", i, err)

			continue
		}

		if sid == id {
			err = session.SetMasterVolume(value)
			if err != nil {
				fmt.Printf("error getting sound session (%v) id: %v", i, err)
			}

			return
		}
	}
}
