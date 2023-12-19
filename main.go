package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/imrishuroy/read-cache/api"
	db "github.com/imrishuroy/read-cache/db/sqlc"
)

func main() {

	fmt.Print("Welcome to ReadCache")

	connPool, err := pgxpool.New(context.Background(), "postgres://root:IWSIWDF2024@localhost:5432/read_cache_db?sslmode=disable")

	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}

	store := db.NewStore(connPool)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal().Msg("cannot create server:")
	}

	err = server.Start("localhost:8080")
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

}
