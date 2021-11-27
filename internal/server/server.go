package server

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
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
		actions = make(chan action.Action)
		hub     = ws.StartHub(newSessionFunc(actions))
	)

	action.Worker(actions, hub.Broadcast())

	http.Handle("/", restrictMethod(http.HandlerFunc(serveIndexPage), http.MethodGet))
	http.Handle("/attributions.html", restrictMethod(addPrefix("/html/", http.FileServer(http.FS(content))), http.MethodGet))
	http.Handle("/roam/", restrictMethod(http.StripPrefix("/roam/", http.HandlerFunc(serveRoamerPage)), http.MethodGet))
	http.Handle("/img/", restrictMethod(http.FileServer(http.FS(img)), http.MethodGet))
	http.Handle("/js/", restrictMethod(http.FileServer(http.FS(js)), http.MethodGet))
	http.Handle("/css/", restrictMethod(http.FileServer(http.FS(css)), http.MethodGet))
	http.Handle("/favicon.ico", restrictMethod(addPrefix("/root", http.FileServer(http.FS(root))), http.MethodGet))

	http.Handle("/ws", restrictMethod(websocketHandler(hub), http.MethodGet))

	log.Printf("Starting Roamer")
	log.Printf("http://127.0.0.1:%v", port)

	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("error running http server: %v", err)
	}
}

func newSessionFunc(actions chan action.Action) ws.NewSessionFunc {
	return func(client *ws.Client) {
		action.ClientSession(client, actions)
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
