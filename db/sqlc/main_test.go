package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/imrishuroy/read-cache/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	// conn, err := sql.Open(dbDriver, dbSource)
	// if err != nil {
	// 	log.Fatal("cannot connect to db:", err)
	// }

	// testQueries = New(conn)

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config")
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(connPool)

	os.Exit(m.Run())
}
