package categories

import (
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ICategoryRepository interface {
	GetAllCategories() ([]models.Category, error)
}

type ICategoriesService interface {
	ListCategories() ([]Category, error)
}

type CategoriesService struct {
	repo ICategoryRepository
}

func NewCategoriesService(repo ICategoryRepository) *CategoriesService {
	return &CategoriesService{repo: repo}
}

func (s *CategoriesService) ListCategories() ([]Category, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		return nil, err
	}
	return toCategoryResponse(categories), nil
}

func toCategoryResponse(categories []models.Category) []Category {
	mapped := make([]Category, len(categories))
	for i, c := range categories {
		mapped[i] = Category{
			Code: c.Code,
			Name: c.Name,
		}
	}
	return mapped
}
