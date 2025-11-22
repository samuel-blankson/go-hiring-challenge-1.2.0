package models

import (
	"gorm.io/gorm"
)

type ProductFilter struct {
	CategoryCode  string  // e.g., "CLOTHING"
	PriceLessThan float64 // products with price < PriceLessThan
}

type IProductsRepository interface {
	GetAllProducts(offset, limit int, filter ProductFilter) ([]Product, int64, error)
	GetProductByCode(code string) (*Product, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

// GetAllProducts supports offset pagination and preloads Category and Variants
func (r *ProductsRepository) GetAllProducts(offset, limit int, filter ProductFilter) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{}).Preload("Variants").Preload("Category")

	// Apply filters
	if filter.CategoryCode != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", filter.CategoryCode)
	}

	if filter.PriceLessThan > 0 {
		query = query.Where("price < ?", filter.PriceLessThan)
	}

	// Count total products for this filter
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductsRepository) GetProductByCode(code string) (*Product, error) {
	var product Product

	err := r.db.Preload("Variants").
		Preload("Category").
		Where("code = ?", code).
		First(&product).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}
