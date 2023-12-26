package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/imrishuroy/read-cache/api"
	db "github.com/imrishuroy/read-cache/db/sqlc"
	"github.com/imrishuroy/read-cache/util"
)

func main() {

	log.Info().Msg("Welcome to ReadCache")
	fmt.Print("Welcome to ReadCache")

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}

	store := db.NewStore(connPool)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal().Msg("cannot create server:")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

}
