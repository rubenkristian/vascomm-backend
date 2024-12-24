package main

import (
	"flag"
	"log"

	"github.com/rubenkristian/backend/configs"
	"github.com/rubenkristian/backend/database"
	"github.com/rubenkristian/backend/internal/app"
	"github.com/rubenkristian/backend/pkg"
)

func main() {
	env := configs.LoadEnv()
	dbConfig := env.LoadDatabaseConfig()
	emailConfig := env.LoadEmailConfig()
	db, err := database.ConnectDatabase(dbConfig)

	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
		return
	}

	emailer, err := pkg.InitializeEmailer(emailConfig)

	if err != nil {
		log.Fatalf("Failed to initialize email %v", err)
		return
	}

	runSeeder := flag.Bool("seed", false, "Run database seeder")
	freshSeeder := flag.Bool("fresh", false, "Run fresh database seeder")
	flag.Parse()

	if *runSeeder {
		db.Seeder(*freshSeeder)
	}

	db.Migrate()

	app.CreateApp(db.DB, emailer, env)
}
