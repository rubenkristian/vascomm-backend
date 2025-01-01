package seeders

import (
	"github.com/rubenkristian/backend/internal/models"
	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB
}

func SeedDB(db *gorm.DB) *Seeder {
	return &Seeder{
		db: db,
	}
}

func (seeder *Seeder) RunSeeder(fresh bool) {
	if fresh {
		seeder.db.Migrator().DropTable(&models.User{}, &models.Product{})
		seeder.db.AutoMigrate(&models.User{}, &models.Product{})
	}
	seeder.productSeeder()
	seeder.userSeeder()
}
