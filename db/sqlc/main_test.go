package db

import (
	"database/sql"
	"github.com/traivok/go-study/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("Could not load configuration:", err)
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Could not connect to db:", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
