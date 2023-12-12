package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:IWSIWDF2024@localhost:5432/read_cache_db?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	// conn, err := sql.Open(dbDriver, dbSource)
	// if err != nil {
	// 	log.Fatal("cannot connect to db:", err)
	// }

	// testQueries = New(conn)

	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(connPool)

	os.Exit(m.Run())
}
