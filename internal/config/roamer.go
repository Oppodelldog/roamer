package config

import (
	_ "embed"
	"fmt"
)

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

func SetSequence(pageId string, sequenceIndex int, sequence string) error {
	var page, pageFound = config.Pages[pageId]
	if !pageFound {
		return fmt.Errorf("page %s not found", pageId)
	}

	if len(page.Actions) < sequenceIndex {
		return fmt.Errorf("sequence %v not found for page %s", sequenceIndex, pageId)
	}

	var action = page.Actions[sequenceIndex]

	action.Sequence = sequence

	config.Pages[pageId].Actions[sequenceIndex] = action

	return nil
}
