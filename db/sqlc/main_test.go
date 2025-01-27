package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver   = "postgres"
	datasource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testStore *Store

func TestMain(m *testing.M) {

	conn, err := sql.Open(dbDriver, datasource)

	if err != nil {
		log.Fatal("error in connecting to db")
	}

	testQueries = New(conn)
	testStore = NewStore(conn)

	os.Exit(m.Run())
}
