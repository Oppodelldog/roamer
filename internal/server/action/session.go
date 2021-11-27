package action

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

		c.ToClient() <- msgText("hello client")
		for {
			select {
			case <-ctx.Done():
				fmt.Print("client session closed\n")

				return
			case msg := <-c.ToServer():
				var envelope Envelope
				var err = json.NewDecoder(bytes.NewBuffer(msg)).Decode(&envelope)
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
