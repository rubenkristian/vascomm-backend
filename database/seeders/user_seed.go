package seeders

import (
	"log"

	"github.com/rubenkristian/backend/internal/models"
)

func (seeder *Seeder) userSeeder() {
	users := []models.User{}

	for _, user := range users {
		if err := seeder.db.Create(&user).Error; err != nil {
			log.Println("Failed to seed user:", err)
			return
		}
	}

	log.Println("Database seeded")
}
