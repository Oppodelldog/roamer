package config

import (
	// embedding default sampleRoamerConfig
	_ "embed"
)

const fileNameUICordsConfig = "roamer-ui-coords.json"

type UIGame struct {
	Slots map[string]Slots
}
type UIConfig struct {
	Games map[string]UIGame
}
type UIRoot map[string]UIConfig

//go:embed default-ui-coords.json
var uiCoordsDefault []byte

var uiConf UIRoot

func GetSlots(version, game, kind string) Slots {
	return uiConf[version].Games[game].Slots[kind]
}
