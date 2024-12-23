package services

import (
	"errors"
	"strings"

	"github.com/rubenkristian/backend/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func InitializeUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (ps *UserService) GetAllUser(take, skip int, search string) ([]models.User, error) {
	var users []models.User

	query := ps.DB.Model(&models.User{}).Limit(take).Offset(skip)

	trimSearch := strings.TrimSpace(search)

	if trimSearch != "" {
		query = query.Where("name LIKE ?", "%"+trimSearch+"%")
	}

	err := query.Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (ps *UserService) CreateUser(product *models.User) error {
	return ps.DB.Create(product).Error
}

func (ps *UserService) UpdateUser(id uint, input *models.User) (*models.User, error) {
	var user models.User

	if err := ps.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	if input.PhoneNumber != "" {
		user.PhoneNumber = input.PhoneNumber
	}

	if input.Password != "" {
		user.Password = input.Password
	}

	if err := ps.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (ps *UserService) DeleteUser(id uint) error {
	return ps.DB.Delete(&models.User{}, id).Error
}
