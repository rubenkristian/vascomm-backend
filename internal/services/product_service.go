package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rubenkristian/backend/commons"
	"github.com/rubenkristian/backend/internal/models"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func InitializeProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		db: db,
	}
}

func (ps *ProductService) GetProduct(id uint) (*models.Product, error) {
	var product models.Product

	if err := ps.db.Find(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product with id %d not found", id)
	}

	return &product, nil
}

func (ps *ProductService) GetAllProduct(pagination *commons.PaginationParams) ([]models.Product, error) {
	pagination.SetParams(10, "asc", "id")
	var products []models.Product

	query := ps.db.Model(&models.Product{}).Limit(pagination.Take).Offset(pagination.Skip)

	trimSearch := strings.TrimSpace(pagination.Search)

	if trimSearch != "" {
		query = query.Where("name LIKE ?", "%"+trimSearch+"%")
	}

	err := query.Order(pagination.SortBy + " " + pagination.Sort).Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (ps *ProductService) CreateProduct(product *models.Product) error {
	return ps.db.Create(product).Error
}

func (ps *ProductService) UpdateProduct(id uint, input *models.Product) (*models.Product, error) {
	var product models.Product

	if err := ps.db.First(&product, id).Error; err != nil {
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

	if err := ps.db.Save(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (ps *ProductService) DeleteProduct(id uint) error {
	return ps.db.Delete(&models.Product{}, id).Error
}
