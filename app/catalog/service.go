package catalog

import "github.com/mytheresa/go-hiring-challenge/models"

type Product struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

type ProductRepository interface {
	GetAllProducts() ([]models.Product, error)
}

type CatalogService struct {
	repo ProductRepository
}

func NewCatalogService(repo ProductRepository) *CatalogService {
	return &CatalogService{repo: repo}
}

func (s *CatalogService) ListProducts() ([]Product, error) {
	products, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	return toCatalogProducts(products), nil
}

func toCatalogProducts(products []models.Product) []Product {
	mapped := make([]Product, len(products))
	for i, p := range products {
		mapped[i] = Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
		}
	}
	return mapped
}
