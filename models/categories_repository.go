package models

import (
	"errors"

	"gorm.io/gorm"
)

var ErrCategoryAlreadyExists = errors.New("category code already exists")

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoriesRepository) CreateCategory(category *Category) error {
	if err := r.db.Create(category).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrCategoryAlreadyExists
		}
		return err
	}

	return nil
}
