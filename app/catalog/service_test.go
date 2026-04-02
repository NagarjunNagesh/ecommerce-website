package catalog

import (
    "errors"
    "testing"

    "github.com/mytheresa/go-hiring-challenge/models"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
)

type fakeIProductRepository struct {
    products []models.Product
    err      error
}

func (f *fakeIProductRepository) GetAllProducts() ([]models.Product, error) {
    if f.err != nil {
        return nil, f.err
    }
    return f.products, nil
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