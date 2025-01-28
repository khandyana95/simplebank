package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/khandyan95/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testStore *Store

func TestMain(m *testing.M) {

	config, err := util.LoadConfing("../..")
	if err != nil {
		log.Fatal("cannot load config file")
	}

	conn, err := sql.Open(config.DBDriver, config.DataSource)

	if err != nil {
		log.Fatal("error in connecting to db")
	}

	testQueries = New(conn)
	testStore = NewStore(conn)

	os.Exit(m.Run())
}
