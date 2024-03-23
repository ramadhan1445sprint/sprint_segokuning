package main

import (
	"log"

	"github.com/ramadhan1445sprint/sprint_segokuning/config"
	database "github.com/ramadhan1445sprint/sprint_segokuning/db"
	"github.com/ramadhan1445sprint/sprint_segokuning/server"
)

func main() {
	config.LoadConfig(".env")

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(80)

	s := server.NewServer(db)
	s.RegisterRoute()

	log.Fatal(s.Run())
}
