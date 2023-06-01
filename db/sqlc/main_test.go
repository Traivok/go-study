package db

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:15432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	connection, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Could not connect to db:", err)
	}

	testQueries = New(connection)

	m.Run()
}