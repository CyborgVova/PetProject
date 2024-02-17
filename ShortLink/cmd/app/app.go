package main

import (
	"flag"
	"shortlink/database"
	"shortlink/service"
)

func main() {
	var storage string
	flag.StringVar(&storage, "s", "inmemory", "-s inmemory or -s database (default: inmemory)")
	flag.Parse()
	db := database.CreateDB(storage)
	go service.RunRest()
	service.RunGrpc(storage, db)
}
