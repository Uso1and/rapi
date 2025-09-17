package main

import (
	"log"
	"resapi/internal/app/routers"
	"resapi/internal/domain/infs/database"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Error Open Connect DB: %v", err)
	}

	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
	r := routers.PubRoutSetup()

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
