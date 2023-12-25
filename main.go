package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"read-cache/api"
	db "read-cache/db/sqlc"
)

func main() {

	log.Info().Msg("Welcome to ReadCache")
	fmt.Print("Welcome to ReadCache")

	// config, err := util.LoadConfig(".")
	// if err != nil {
	// 	log.Fatal().Msg("cannot load config")
	// }

	//connPool, err := pgxpool.New(context.Background(), config.DBSource)
	connPool, err := pgxpool.New(context.Background(), "postgres://postgres:7EQERkvXwFYUcdidVxUd@read-cache.cf48iqcewxbw.ap-south-1.rds.amazonaws.com:5432/read_cache_db")
	//
	// dbUrl := os.Getenv("DATABASE_URL")
	//fmt.Println(dbUrl)
	//connPool, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}

	store := db.NewStore(connPool)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal().Msg("cannot create server:")
	}

	err = server.Start(":8080")
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

}
