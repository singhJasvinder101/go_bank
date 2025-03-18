package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	server "github.com/singhJasvinder101/go_bank/api"
	db "github.com/singhJasvinder101/go_bank/db/sqlc"
	"github.com/singhJasvinder101/go_bank/utils"
)

func main() {
    
    env_config, err := utils.LoadConfig([]string{".", "/app"})
    if err != nil {
        log.Fatal("cannot load config: ", err)
    }


    config, err := pgxpool.ParseConfig(env_config.DB_SOURCE)
    if err != nil {
        log.Fatal("cannot parse db config: ", err)
    }

    conn, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatal("cannot connect to db: ", err)
    }

    store := db.NewStore(conn)
    srv := server.NewServer(store)

    err = srv.Start(env_config.ADDRESS)
    if err != nil {
        log.Fatal("cannot start server: ", err)
    }
}