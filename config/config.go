package config

import (
	// embedding default config
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Game struct {
	Slots map[string]Slots
}
type Config struct {
	Games map[string]Game
}
type Root map[string]Config

const fileName = "roamer-config.json"
const filePerm = 0600

//go:embed default.json
var defaultConfig []byte

var conf Root

func Load() error {
	ensureConfig()

	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("cannot load config: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("cannot close config: %v\n", err)
		}
	}()

	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		return fmt.Errorf("cannot decode config: %w", err)
	}

	return nil
}

func ensureConfig() {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err := ioutil.WriteFile(fileName, defaultConfig, filePerm)
		if err != nil {
			fmt.Printf("cannot write default config: %v", err)
		}
	}
}

func Save() {
	f, err := os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE, filePerm)
	if err != nil {
		fmt.Printf("cannot open config file for writing: %v", err)
		return
	}

	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("cannot close config: %v", err)
		}
	}()

	err = json.NewEncoder(f).Encode(conf)
	if err != nil {
		fmt.Printf("cannot encode config: %v", err)
	}
}

func GetSlots(version, game, kind string) Slots {
	return conf[version].Games[game].Slots[kind]
}
