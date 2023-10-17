package action

import (
	"github.com/Oppodelldog/roamer/internal/logger"
)

func StartLoggerWorker(_ <-chan Action, broadcast chan<- []byte) {
	var buffer = make(chan string, 10000)

	logger.AddHook(func(s string) { buffer <- s })

	go func() {
		for s := range buffer {
			broadcast <- msgLogMessage(s)
		}
	}()
}
