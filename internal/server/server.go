package server

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/Oppodelldog/roamer/internal/server/action"
	"github.com/Oppodelldog/roamer/internal/server/ws"
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

const port = 10982

func Start() {
	var (
		actions          = make(chan action.Action)
		sequencerActions = make(chan action.Action)
		soundActions     = make(chan action.Action)
		loggerActions    = make(chan action.Action)
		hub              = ws.StartHub(newSessionFunc(actions, sequencerActions, soundActions, loggerActions))
	)

	action.StartConfigWorker(actions, hub.Broadcast())
	action.StartSequencerWorker(sequencerActions, hub.Broadcast())
	action.StartSoundSettingsWorker(soundActions, hub.Broadcast())
	action.StartLoggerWorker(loggerActions, hub.Broadcast())

	http.Handle("/", restrictMethod(http.HandlerFunc(serveIndexPage), http.MethodGet))
	http.Handle("/attributions.html", restrictMethod(addPrefix("/html/", http.FileServer(http.FS(contentFS()))), http.MethodGet))
	http.Handle("/img/", restrictMethod(http.FileServer(http.FS(imgFS())), http.MethodGet))
	http.Handle("/js/", restrictMethod(http.FileServer(http.FS(jsFS())), http.MethodGet))
	http.Handle("/css/", restrictMethod(http.FileServer(http.FS(cssFS())), http.MethodGet))
	http.Handle("/favicon.ico", restrictMethod(addPrefix("/root", http.FileServer(http.FS(root))), http.MethodGet))
	http.Handle("/manifest.json", restrictMethod(ManifestHandler(Manifest{
		Name:            "Roamer",
		ShortName:       "Roamer",
		ThemeColor:      "#ffffff",
		BackgroundColor: "#000000",
		Display:         "browser",
		Scope:           "/",
		StartUrl:        "/",
	}), http.MethodGet))

	http.Handle("/ws", restrictMethod(websocketHandler(hub), http.MethodGet))

	log.Printf("Starting Roamer")
	log.Printf("http://127.0.0.1:%v", port)

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil) //nolint:gosec
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("error running http server: %v", err)
	}
}

func newSessionFunc(actions, sequencerActions, soundActions, loggerActions chan action.Action) ws.NewSessionFunc {
	return func(client *ws.Client) {
		action.ClientSession(client, actions, sequencerActions, soundActions, loggerActions)
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

func websocketHandler(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	}
}

func filesystem(envVar string, fs fs.FS) fs.FS {
	if absolutePath, ok := os.LookupEnv(envVar); ok {
		return os.DirFS(absolutePath)
	}

	return fs
}

func contentFS() fs.FS {
	return assetFS(content)
}

func cssFS() fs.FS {
	return assetFS(css)
}

func jsFS() fs.FS {
	return assetFS(js)
}

func imgFS() fs.FS {
	return assetFS(img)
}

func assetFS(efs embed.FS) fs.FS {
	return filesystem("ROAMER_ASSETS", efs)
}
