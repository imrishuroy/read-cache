package main

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/imrishuroy/read-cache/api"
)

func main() {

	fmt.Print("Welcome to ReadCache")

	server, err := api.NewServer()
	if err != nil {
		log.Fatal().Msg("cannot create server:")
	}

	err = server.Start("localhost:8080")
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

}
