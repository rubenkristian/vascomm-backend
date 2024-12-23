package services

import (
	"errors"
	"strings"

	"github.com/rubenkristian/backend/internal/models"
	"gorm.io/gorm"
)

type ProductService struct {
	DB *gorm.DB
}

func InitializeProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		DB: db,
	}
}

func (ps *ProductService) GetAllProduct(take, skip int, search string) ([]models.Product, error) {
	var products []models.Product
	query := ps.DB.Model(&models.Product{}).Limit(take).Offset(skip)

	trimSearch := strings.TrimSpace(search)

	if trimSearch != "" {
		query = query.Where("name LIKE ?", "%"+trimSearch+"%")
	}

	err := query.Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductService) CreateProduct(product *models.Product) error {
	return ps.DB.Create(product).Error
}

func (ps *ProductService) UpdateProduct(id uint, input *models.Product) (*models.Product, error) {
	var product models.Product

	if err := ps.DB.First(&product, id).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if input.Name != "" {
		product.Name = input.Name
	}

	if input.Description != "" {
		product.Description = input.Description
	}

	if input.Price != 0 {
		product.Price = input.Price
	}

	if err := ps.DB.Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (ps *ProductService) DeleteProduct(id uint) error {
	return ps.DB.Delete(&models.Product{}, id).Error
}
