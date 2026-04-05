package models

import (
	"gorm.io/gorm"
)

type ProductPage struct {
	Products []Product
	Total    int64
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) ListProducts(offset, limit int) (ProductPage, error) {
	var products []Product
	var total int64

	if err := r.db.Model(&Product{}).Count(&total).Error; err != nil {
		return ProductPage{}, err
	}

	if err := r.db.Preload("Category").Preload("Variants").Order("code ASC").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return ProductPage{}, err
	}

	return ProductPage{Products: products, Total: total}, nil
}

func (r *ProductsRepository) GetProductByCode(code string) (*Product, error) {
	var product Product
	if err := r.db.Preload("Category").Preload("Variants").Where("code = ?", code).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
