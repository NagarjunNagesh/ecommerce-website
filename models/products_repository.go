package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductListOptions struct {
	Offset        int
	Limit         int
	CategoryCode  string
	PriceLessThan *decimal.Decimal
}

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

func (r *ProductsRepository) ListProducts(options ProductListOptions) (ProductPage, error) {
	var products []Product
	var total int64

	query := r.applyProductListFilters(r.db.Model(&Product{}), options).Session(&gorm.Session{})

	if err := query.Count(&total).Error; err != nil {
		return ProductPage{}, err
	}

	if err := query.Preload("Category").Preload("Variants").Order("code ASC").Offset(options.Offset).Limit(options.Limit).Find(&products).Error; err != nil {
		return ProductPage{}, err
	}

	return ProductPage{Products: products, Total: total}, nil
}

func (r *ProductsRepository) applyProductListFilters(db *gorm.DB, options ProductListOptions) *gorm.DB {
	if options.CategoryCode != "" {
		db = db.Joins("JOIN categories ON categories.id = products.category_id").Where("categories.code = ?", options.CategoryCode)
	}
	if options.PriceLessThan != nil {
		db = db.Where("products.price < ?", *options.PriceLessThan)
	}
	return db
}

func (r *ProductsRepository) GetProductByCode(code string) (*Product, error) {
	var product Product
	if err := r.db.Preload("Category").Preload("Variants").Where("code = ?", code).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
