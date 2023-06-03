package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/traivok/go-study/api"
	db "github.com/traivok/go-study/db/sqlc"
	"github.com/traivok/go-study/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load configuration:", config)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Could not connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Could not start server:", err)

	}
}
