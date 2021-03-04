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

	seq.EnqueueSequence(sequences.GetSequenceFunc(Sequence))
}

func hStop(w http.ResponseWriter, _ *http.Request) {
	err := seq.Pause()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Server Error: %v", err.Error())))
	}
}

func hState(w http.ResponseWriter, _ *http.Request) {
	var state = struct {
		Sequence string
		IsPaused bool
	}{
		Sequence: Sequence,
		IsPaused: seq.IsPaused()}

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
