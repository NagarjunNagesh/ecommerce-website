package catalog

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
)

var productCodeRE = regexp.MustCompile(`^PROD\d{3,}$`)

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

func (h *CatalogHandler) HandleGetDetail(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	log.Printf("catalog detail handler: received request for code=%s", code)

	if err := validateProductCode(code); err != nil {
		log.Printf("catalog detail handler: invalid product code '%s': %v", code, err)
        api.ErrorResponse(w, http.StatusBadRequest, err.Error())
        return
    }

	log.Printf("catalog detail handler: %s %s code=%s", r.Method, r.URL.Path, code)

	product, err := h.service.GetProductDetail(code)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			log.Printf("catalog detail handler: product not found for code=%s", code)
			api.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		log.Printf("catalog detail handler: failed to get product detail for code=%s: %v", code, err)
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("catalog detail handler: returning product %s with %d variants", product.Code, len(product.Variants))

	api.OKResponse(w, DetailResponse{Product: *product})
}

func validateProductCode(code string) error {
	code = strings.TrimSpace(code)

	switch {
	case code == "":
		return errors.New("code is required")
	case !productCodeRE.MatchString(code):
		return errors.New("invalid product code format")
	default:
		return nil
	}
}
