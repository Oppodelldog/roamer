package server

import (
	"net/http"
	"path"

	"github.com/Oppodelldog/roamer/internal/pages/index"
	"github.com/Oppodelldog/roamer/internal/pages/roamerpage"
)

func serveIndexPage(writer http.ResponseWriter, request *http.Request) {
	var err = index.Render(content, writer)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func serveRoamerPage(writer http.ResponseWriter, request *http.Request) {
	var (
		pageId = path.Base(request.URL.Path)
		err    = roamerpage.Render(pageId, content, writer)
	)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
