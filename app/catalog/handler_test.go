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
const databaseUnavailable = "database unavailable"

type fakeProductRepository struct {
	products        []models.Product
	total           int64
	product         *models.Product
	err             error
	receivedOptions models.ProductListOptions
}

func (f *fakeProductRepository) ListProducts(options models.ProductListOptions) (models.ProductPage, error) {
	f.receivedOptions = options
	if f.err != nil {
		return models.ProductPage{}, f.err
	}
	return models.ProductPage{Products: f.products, Total: f.total}, nil
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
			total: 8,
		}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath, nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, applicationJSON, recorder.Header().Get(contentType))
		assert.Equal(t, 0, repo.receivedOptions.Offset)
		assert.Equal(t, 10, repo.receivedOptions.Limit)
		assert.JSONEq(t, `{"products":[{"code":"PROD001","price":99.99,"category":{"code":"clothing","name":"Clothing"}},{"code":"PROD002","price":120,"category":{"code":"shoes","name":"Shoes"}}],"total":8}`, recorder.Body.String())
	})

	t.Run("uses explicit pagination params", func(t *testing.T) {
		repo := &fakeProductRepository{
			products: []models.Product{
				{Code: "PROD003", Price: decimal.NewFromFloat(59.99), Category: &models.Category{Code: "accessories", Name: "Accessories"}},
			},
			total: 8,
		}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"?offset=2&limit=1", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, 2, repo.receivedOptions.Offset)
		assert.Equal(t, 1, repo.receivedOptions.Limit)
		assert.JSONEq(t, `{"products":[{"code":"PROD003","price":59.99,"category":{"code":"accessories","name":"Accessories"}}],"total":8}`, recorder.Body.String())
	})

	t.Run("forwards category and price filters", func(t *testing.T) {
		repo := &fakeProductRepository{total: 1}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"?category=Clothing&priceLessThan=100", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "clothing", repo.receivedOptions.CategoryCode)
		assert.NotNil(t, repo.receivedOptions.PriceLessThan)
		assert.True(t, repo.receivedOptions.PriceLessThan.Equal(decimal.RequireFromString("100")))
	})

	t.Run("returns bad request for invalid priceLessThan", func(t *testing.T) {
		repo := &fakeProductRepository{}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"?priceLessThan=cheap", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "priceLessThan must be a valid number")
	})

	t.Run("returns bad request for non-positive priceLessThan", func(t *testing.T) {
		repo := &fakeProductRepository{}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"?priceLessThan=0", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "priceLessThan must be greater than 0")
	})

	t.Run("returns internal server error when repository fails", func(t *testing.T) {
		repo := &fakeProductRepository{err: errors.New(databaseUnavailable)}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath, nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), databaseUnavailable)
	})

	t.Run("returns bad request for invalid offset", func(t *testing.T) {
		repo := &fakeProductRepository{}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"?offset=abc", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "offset must be a valid integer")
	})

	t.Run("returns bad request for negative offset", func(t *testing.T) {
		repo := &fakeProductRepository{}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"?offset=-1", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "offset must be greater than or equal to 0")
	})

	t.Run("returns bad request for invalid limit", func(t *testing.T) {
		repo := &fakeProductRepository{}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"?limit=0", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "limit must be between 1 and 100")
	})

	t.Run("returns bad request for limit above maximum", func(t *testing.T) {
		repo := &fakeProductRepository{}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"?limit=101", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "limit must be between 1 and 100")
	})

	t.Run("returns bad request for non numeric limit", func(t *testing.T) {
		repo := &fakeProductRepository{}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"?limit=ten", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "limit must be a valid integer")
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
		request := httptest.NewRequest(http.MethodGet, catalogPath+"/PROD001", nil)
		request.SetPathValue("code", "PROD001")

		handler.HandleGetDetail(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, applicationJSON, recorder.Header().Get(contentType))
		assert.JSONEq(t, `{"product":{"code":"PROD001","price":99.99,"category":{"code":"clothing","name":"Clothing"},"variants":[{"name":"Variant A","sku":"SKU001A","price":89.99},{"name":"Variant B","sku":"SKU001B","price":99.99}]}}`, recorder.Body.String())
	})

	t.Run("returns not found when product does not exist", func(t *testing.T) {
		repo := &fakeProductRepository{product: nil}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"/PROD999", nil)
		request.SetPathValue("code", "PROD999")

		handler.HandleGetDetail(recorder, request)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("returns bad request for invalid product code format", func(t *testing.T) {
		handler := NewCatalogHandler(&fakeProductRepository{})

		codes := []string{"invalid", "P123", "", "PROD12"}
		for _, code := range codes {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, catalogPath+"/"+code, nil)
			request.SetPathValue("code", code)

			handler.HandleGetDetail(recorder, request)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("returns internal server error when repository fails", func(t *testing.T) {
		repo := &fakeProductRepository{err: errors.New(databaseUnavailable)}
		handler := NewCatalogHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, catalogPath+"/PROD001", nil)
		request.SetPathValue("code", "PROD001")

		handler.HandleGetDetail(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Contains(t, recorder.Body.String(), databaseUnavailable)
	})
}
