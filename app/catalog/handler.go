package catalog

import (
	"log"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
)

type Response struct {
	Products []Product `json:"products"`
}

type CatalogHandler struct {
	service ICatalogService
}

func NewCatalogHandler(r IProductRepository) *CatalogHandler {
	return &CatalogHandler{
		service: NewCatalogService(r),
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("catalog handler: %s %s", r.Method, r.URL.Path)

	products, err := h.service.ListProducts()
	if err != nil {
		log.Printf("catalog handler: failed to list products: %v", err)
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("catalog handler: returning %d products", len(products))

	api.OKResponse(w, Response{Products: products})
}
