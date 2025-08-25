package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

// var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	connStr := "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	var err error
	testDB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	//testQueries = New(conn)

	os.Exit(m.Run())
}
