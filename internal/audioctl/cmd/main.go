package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	sound "github.com/Oppodelldog/roamer/internal/audioctl"
	"github.com/moutend/go-wca/pkg/wca"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	dev, err := sound.DefaultDevice()
	if err != nil {
		log.Fatal(err)

		return
	}

	defer dev.Release()

	for _, session := range dev.Sessions {
		name, err := session.GetDisplayNameEnhanced()
		if err != nil {
			log.Println(err)

			return
		}

		if name == "Chrome" {
			m, err := session.GetMute()
			if err != nil {
				log.Println(err)

				return
			}

			err = session.SetMute(!m)
			if err != nil {
				log.Println(err)

				return
			}
		}

		fmt.Println(name)
	}

	err = sound.RegisterDeviceEvents(wca.IMMNotificationClientCallback{
		OnDefaultDeviceChanged: onDefaultDeviceChanged,
		OnDeviceAdded:          onDeviceAdded,
		OnDeviceRemoved:        onDeviceRemoved,
		OnDeviceStateChanged:   onDeviceStateChanged,
		OnPropertyValueChanged: onPropertyValueChanged,
	})

	if err != nil {
		log.Println(err)

		return
	}

	<-quit
}

func onDefaultDeviceChanged(flow wca.EDataFlow, role wca.ERole, pwstrDeviceId string) error {
	fmt.Printf("Called OnDefaultDeviceChanged\t(%v, %v, %q)\n", flow, role, pwstrDeviceId)

	return nil
}

func onDeviceAdded(pwstrDeviceId string) error {
	fmt.Printf("Called OnDeviceAdded\t(%q)\n", pwstrDeviceId)

	return nil
}

func onDeviceRemoved(pwstrDeviceId string) error {
	fmt.Printf("Called OnDeviceRemoved\t(%q)\n", pwstrDeviceId)

	return nil
}

func onDeviceStateChanged(pwstrDeviceId string, dwNewState uint64) error {
	fmt.Printf("Called OnDeviceStateChanged\t(%q, %v)\n", pwstrDeviceId, dwNewState)

	return nil
}

func onPropertyValueChanged(pwstrDeviceId string, key uint64) error {
	fmt.Printf("Called OnPropertyValueChanged\t(%q, %v)\n", pwstrDeviceId, key)

	return nil
}
