package server

import (
	"io"
	"net/http"
)

func serveIndexPage(writer http.ResponseWriter, _ *http.Request) {
	f, err := contentFS().Open("html/index.html")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(writer, f)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
