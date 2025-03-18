package db

import (
	"context"
	"log"
	"os"

	// "os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/singhJasvinder101/go_bank/utils"
)

var testQueries *Queries
var testDB *pgxpool.Pool


func TestMain(m *testing.M) {
    env_config, err := utils.LoadConfig("../../")
    if err != nil {
        log.Fatal("cannot load config: ", err)
    }

    config, err := pgxpool.ParseConfig(env_config.DB_SOURCE)
    if err != nil {
        log.Fatal("cannot parse db config: ", err)
    }

    // pgxpool to create a connection pool (compatible with sqlc)
    testDB, err = pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatal("cannot connect to db: ", err)
    }

    testQueries = New(testDB)
    os.Exit(m.Run()) 
}