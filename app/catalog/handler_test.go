package catalog

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const contentType = "Content-Type"
const applicationJSON = "application/json"
const catalogPath = "/catalog"

type fakeProductRepository struct {
	products []models.Product
	product  *models.Product
	err      error
}

func (f *fakeProductRepository) GetAllProducts() ([]models.Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.products, nil
}

func (f *fakeProductRepository) GetProductByCode(code string) (*models.Product, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.product == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return f.product, nil
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
		request := httptest.NewRequest(http.MethodGet, catalogPath, nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, applicationJSON, recorder.Header().Get(contentType))
		assert.JSONEq(t, `{"products":[{"code":"PROD001","price":99.99,"category":{"code":"clothing","name":"Clothing"}},{"code":"PROD002","price":120,"category":{"code":"shoes","name":"Shoes"}}]}`, recorder.Body.String())
	})

	t.Run("returns internal server error when repository fails", func(t *testing.T) {
		repo := &fakeProductRepository{err: errors.New("database unavailable")}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath, nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "database unavailable")
	})
}

func TestCatalogHandlerHandleGetDetail(t *testing.T) {
	t.Run("returns product details as json", func(t *testing.T) {
		repo := &fakeProductRepository{
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
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath + "/PROD001", nil)
		request.SetPathValue("code", "PROD001")

		handler.HandleGetDetail(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, applicationJSON, recorder.Header().Get(contentType))
		assert.JSONEq(t, `{"product":{"code":"PROD001","price":99.99,"category":{"code":"clothing","name":"Clothing"},"variants":[{"name":"Variant A","sku":"SKU001A","price":89.99},{"name":"Variant B","sku":"SKU001B","price":99.99}]}}`, recorder.Body.String())
	})

	t.Run("returns not found when product does not exist", func(t *testing.T) {
		repo := &fakeProductRepository{}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath + "/PROD999", nil)
		request.SetPathValue("code", "PROD999")

		handler.HandleGetDetail(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "product not found")
	})

	t.Run("returns internal server error when repository fails", func(t *testing.T) {
		repo := &fakeProductRepository{err: errors.New("database unavailable")}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath + "/PROD001", nil)
		request.SetPathValue("code", "PROD001")

		handler.HandleGetDetail(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "database unavailable")
	})
}
