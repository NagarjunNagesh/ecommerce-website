package categories

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

var (
	ErrCategoryAlreadyExists = errors.New("category code already exists")
	ErrInvalidCategoryInput  = errors.New("code and name are required")
	ErrInvalidCategoryCode   = errors.New("code must contain only lowercase letters, numbers, and hyphens")

	categoryCodeRE = regexp.MustCompile(`^[a-z0-9-]+$`)
)

type ICategoryRepository interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category *models.Category) error
}

type ICategoriesService interface {
	ListCategories() ([]api.CategoryResponse, error)
	CreateCategory(req CreateCategoryRequest) (*api.CategoryResponse, error)
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

func (s *CategoriesService) CreateCategory(req CreateCategoryRequest) (*api.CategoryResponse, error) {
	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)

	if code == "" || name == "" {
		return nil, ErrInvalidCategoryInput
	}
	if !categoryCodeRE.MatchString(code) {
		return nil, ErrInvalidCategoryCode
	}

	category := &models.Category{
		Code: code,
		Name: name,
	}

	if err := s.repo.CreateCategory(category); err != nil {
		if errors.Is(err, models.ErrCategoryAlreadyExists) {
			return nil, ErrCategoryAlreadyExists
		}
		return nil, err
	}

	return &api.CategoryResponse{
		Code: category.Code,
		Name: category.Name,
	}, nil
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
