package main

import (
	"fmt"
	"github.com/Oppodelldog/roamer/internal/audioctl"
	"github.com/moutend/go-wca/pkg/wca"
	"log"
	"os"
	"os/signal"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	if dev, err := sound.DefaultDevice(); err != nil {
		log.Fatal(err)
	} else {
		defer dev.Release()
		for _, session := range dev.Sessions {
			name, err := session.GetDisplayNameEnhanced()
			if err != nil {
				log.Fatal(err)
			}
			if name == "Chrome" {
				m, err := session.GetMute()
				if err != nil {
					log.Fatal(err)
				}
				err = session.SetMute(!m)
				if err != nil {
					log.Fatal(err)
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
			log.Fatal(err)
		}

		<-quit
	}
	return
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
