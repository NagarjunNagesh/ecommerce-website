package catalog

import (
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type fakeIProductRepository struct {
	products []models.Product
	product  *models.Product
	err      error
}

func (f *fakeIProductRepository) GetAllProducts() ([]models.Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.products, nil
}

func (f *fakeIProductRepository) GetProductByCode(code string) (*models.Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.product == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return f.product, nil
}

func TestCatalogServiceListProducts(t *testing.T) {
	t.Run("maps repository products to catalog products", func(t *testing.T) {
		repo := &fakeIProductRepository{
			products: []models.Product{
				{Code: "PROD001", Price: decimal.NewFromFloat(99.99)},
				{Code: "PROD002", Price: decimal.NewFromFloat(120.00)},
			},
		}
		service := NewCatalogService(repo)

		products, err := service.ListProducts()

		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Equal(t, "PROD001", products[0].Code)
		assert.Equal(t, 99.99, products[0].Price)
		assert.Equal(t, "PROD002", products[1].Code)
		assert.Equal(t, 120.0, products[1].Price)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := &fakeIProductRepository{err: errors.New("database unavailable")}
		service := NewCatalogService(repo)

		products, err := service.ListProducts()

		assert.Nil(t, products)
		assert.EqualError(t, err, "database unavailable")
	})
}

func TestCatalogServiceGetProductDetail(t *testing.T) {
	t.Run("maps product details and inherits missing variant prices", func(t *testing.T) {
		repo := &fakeIProductRepository{
			product: &models.Product{
				Code:  "PROD001",
				Price: decimal.NewFromFloat(99.99),
				Category: &models.Category{
					Code: "clothing",
					Name: "Clothing",
				},
				Variants: []models.Variant{
					{Name: "Variant A", SKU: "SKU001A", Price: decimal.NewFromFloat(89.99)},
					{Name: "Variant B", SKU: "SKU001B"},
				},
			},
		}
		service := NewCatalogService(repo)

		product, err := service.GetProductDetail("PROD001")

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, "PROD001", product.Code)
		assert.Equal(t, 99.99, product.Price)
		assert.NotNil(t, product.Category)
		assert.Equal(t, "clothing", product.Category.Code)
		assert.Len(t, product.Variants, 2)
		assert.Equal(t, 89.99, product.Variants[0].Price)
		assert.Equal(t, 99.99, product.Variants[1].Price)
	})

	t.Run("returns product not found when repository has no match", func(t *testing.T) {
		repo := &fakeIProductRepository{}
		service := NewCatalogService(repo)

		product, err := service.GetProductDetail("UNKNOWN")

		assert.Nil(t, product)
		assert.ErrorIs(t, err, ErrProductNotFound)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := &fakeIProductRepository{err: errors.New("database unavailable")}
		service := NewCatalogService(repo)

		product, err := service.GetProductDetail("PROD001")

		assert.Nil(t, product)
		assert.EqualError(t, err, "database unavailable")
	})
}
