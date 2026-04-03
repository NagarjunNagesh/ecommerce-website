package categories

import (
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ICategoryRepository interface {
	GetAllCategories() ([]models.Category, error)
}

type ICategoriesService interface {
	ListCategories() ([]api.CategoryResponse, error)
}

type CategoriesService struct {
	repo ICategoryRepository
}

func NewCategoriesService(repo ICategoryRepository) *CategoriesService {
	return &CategoriesService{repo: repo}
}

func (s *CategoriesService) ListCategories() ([]api.CategoryResponse, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		return nil, err
	}
	return toCategoryResponse(categories), nil
}

func toCategoryResponse(categories []models.Category) []api.CategoryResponse {
	mapped := make([]api.CategoryResponse, len(categories))
	for i, c := range categories {
		mapped[i] = api.CategoryResponse{
			Code: c.Code,
			Name: c.Name,
		}
	}
	return mapped
}
