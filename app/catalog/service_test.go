package catalog

import (
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const catalogDatabaseUnavailable = "database unavailable"
const productCodeOne = "PROD001"
const productCodeTwo = "PROD002"

type fakeIProductRepository struct {
	products       []models.Product
	total          int64
	product        *models.Product
	err            error
	receivedOffset int
	receivedLimit  int
}

func (f *fakeIProductRepository) ListProducts(offset, limit int) (models.ProductPage, error) {
	f.receivedOffset = offset
	f.receivedLimit = limit
	if f.err != nil {
		return models.ProductPage{}, f.err
	}
	return models.ProductPage{Products: f.products, Total: f.total}, nil
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
				{Code: productCodeOne, Price: decimal.NewFromFloat(99.99)},
				{Code: productCodeTwo, Price: decimal.NewFromFloat(120.00)},
			},
			total: 8,
		}
		service := NewCatalogService(repo)

		products, total, err := service.ListProducts(2, 5)

		assert.NoError(t, err)
		assert.Equal(t, int64(8), total)
		assert.Len(t, products, 2)
		assert.Equal(t, 2, repo.receivedOffset)
		assert.Equal(t, 5, repo.receivedLimit)
		assert.Equal(t, productCodeOne, products[0].Code)
		assert.Equal(t, 99.99, products[0].Price)
		assert.Equal(t, productCodeTwo, products[1].Code)
		assert.Equal(t, 120.0, products[1].Price)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := &fakeIProductRepository{err: errors.New(catalogDatabaseUnavailable)}
		service := NewCatalogService(repo)

		products, total, err := service.ListProducts(0, 10)

		assert.Nil(t, products)
		assert.Equal(t, int64(0), total)
		assert.EqualError(t, err, catalogDatabaseUnavailable)
	})
}

func TestCatalogServiceGetProductDetail(t *testing.T) {
	t.Run("maps product details and inherits missing variant prices", func(t *testing.T) {
		repo := &fakeIProductRepository{
			product: &models.Product{
				Code:  productCodeOne,
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

		product, err := service.GetProductDetail(productCodeOne)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, productCodeOne, product.Code)
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
		repo := &fakeIProductRepository{err: errors.New(catalogDatabaseUnavailable)}
		service := NewCatalogService(repo)

		product, err := service.GetProductDetail(productCodeOne)

		assert.Nil(t, product)
		assert.EqualError(t, err, catalogDatabaseUnavailable)
	})
}
