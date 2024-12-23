package seeders

import (
	"log"

	"github.com/rubenkristian/backend/internal/models"
)

func (seeder *Seeder) productSeeder() {
	products := []models.Product{}

	for _, product := range products {
		if err := seeder.db.Create(&product).Error; err != nil {
			log.Println("Failed to seed user:", err)
			return
		}
	}

	log.Println("Database seeded")
}
