package categories

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

const contentType = "Content-Type"
const applicationJSON = "application/json"
const categoriesPath = "/categories"
const handlerTestCategoryCode = "test"
const handlerTestCategoryName = "Test Category"
const handlerTestDatabaseUnavailable = "database unavailable"

type fakeCategoryRepository struct {
	categories      []models.Category
	err             error
	createdCategory *models.Category
}

func (f *fakeCategoryRepository) GetAllCategories() ([]models.Category, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.categories, nil
}

func (f *fakeCategoryRepository) CreateCategory(category *models.Category) error {
	if f.err != nil {
		return f.err
	}
	f.createdCategory = category
	return nil
}

func TestCategoriesHandlerHandleGet(t *testing.T) {
	t.Run("returns categories as json", func(t *testing.T) {
		repo := &fakeCategoryRepository{
			categories: []models.Category{
				{Code: "clothing", Name: "Clothing"},
				{Code: "shoes", Name: "Shoes"},
			},
		}
		handler := NewCategoriesHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, categoriesPath, nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, applicationJSON, recorder.Header().Get(contentType))
		assert.JSONEq(t, `[{"code":"clothing","name":"Clothing"},{"code":"shoes","name":"Shoes"}]`, recorder.Body.String())
	})

	t.Run("returns internal server error when repository fails", func(t *testing.T) {
		repo := &fakeCategoryRepository{err: errors.New(handlerTestDatabaseUnavailable)}
		handler := NewCategoriesHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, categoriesPath, nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, applicationJSON, recorder.Header().Get(contentType))
		assert.JSONEq(t, `{"error":"database unavailable"}`, recorder.Body.String())
	})
}

func TestCategoriesHandlerHandlePost(t *testing.T) {
	t.Run("creates a category", func(t *testing.T) {
		repo := &fakeCategoryRepository{}
		handler := NewCategoriesHandler(repo)

		reqBody, _ := json.Marshal(CreateCategoryRequest{Code: handlerTestCategoryCode, Name: handlerTestCategoryName})
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, categoriesPath, bytes.NewBuffer(reqBody))

		handler.HandlePost(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Equal(t, applicationJSON, recorder.Header().Get(contentType))
		assert.JSONEq(t, `{"code":"test","name":"Test Category"}`, recorder.Body.String())
	})

	t.Run("returns bad request for invalid json", func(t *testing.T) {
		repo := &fakeCategoryRepository{}
		handler := NewCategoriesHandler(repo)

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, categoriesPath, bytes.NewBufferString(`{"invalid":json`))

		handler.HandlePost(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.JSONEq(t, `{"error":"invalid request body"}`, recorder.Body.String())
	})

	t.Run("returns bad request for missing fields", func(t *testing.T) {
		repo := &fakeCategoryRepository{}
		handler := NewCategoriesHandler(repo)

		reqBody, _ := json.Marshal(CreateCategoryRequest{Code: "", Name: handlerTestCategoryName})
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, categoriesPath, bytes.NewBuffer(reqBody))

		handler.HandlePost(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.JSONEq(t, `{"error":"code and name are required"}`, recorder.Body.String())
	})

	t.Run("normalizes category code to lowercase", func(t *testing.T) {
		repo := &fakeCategoryRepository{}
		handler := NewCategoriesHandler(repo)

		reqBody, _ := json.Marshal(CreateCategoryRequest{Code: "  Shoes123  ", Name: handlerTestCategoryName})
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, categoriesPath, bytes.NewBuffer(reqBody))

		handler.HandlePost(recorder, request)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.NotNil(t, repo.createdCategory)
		assert.Equal(t, "shoes123", repo.createdCategory.Code)
		assert.JSONEq(t, `{"code":"shoes123","name":"Test Category"}`, recorder.Body.String())
	})

	t.Run("returns bad request for non-alphanumeric category code", func(t *testing.T) {
		repo := &fakeCategoryRepository{}
		handler := NewCategoriesHandler(repo)

		reqBody, _ := json.Marshal(CreateCategoryRequest{Code: "shoe-items!", Name: handlerTestCategoryName})
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, categoriesPath, bytes.NewBuffer(reqBody))

		handler.HandlePost(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.JSONEq(t, `{"error":"code must contain only lowercase letters and numbers"}`, recorder.Body.String())
	})

	t.Run("returns bad request for unknown fields", func(t *testing.T) {
		repo := &fakeCategoryRepository{}
		handler := NewCategoriesHandler(repo)

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, categoriesPath, bytes.NewBufferString(`{"code":"test", "name":"Test", "extra":"field"}`))

		handler.HandlePost(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.JSONEq(t, `{"error":"invalid request body"}`, recorder.Body.String())
	})

	t.Run("returns bad request (400) for duplicate key per user preference", func(t *testing.T) {
		repo := &fakeCategoryRepository{err: models.ErrCategoryAlreadyExists}
		handler := NewCategoriesHandler(repo)

		reqBody, _ := json.Marshal(CreateCategoryRequest{Code: handlerTestCategoryCode, Name: handlerTestCategoryName})
		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, categoriesPath, bytes.NewBuffer(reqBody))

		handler.HandlePost(recorder, request)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.JSONEq(t, `{"error":"category code already exists"}`, recorder.Body.String())
	})
}
