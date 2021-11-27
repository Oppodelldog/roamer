package main

import (
	"github.com/Oppodelldog/roamer/internal/config"
	"github.com/Oppodelldog/roamer/internal/server"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}
	server.Start()
}
