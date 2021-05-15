package server

import (
	"github.com/Oppodelldog/roamer/internal/pages/index"
	"github.com/Oppodelldog/roamer/internal/pages/roamerpage"
	"net/http"
	"path"
)

func serveIndexPage(writer http.ResponseWriter, request *http.Request) {
	err := index.Render(content, writer)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func serveRoamerPage(writer http.ResponseWriter, request *http.Request) {
	pageId := path.Base(request.URL.Path)
	err := roamerpage.Render(pageId, content, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
