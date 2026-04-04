package catalog

import (
	"errors"
	"log"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

var ErrProductNotFound = errors.New("product not found")

type IProductRepository interface {
	GetAllProducts() ([]models.Product, error)
	GetProductByCode(code string) (*models.Product, error)
}

type ICatalogService interface {
	ListProducts() ([]Product, error)
	GetProductDetail(code string) (*ProductDetail, error)
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

func (s *CatalogService) GetProductDetail(code string) (*ProductDetail, error) {
	log.Printf("catalog service: getting product detail for code=%s", code)

	product, err := s.repo.GetProductByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("catalog service: product not found for code=%s", code)
			return nil, ErrProductNotFound
		}

		log.Printf("catalog service: failed to fetch product detail for code=%s: %v", code, err)
		return nil, err
	}

	detail := toProductDetail(*product)
	return &detail, nil
}

func toCatalogProducts(products []models.Product) []Product {
	mapped := make([]Product, len(products))
	for i, p := range products {
		mapped[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: toCatalogCategory(p.Category),
		}
	}
	return mapped
}

func toProductDetail(product models.Product) ProductDetail {
	variants := make([]ProductVariant, len(product.Variants))
	for i, variant := range product.Variants {
		price := variant.Price
		if price.IsZero() {
			log.Printf("catalog service: variant price is zero, inheriting product price for variant=%s", variant.SKU)
			price = product.Price
		}

		variants[i] = ProductVariant{
			Name:  variant.Name,
			SKU:   variant.SKU,
			Price: price.InexactFloat64(),
		}
	}

	return ProductDetail{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: toCatalogCategory(product.Category),
		Variants: variants,
	}
}

func toCatalogCategory(category *models.Category) *api.CategoryResponse {
	if category == nil {
		return nil
	}

	return &api.CategoryResponse{
		Code: category.Code,
		Name: category.Name,
	}
}
