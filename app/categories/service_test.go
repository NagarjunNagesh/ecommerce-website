package categories

import (
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

type fakeCategoryRepositoryForService struct {
	categories []models.Category
	err        error
}

func (f *fakeCategoryRepositoryForService) GetAllCategories() ([]models.Category, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.categories, nil
}

func TestCategoriesServiceListCategories(t *testing.T) {
	t.Run("maps repository categories to api response", func(t *testing.T) {
		repo := &fakeCategoryRepositoryForService{
			categories: []models.Category{
				{Code: "clothing", Name: "Clothing"},
				{Code: "shoes", Name: "Shoes"},
			},
		}
		service := NewCategoriesService(repo)

		categories, err := service.ListCategories()

		assert.NoError(t, err)
		assert.Len(t, categories, 2)
		assert.Equal(t, "clothing", categories[0].Code)
		assert.Equal(t, "Clothing", categories[0].Name)
		assert.Equal(t, "shoes", categories[1].Code)
		assert.Equal(t, "Shoes", categories[1].Name)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		repo := &fakeCategoryRepositoryForService{err: errors.New("database unavailable")}
		service := NewCategoriesService(repo)

		categories, err := service.ListCategories()

		assert.Nil(t, categories)
		assert.EqualError(t, err, "database unavailable")
	})
}
