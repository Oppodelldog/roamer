package server

import (
	"embed"
	"errors"
	"fmt"
	key2 "github.com/Oppodelldog/roamer/internal/key"
	sequencer2 "github.com/Oppodelldog/roamer/internal/sequencer"
	"log"
	"net/http"
	"path"
	"strconv"
)

//go:embed html
var content embed.FS

//go:embed img
var img embed.FS

//go:embed js
var js embed.FS

//go:embed css
var css embed.FS

//go:embed root/favicon.ico
var root embed.FS

var seq *sequencer2.Sequencer

const port = 10982

func Start() {
	http.Handle("/", restrictMethod(http.HandlerFunc(serveIndexPage), http.MethodGet))
	http.Handle("/attributions.html", restrictMethod(addPrefix("/html/", http.FileServer(http.FS(content))), http.MethodGet))
	http.Handle("/roam/", restrictMethod(http.StripPrefix("/roam/", http.HandlerFunc(serveRoamerPage)), http.MethodGet))
	http.Handle("/img/", restrictMethod(http.FileServer(http.FS(img)), http.MethodGet))
	http.Handle("/js/", restrictMethod(http.FileServer(http.FS(js)), http.MethodGet))
	http.Handle("/css/", restrictMethod(http.FileServer(http.FS(css)), http.MethodGet))
	http.Handle("/favicon.ico", restrictMethod(addPrefix("/root", http.FileServer(http.FS(root))), http.MethodGet))
	http.Handle("/set/", restrictMethod(http.HandlerFunc(hSet), http.MethodPost))
	http.Handle("/pause", restrictMethod(http.HandlerFunc(hPause), http.MethodPost))
	http.Handle("/abort", restrictMethod(http.HandlerFunc(hAbort), http.MethodPost))
	http.Handle("/state", restrictMethod(http.HandlerFunc(hState), http.MethodGet))

	seq = sequencer2.New(1)
	seq.BeforeSequence(func() {
		key2.ResetPressed()
	})

	log.Printf("Starting Rust-Roamer")
	log.Printf("http://127.0.0.1:%v", port)

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("error running http server: %v", err)
	}
}

func addPrefix(s string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		request.URL.Path = path.Join(s, request.URL.Path)
		h.ServeHTTP(writer, request)
	})
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
