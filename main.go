package main

import (
	"context"

	"github.com/imrishuroy/read-cache-api/api"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	db "github.com/imrishuroy/read-cache-api/db/sqlc"
	"github.com/imrishuroy/read-cache-api/util"
)

func main() {

	log.Info().Msg("Welcome to ReadCache")

	// loading configurations
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	// db connection
	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}
	defer connPool.Close() // close db connection

	// creating a Firebase app instance
	//opt := option.WithCredentialsFile("./service-account-key.json")
	// app, err := firebase.NewApp(context.Background(), nil, opt)
	// if err != nil {

	// 	log.Fatal().Msg("Failed to create Firebase app")
	// }

	// fmt.Println("fb connection done ", app)

	// auth, err := app.Auth(context.Background())
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to create Firebase auth client")
	// }

	// store creation
	store := db.NewStore(connPool)

	// api server setup
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	// start the server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

}
