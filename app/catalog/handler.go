package catalog

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
)

type Response struct {
	Products []Product `json:"products"`
}

type CatalogHandler struct {
	service *CatalogService
}

func NewCatalogHandler(r ProductRepository) *CatalogHandler {
	return &CatalogHandler{
		service: NewCatalogService(r),
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, Response{Products: products})
}
