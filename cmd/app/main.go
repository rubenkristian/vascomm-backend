package main

import (
	"log"

	"github.com/rubenkristian/backend/configs"
	"github.com/rubenkristian/backend/database"
	"github.com/rubenkristian/backend/internal/app"
	"github.com/rubenkristian/backend/utils"
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

	emailer, err := utils.InitializeEmailer(emailConfig)

	if err != nil {
		log.Fatalf("Failed to initialize email %v", err)
		return
	}

	app.CreateApp(db.Conn, emailer, env)
}
