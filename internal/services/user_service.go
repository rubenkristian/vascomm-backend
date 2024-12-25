package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-gomail/gomail"
	"github.com/rubenkristian/backend/internal/models"
	"github.com/rubenkristian/backend/utils"
	"gorm.io/gorm"
)

type UserService struct {
	authGenerator *utils.AuthToken
	db            *gorm.DB
	emailer       *utils.Emailer
}

func InitializeUserService(authGenerator *utils.AuthToken, db *gorm.DB, emailer *utils.Emailer) *UserService {
	return &UserService{
		authGenerator: authGenerator,
		db:            db,
		emailer:       emailer,
	}
}

func (userService *UserService) UserLogin(username, password string) (*models.User, error) {
	var user models.User

	if err := userService.db.Where("phone_number = ? OR email = ? AND role = ?", username, username, "user").First(&user).Error; err != nil {
		return nil, fmt.Errorf("email/phone number not match")
	}

	if err := userService.authGenerator.ValidatePassword(password, user.Password); err != nil {
		return nil, fmt.Errorf("invalid password for this user")
	}

	return &user, nil
}

func (userService *UserService) AdminLogin(username, password string) (*models.User, error) {
	var user models.User

	if err := userService.db.Where("phone_number = ? OR email = ? AND role = ?", username, username, "admin").First(&user).Error; err != nil {
		return nil, fmt.Errorf("email/phone number not match")
	}

	if err := userService.authGenerator.ValidatePassword(password, user.Password); err != nil {
		return nil, fmt.Errorf("invalid password for this admin")
	}

	return &user, nil
}

func (userService *UserService) RegisterFromGoogle(name, email, phoneNumber string) *models.User {
	var user models.User

	err := userService.db.Where("phone_number = ?", phoneNumber).Or("email = ?", email).First(&user).Error

	if err != nil {
		user = models.User{
			Name:        name,
			Email:       email,
			PhoneNumber: phoneNumber,
		}
		userService.db.Create(&user)
	}

	return &user
}

func (us *UserService) GetUser(id uint) (*models.User, error) {
	var user models.User

	if err := us.db.Find(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return &user, nil
}

func (us *UserService) GetAllUser(take, skip int, search string) ([]models.User, error) {
	var users []models.User

	query := us.db.Model(&models.User{}).Limit(take).Offset(skip)

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

func (userService *UserService) CreateUser(userInput *models.User) error {
	password, hash, err := userService.authGenerator.GeneratePassword(12)

	if err != nil {
		return err
	}

	email := gomail.NewMessage()
	email.SetHeader("Subject", "Password Login")
	email.SetBody("text/html", "Here is your password <b>"+password+"</b>")

	if err := userService.emailer.SendEmail(userInput.Email, email); err != nil {
		return err
	}

	userInput.Password = hash
	userInput.Role = "user"

	return userService.db.Create(userInput).Error
}

func (us *UserService) UpdateUser(id uint, userInput *models.User) (*models.User, error) {
	var user models.User

	if err := us.db.First(&user, id).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if userInput.Name != "" {
		user.Name = userInput.Name
	}

	if userInput.Email != "" {
		user.Email = userInput.Email
	}

	if userInput.PhoneNumber != "" {
		user.PhoneNumber = userInput.PhoneNumber
	}

	if userInput.Password != "" {
		user.Password = userInput.Password
	}

	if err := us.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserService) DeleteUser(id uint) error {
	return us.db.Delete(&models.User{}, id).Error
}
