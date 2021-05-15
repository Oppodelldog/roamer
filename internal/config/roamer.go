package config

import _ "embed"

const fileNameRoamerConfig = "roamer-config.json"

//go:embed sample-roamer-config.json
var sampleRoamerConfig []byte

var config Config

type (
	Config struct {
		Title       string `json:"Title"`
		WelcomeText string `json:"WelcomeText"`
		Pages       Pages  `json:"Pages"`
	}
	Pages map[string]Page
	Page  struct {
		TitleShort string   `json:"TitleShort"`
		Title      string   `json:"Title"`
		CSSFile    string   `json:"CssFile"`
		Columns    []int    `json:"Columns"`
		Actions    []Action `json:"Actions"`
	}
	Action struct {
		Icon     string `json:"Icon"`
		Action   string `json:"Action"`
		Caption  string `json:"Caption"`
		Sequence string `json:"Sequence"`
	}
)

func Roamer() Config {
	return config
}

func RoamerPage(pageId string) (Page, bool) {
	for id, page := range config.Pages {
		if id != pageId {
			continue
		}
		return page, true
	}

	return Page{}, false
}
