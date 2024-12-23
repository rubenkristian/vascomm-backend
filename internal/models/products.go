package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:200;not null" json:"name" validate:"required,min=3"`
	Description string         `gorm:"size:500;not null" json:"description" validate:"omitempty,min=10"`
	Price       float64        `gorm:"not null" json:"price" validate:"required,gte=0"`
	Image       string         `json:"image"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
