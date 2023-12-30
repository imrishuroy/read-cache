package main

import (
	"context"
	"fmt"

	"github.com/imrishuroy/read-cache-api/api"
	"github.com/imrishuroy/read-cache-api/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	db "github.com/imrishuroy/read-cache-api/db/sqlc"
	"github.com/imrishuroy/read-cache-api/util"

	firebase "firebase.google.com/go/v4"

	"google.golang.org/api/option"
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

	opt := option.WithCredentialsFile("./service-account-key.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		// log.Fatal().Msg("Failed to create Firebase app : %v", err)
		log.Fatal().Msg("Failed to create Firebase app")
	}

	fmt.Println("fb connection done ", app)

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Firebase auth client")
	}

	authService := &auth.AuthService{
		FireAuth: authClient,
	}

	// store creation
	store := db.NewStore(connPool)

	// api server setup
	server, err := api.NewServer(config, store, authService)
	if err != nil {
		log.Fatal().Msg("cannot create server:")
	}

	// start the server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

}
