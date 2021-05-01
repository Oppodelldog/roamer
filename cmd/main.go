package main

import (
	config2 "github.com/Oppodelldog/roamer/internal/config"
	server2 "github.com/Oppodelldog/roamer/internal/server"
)

func main() {
	if err := config2.Load(); err != nil {
		panic(err)
	}

	server2.Start()
}
