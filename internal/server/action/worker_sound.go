package action

import (
	"fmt"
)

func StartSoundSettingsWorker(actions <-chan Action, _ chan<- []byte) {
	go func() {
		for action := range actions {
			switch v := action.(type) {
			case LoadSoundSettings:
				v.Response <- msgSoundSettings(getSoundSettings())
			case SetSoundVolume:
				setSoundVol(v.Id, v.Value)
			case SetMainSoundVolume:
				setMainSoundVol(v.Value)
			default:
				fmt.Printf("unknown action for sound-settings worker: %T\n", action)
			}
		}
	}()
}
