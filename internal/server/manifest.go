package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Manifest holds data of the manifest
//{
//  "name": "MyService",
//  "short_name": "MySvc",
//  "theme_color": "rgba(220,136,68,0.13)",
//  "background_color": "#ffffff",
//  "display": "browser",
//  "scope": "/",
//  "start_url": "/"
//  "icons": [
// 		{
// 			"src": "/images/icons-192.png",
// 			"type": "image/png",
// 			"sizes": "192x192"
// 		},
// 		{
// 			"src": "/images/icons-512.png",
// 			"type": "image/png",
// 			"sizes": "512x512"
// 		}
//  ],
//}

type Manifest struct {
	Name            string         `json:"name,omitempty"`
	ShortName       string         `json:"short_name,omitempty"`
	ThemeColor      string         `json:"theme_color,omitempty"`
	BackgroundColor string         `json:"background_color,omitempty"`
	Display         string         `json:"display,omitempty"`
	Scope           string         `json:"scope,omitempty"`
	StartUrl        string         `json:"start_url,omitempty"`
	Icons           []ManifestIcon `json:"icons,omitempty"`
}

type ManifestIcon struct {
	Src   string `json:"src,omitempty"`
	Type  string `json:"type,omitempty"`
	Sizes string `json:"sizes,omitempty"`
}

func ManifestHandler(m Manifest) http.Handler {
	var bytes, err = json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("error marshalling Config to json: %v", err))
	}

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("content-type", "application/json")
		_, _ = writer.Write(bytes)
	})
}
