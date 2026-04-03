package categories

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

type fakeCategoryRepository struct {
	categories []models.Category
	err        error
}

func (f *fakeCategoryRepository) GetAllCategories() ([]models.Category, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.categories, nil
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
		request := httptest.NewRequest(http.MethodGet, "/categories", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.JSONEq(t, `[{"code":"clothing","name":"Clothing"},{"code":"shoes","name":"Shoes"}]`, recorder.Body.String())
	})

	t.Run("returns internal server error when repository fails", func(t *testing.T) {
		repo := &fakeCategoryRepository{err: errors.New("database unavailable")}
		handler := NewCategoriesHandler(repo)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/categories", nil)

		handler.HandleGet(recorder, request)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"error":"database unavailable"}`, recorder.Body.String())
	})
}
