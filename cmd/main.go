package main

import (
	"rust-roamer/config"
	"rust-roamer/server"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	server.Start()
}
