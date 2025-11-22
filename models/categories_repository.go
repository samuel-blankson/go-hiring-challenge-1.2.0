package models

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrDuplicateCategory = errors.New("category with this code already exists")

type ICategoriesRepository interface {
	GetAllCategories() ([]Category, error)
	CreateCategory(category *Category) error
}

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{db: db}
}

// GetAllCategories fetch's all categories in the database.
func (r *CategoriesRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepository) CreateCategory(category *Category) error {
	err := r.db.Create(category).Error
	if err != nil {

		// detect duplicate key error from PostgreSQL
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique key violation
				return ErrDuplicateCategory
			}
		}

		return err
	}

	return nil
}
