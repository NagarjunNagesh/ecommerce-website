package catalog

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type fakeProductRepository struct {
	products []models.Product
	err      error
}

func (f *fakeProductRepository) GetAllProducts() ([]models.Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.products, nil
}

func TestCatalogHandlerHandleGet(t *testing.T) {
	t.Run("returns products as json", func(t *testing.T) {
		repo := &fakeProductRepository{
			products: []models.Product{
				{Code: "PROD001", Price: decimal.NewFromFloat(99.99), Category: &models.Category{Code: "clothing", Name: "Clothing"}},
				{Code: "PROD002", Price: decimal.NewFromFloat(120.00), Category: &models.Category{Code: "shoes", Name: "Shoes"}},
			},
		}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/catalog", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"products":[{"code":"PROD001","price":99.99,"category":{"code":"clothing","name":"Clothing"}},{"code":"PROD002","price":120,"category":{"code":"shoes","name":"Shoes"}}]}`, recorder.Body.String())
	})

	t.Run("returns internal server error when repository fails", func(t *testing.T) {
		repo := &fakeProductRepository{err: errors.New("database unavailable")}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/catalog", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "database unavailable")
	})
}
