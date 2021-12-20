package action

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Oppodelldog/roamer/internal/audioctl"

	"github.com/Oppodelldog/roamer/internal/config"

	"github.com/Oppodelldog/roamer/internal/server/ws"
)

type Action interface{}

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

				var message interface{}

				switch envelope.Type {
				case seqState:
					var state SequenceState
					err = json.NewDecoder(bytes.NewBuffer(envelope.Payload)).Decode(&state)
					message = state
				case seqClearSequence:
					var clearSequence SequenceClearSequence
					err = json.NewDecoder(bytes.NewBuffer(envelope.Payload)).Decode(&clearSequence)
					message = clearSequence
				case seqSetConfigSequence:
					var setConfigSequence SequenceSetConfigSequence
					err = json.NewDecoder(bytes.NewBuffer(envelope.Payload)).Decode(&setConfigSequence)
					message = setConfigSequence
				case seqPause:
					var pause SequencePause
					err = json.NewDecoder(bytes.NewBuffer(envelope.Payload)).Decode(&pause)
					message = pause
				case seqAbort:
					var abort SequenceAbort
					err = json.NewDecoder(bytes.NewBuffer(envelope.Payload)).Decode(&abort)
					message = abort
				case seqSave:
					var save SequenceSave
					err = json.NewDecoder(bytes.NewBuffer(envelope.Payload)).Decode(&save)
					save.Response = c.ToClient()
					message = save
				case loadSoundSettings:
					var lss LoadSoundSettings
					err = json.NewDecoder(bytes.NewBuffer(envelope.Payload)).Decode(&lss)
					lss.Response = c.ToClient()
					message = lss
				case setSoundVolume:
					var ssv SetSoundVolume
					err = json.NewDecoder(bytes.NewBuffer(envelope.Payload)).Decode(&ssv)
					message = ssv
				default:
					fmt.Printf("unknown message type: %v", envelope.Type)
				}

				if err != nil {
					fmt.Printf("error decoding client message '%T': %v", message, err)
					break
				}
				actions <- message
			}
		}
	}()
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

	return settings
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
