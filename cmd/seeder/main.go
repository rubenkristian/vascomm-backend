package main

import (
	"flag"
	"log"

	"github.com/rubenkristian/backend/configs"
	"github.com/rubenkristian/backend/database"
)

func main() {
	env := configs.LoadEnv()
	dbConfig := env.LoadDatabaseConfig()
	db, err := database.ConnectDatabase(dbConfig)

	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
		return
	}

	freshSeeder := flag.Bool("fresh", false, "Run fresh database seeder")
	flag.Parse()

	db.Seeder(*freshSeeder)
}
