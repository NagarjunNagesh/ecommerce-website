package catalog

import "github.com/mytheresa/go-hiring-challenge/app/api"

// Product represents the basic product info for lists
type Product struct {
	Code     string                `json:"code"`
	Price    float64               `json:"price"`
	Category *api.CategoryResponse `json:"category"`
}

// ProductDetail represents the detailed product info including variants
type ProductDetail struct {
	Code     string                `json:"code"`
	Price    float64               `json:"price"`
	Category *api.CategoryResponse `json:"category"`
	Variants []ProductVariant      `json:"variants"`
}

// ProductVariant represents a single product variant detail
type ProductVariant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

// Response represents the catalog list response
type Response struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}

// DetailResponse represents the product detail response
type DetailResponse struct {
	Product ProductDetail `json:"product"`
}
