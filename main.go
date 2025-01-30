package main

import (
	"log"

	"github.com/khandyan95/simplebank/api"
)

func main() {

	server, err := api.NewServer()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	if err := server.Start(); err != nil {
		log.Fatal("error in running server")
		return
	}
}
