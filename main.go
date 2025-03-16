package main

import (
    "context"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
    server "github.com/singhJasvinder101/go_bank/api"
    db "github.com/singhJasvinder101/go_bank/db/sqlc"
)

const (
    dbSource = "postgresql://postgres:123@localhost:5432/go_bank?sslmode=disable"
    address  = "localhost:3000"
)

func main() {
    config, err := pgxpool.ParseConfig(dbSource)
    if err != nil {
        log.Fatal("cannot parse db config: ", err)
    }

    conn, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatal("cannot connect to db: ", err)
    }

    store := db.NewStore(conn)
    srv := server.NewServer(store)

    err = srv.Start(address)
    if err != nil {
        log.Fatal("cannot start server: ", err)
    }
}