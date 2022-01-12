package action

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/audioctl"
)

func getSoundSettings() SoundSettings {
	var (
		settings      = SoundSettings{}
		soundSessions []SoundSession
	)

	device, err := audioctl.DefaultDevice()
	if err != nil {
		fmt.Println(err)
		return settings
	}

	defer device.Release()

	for i, session := range device.Sessions {
		var (
			err         error
			sessionId   string
			displayName string
			iconPath    string
			volume      float32
		)

		sessionId, err = session.GetSessionInstanceIdentifier()
		if err != nil {
			fmt.Printf("error getting sound session (%v) display name: %v", i, err)
			continue
		}

		displayName, err = session.GetDisplayNameEnhanced()
		if err != nil {
			fmt.Printf("error getting sound session (%v) display name: %v", i, err)
			continue
		}

		iconPath, err = session.GetIconPath()
		if err != nil {
			fmt.Printf("error getting sound session (%v) icon path: %v", i, err)
			continue
		}

		volume, err = session.GetMasterVolume()
		if err != nil {
			fmt.Printf("error getting sound session (%v) volume: %v", i, err)
			continue
		}

		soundSessions = append(soundSessions, SoundSession{
			Id:    sessionId,
			Name:  displayName,
			Icon:  iconPath,
			Value: volume,
		})
	}

	settings.Sessions = soundSessions
	settings.MainSession = getMasterSoundSession(device)

	return settings
}

func getMasterSoundSession(dev *audioctl.Device) SoundSession {
	scalar, err := dev.GetMasterVolumeLevelScalar()
	if err != nil {
		fmt.Printf("cannot get main volume: %v\n", err)
		return SoundSession{}
	}

	return SoundSession{
		Name:  "Main Volume",
		Icon:  "",
		Value: scalar,
	}
}

func setMainSoundVol(value float32) {
	device, err := audioctl.DefaultDevice()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer device.Release()

	err = device.SetMasterVolumeLevelScalar(value)
	if err != nil {
		fmt.Println(err)
	}
}

func setSoundVol(id string, value float32) {
	device, err := audioctl.DefaultDevice()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer device.Release()

	for i, session := range device.Sessions {
		sid, err := session.GetSessionInstanceIdentifier()
		if err != nil {
			fmt.Printf("error getting sound session (%v) id: %v", i, err)

			continue
		}

		if sid == id {
			err = session.SetMasterVolume(value)
			if err != nil {
				fmt.Printf("error getting sound session (%v) id: %v", i, err)
			}

			return
		}
	}
}
