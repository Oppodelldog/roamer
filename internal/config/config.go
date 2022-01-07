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
	return loadConfig(fileNameRoamerConfig, &config, sampleRoamerConfig)
}

func Save() error {
	return saveConfig(fileNameRoamerConfig, config)
}

func NewPage() error {
	var id = fmt.Sprintf("page_%v", len(config.Pages)+1)
	if _, exists := config.Pages[id]; exists {
		return fmt.Errorf("the page (%s) already exists", id)
	}

	config.Pages[id] = Page{
		TitleShort: "New Page",
		Title:      "New Page",
	}

	return Save()
}

func DeletePage(id string) error {
	if _, exists := config.Pages[id]; !exists {
		return fmt.Errorf("the page (%s) does not exists", id)
	}

	delete(config.Pages, id)

	return Save()
}

func NewSequence(pageId string) error {
	page, exists := config.Pages[pageId]
	if !exists {
		return fmt.Errorf("the page (%s) does not exists", pageId)
	}

	page.Actions = append(page.Actions, Action{Caption: "New Macro"})
	config.Pages[pageId] = page

	return Save()
}

func DeleteSequence(pageId string, seq int) error {
	page, exists := config.Pages[pageId]
	if !exists {
		return fmt.Errorf("the page (%s) does not exists", pageId)
	}

	page.Actions = append(page.Actions[:seq], page.Actions[seq+1:]...)
	config.Pages[pageId] = page

	return Save()
}

func SavePages(pages Pages) error {
	config.Pages = pages

	return Save()
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

func saveConfig(filename string, c Config) error {
	ensureConfig(filename, nil)

	f, err := os.Create(getConfigFilePath(filename))
	if err != nil {
		return fmt.Errorf("cannot open config file for writing: %w", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("cannot close '%s': %v\n", filename, err)
		}
	}()

	err = json.NewEncoder(f).Encode(c)
	if err != nil {
		return fmt.Errorf("cannot encode '%s': %w", filename, err)
	}

	return nil
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
