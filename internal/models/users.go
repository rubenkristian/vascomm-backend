package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name" validate:"required,min=2"`
	PhoneNumber string         `gorm:"size:16;unique" json:"phone_number" validate:"required,min=10"`
	Email       string         `gorm:"size:100;unique" json:"email" validate:"required,email"`
	Password    string         `json:"-" validate:"omitempty,min=6"`
	Role        string         `json:"role"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
