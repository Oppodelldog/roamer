package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const filePerm = 0600
const appDataFolderName = "roamer"

func Load() error {
	err := loadConfig(fileNameRoamerConfig, &config, sampleRoamerConfig)
	if err != nil {
		return err
	}

	return loadConfig(fileNameUICordsConfig, &uiConf, uiCoordsDefault)
}

func loadConfig(filename string, data interface{}, defaultData []byte) error {
	ensureConfig(filename, defaultData)

	f, err := os.Open(getConfigFilePath(filename))
	if err != nil {
		return fmt.Errorf("cannot load sampleRoamerConfig: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("cannot close '%s': %v\n", filename, err)
		}
	}()

	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return fmt.Errorf("cannot decode '%s': %w", filename, err)
	}

	return nil
}

func ensureConfig(filename string, defaultData []byte) {
	var (
		appDataFolder = getAppDataFolder()
		err           = os.MkdirAll(appDataFolder, filePerm)
	)
	if err != nil {
		panic(fmt.Sprintf("unable to create app data folder '%s': %v", appDataFolder, err))
	}

	configFilePath := getConfigFilePath(filename)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		err := ioutil.WriteFile(configFilePath, defaultData, filePerm)
		if err != nil {
			fmt.Printf("cannot write default '%s': %v", filename, err)
		}
	}
}

func getConfigFilePath(filename string) string {
	appDataFolder := getAppDataFolder()
	configFilePath := filepath.Join(appDataFolder, filename)

	return configFilePath
}

func getAppDataFolder() string {
	appDataDir, ok := os.LookupEnv("APPDATA")
	if !ok {
		panic("could not get APPDATA variable from env")
	}

	appDataFolder := filepath.Join(appDataDir, appDataFolderName)

	return appDataFolder
}
