package categories

import (
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

const serviceTestCategoryCode = "test"
const serviceTestCategoryName = "Test Category"
const serviceTestDatabaseUnavailable = "database unavailable"

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

func (f *fakeCategoryRepositoryForService) CreateCategory(category *models.Category) error {
	if f.err != nil {
		return f.err
	}
	return nil
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
		repo := &fakeCategoryRepositoryForService{err: errors.New(serviceTestDatabaseUnavailable)}
		service := NewCategoriesService(repo)

		categories, err := service.ListCategories()

		assert.Nil(t, categories)
		assert.EqualError(t, err, serviceTestDatabaseUnavailable)
	})
}

func TestCategoriesServiceCreateCategory(t *testing.T) {
	t.Run("successfully creates a category", func(t *testing.T) {
		repo := &fakeCategoryRepositoryForService{}
		service := NewCategoriesService(repo)

		req := CreateCategoryRequest{Code: serviceTestCategoryCode, Name: serviceTestCategoryName}
		res, err := service.CreateCategory(req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, serviceTestCategoryCode, res.Code)
		assert.Equal(t, serviceTestCategoryName, res.Name)
	})

	t.Run("returns error when code is empty", func(t *testing.T) {
		repo := &fakeCategoryRepositoryForService{}
		service := NewCategoriesService(repo)

		req := CreateCategoryRequest{Code: "  ", Name: serviceTestCategoryName}
		res, err := service.CreateCategory(req)

		assert.Nil(t, res)
		assert.ErrorIs(t, err, ErrInvalidCategoryInput)
	})

	t.Run("returns error when name is empty", func(t *testing.T) {
		repo := &fakeCategoryRepositoryForService{}
		service := NewCategoriesService(repo)

		req := CreateCategoryRequest{Code: serviceTestCategoryCode, Name: ""}
		res, err := service.CreateCategory(req)

		assert.Nil(t, res)
		assert.ErrorIs(t, err, ErrInvalidCategoryInput)
	})

	t.Run("returns error when category code already exists", func(t *testing.T) {
		repo := &fakeCategoryRepositoryForService{err: models.ErrCategoryAlreadyExists}
		service := NewCategoriesService(repo)

		req := CreateCategoryRequest{Code: serviceTestCategoryCode, Name: serviceTestCategoryName}
		res, err := service.CreateCategory(req)

		assert.Nil(t, res)
		assert.ErrorIs(t, err, ErrCategoryAlreadyExists)
	})
}
