package server

import (
	"embed"
	"log"
	"net/http"
	"rust-roamer/sequencer"
	"strconv"
)

//go:embed html/index.html
var content embed.FS

var seq *sequencer.Sequencer

const port = 10982

func Start() {
	http.Handle("/html/", restrictMethod(http.FileServer(http.FS(content)), http.MethodGet))
	http.Handle("/set/", restrictMethod(http.HandlerFunc(hSet), http.MethodPost))
	http.Handle("/pause", restrictMethod(http.HandlerFunc(hStop), http.MethodPost))
	http.Handle("/state", restrictMethod(http.HandlerFunc(hState), http.MethodGet))

	seq = sequencer.New(1)

	log.Printf("Starting Rust-Roamer")
	log.Printf("http://127.0.0.1:%v/html/", port)

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func restrictMethod(h http.Handler, allowedMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, allowedMethod := range allowedMethods {
			if r.Method == allowedMethod {
				h.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}
