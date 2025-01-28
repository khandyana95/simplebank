package main

import (
	"database/sql"
	"log"

	"github.com/khandyan95/simplebank/api"
	db "github.com/khandyan95/simplebank/db/sqlc"
	"github.com/khandyan95/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfing(".")
	if err != nil {
		log.Fatal("cannot load config file")
	}

	dbConn, err := sql.Open(config.DBDriver, config.DataSource)

	if err != nil {
		log.Fatal("error in connecting to DB")
		return
	}

	server := api.NewServer(db.NewStore(dbConn))

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("error in running server")
		return
	}

}
