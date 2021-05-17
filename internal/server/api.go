package server

import (
	"encoding/json"
	"fmt"
	"github.com/Oppodelldog/roamer/internal/config"
	script2 "github.com/Oppodelldog/roamer/internal/script"
	"github.com/Oppodelldog/roamer/internal/sequencer"
	"github.com/Oppodelldog/roamer/internal/sequences"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

var Sequence string

func hSet(_ http.ResponseWriter, r *http.Request) {
	Sequence = path.Base(r.URL.Path)

	seq.EnqueueSequence(sequences.NewBuildInSequenceFunc(Sequence))
}

func hSetConfigSeq(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Split(r.URL.Path, "/")
	pageId := urlPath[len(urlPath)-2]
	actionIdx, err := strconv.Atoi(urlPath[len(urlPath)-1])
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing action index: %v", err), http.StatusBadRequest)
		return
	}

	page, ok := config.RoamerPage(pageId)
	if !ok {
		http.Error(w, fmt.Sprintf("roamer-page '%s' not found", pageId), http.StatusNotFound)
		return
	}

	if len(page.Actions) <= actionIdx {
		http.Error(w, fmt.Sprintf("roamer-page '%s' not found", pageId), http.StatusNotFound)
		return
	}

	seq.EnqueueSequence(func() []sequencer.Elem {
		seq, err := script2.Parse(page.Actions[actionIdx].Sequence)
		if err != nil {
			panic(err)
		}

		return seq
	})
}

func hPause(w http.ResponseWriter, _ *http.Request) {
	err := seq.Pause()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(fmt.Sprintf("Server Error: %v", err.Error())))
		if err != nil {
			log.Printf("cannot write error: %v", err)
		}
	}
}

func hAbort(_ http.ResponseWriter, _ *http.Request) {
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
