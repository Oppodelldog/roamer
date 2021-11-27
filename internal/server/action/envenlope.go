package action

import (
	"encoding/json"
	"fmt"
)

type Envelope struct {
	Type    string
	Payload json.RawMessage
}

func jsonEnvelope(messageType string, data interface{}) []byte {
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("cannot precompute payload type %T: %v\n", data, err)
	}

	b, err := json.Marshal(Envelope{Type: messageType, Payload: payload})
	if err != nil {
		fmt.Printf("cannot json encode message for payload type %T: %v", data, err)
	}

	return b
}
