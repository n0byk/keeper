package main

import (
	"log"
	"os"
	"os/signal"
	"storage_service/distribution"
	"storage_service/pgsql"
)

func main() {
	// create a thread-safe database instance for use with the router
	log.Println("connecting to database")
	db := pgsql.NewDatabase()

	db.AddUser("asdasd", "asdasd")
	log.Println("connecting to distribution")
	nc := distribution.NewConnection()

	nc.InitSubscribe("keeper", "keeper")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Printf("Draining...")
	db.Close()
	nc.Drain()
	log.Fatalf("Exiting")
}
