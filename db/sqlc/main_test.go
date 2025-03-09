package db

import (
    "context"
    "log"
    "os"
    "testing"

    "github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

const dbSource = "postgresql://postgres:123@localhost:5432/go_bank?sslmode=disable"

func TestMain(m *testing.M) {
    config, err := pgxpool.ParseConfig(dbSource)
    if err != nil {
        log.Fatal("cannot parse db config: ", err)
    }

    testDB, err = pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatal("cannot connect to db: ", err)
    }

    testQueries = New(testDB)

    os.Exit(m.Run())
}