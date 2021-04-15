package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"rust-roamer/sequences"
)

var Sequence string

func hSet(_ http.ResponseWriter, r *http.Request) {
	Sequence = path.Base(r.URL.Path)

	seq.EnqueueSequence(sequences.NewSequenceFunc(Sequence))
}

func hPause(w http.ResponseWriter, _ *http.Request) {
	err := seq.Pause()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Server Error: %v", err.Error())))
	}
}

func hAbort(w http.ResponseWriter, _ *http.Request) {
	seq.Abort()
	Sequence = ""
}

func hState(w http.ResponseWriter, _ *http.Request) {
	var state = struct {
		Sequence    string
		IsPaused    bool
		HasSequence bool
	}{
		Sequence:    Sequence,
		IsPaused:    seq.IsPaused(),
		HasSequence: seq.HasSequence()}

	jsonBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		panic(err)
	}
}
