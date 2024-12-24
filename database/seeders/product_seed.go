package seeders

import (
	"log"

	"github.com/rubenkristian/backend/internal/models"
)

func (seeder *Seeder) productSeeder() {
	products := []models.Product{{
		Name:        "Product 1",
		Description: "This is product 1",
		Price:       18000,
	}, {
		Name:        "Product 2",
		Description: "This is product 2",
		Price:       20000,
	}, {
		Name:        "Product 3",
		Description: "This is product 3",
		Price:       38000,
	}}

	for _, product := range products {
		if err := seeder.db.Create(&product).Error; err != nil {
			log.Println("Failed to seed user:", err)
			return
		}
	}

	log.Println("Database table product seeded")
}
