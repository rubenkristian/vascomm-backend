package seeders

import (
	"log"

	"github.com/rubenkristian/backend/internal/models"
	"github.com/rubenkristian/backend/utils"
)

func (seeder *Seeder) userSeeder() {
	password, err := utils.HashPasssword("password")

	if err != nil {
		log.Println("Failed to seed user:", err)
		return
	}

	users := []models.User{{
		Name:        "user",
		PhoneNumber: "62800000000",
		Email:       "user@gmail.com",
		Password:    password,
		Role:        "user",
	}, {
		Name:        "admin",
		PhoneNumber: "62800000001",
		Email:       "admin@gmail.com",
		Password:    password,
		Role:        "admin",
	}}

	for _, user := range users {
		if err := seeder.db.Create(&user).Error; err != nil {
			log.Println("Failed to seed user:", err)
			return
		}
	}

	log.Println("Database user table seeded")
}
