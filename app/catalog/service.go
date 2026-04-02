package catalog

import (
	"log"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Product struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

type IProductRepository interface {
	GetAllProducts() ([]models.Product, error)
}

type ICatalogService interface {
	ListProducts() ([]Product, error)
}

type CatalogService struct {
	repo IProductRepository
}

func NewCatalogService(repo IProductRepository) *CatalogService {
	return &CatalogService{repo: repo}
}

func (s *CatalogService) ListProducts() ([]Product, error) {
	log.Printf("catalog service: listing products")

	products, err := s.repo.GetAllProducts()
	if err != nil {
		log.Printf("catalog service: failed to fetch products: %v", err)
		return nil, err
	}

	log.Printf("catalog service: fetched %d products", len(products))

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
