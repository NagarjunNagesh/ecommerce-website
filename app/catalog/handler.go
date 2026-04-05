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

const (
	paramOffset        = "offset"
	paramLimit         = "limit"
	paramCategory      = "category"
	paramPriceLessThan = "priceLessThan"

	defaultOffset = 0
	defaultLimit  = 10
	minLimit      = 1
	maxLimit      = 100
)

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

	input, err := parseListProductsInput(r)
	if err != nil {
		log.Printf("catalog handler: invalid catalog list params: %v", err)
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	products, total, err := h.service.ListProducts(input)
	if err != nil {
		log.Printf("catalog handler: failed to list products: %v", err)
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("catalog handler: returning %d products out of %d total", len(products), total)

	api.OKResponse(w, Response{Products: products, Total: total})
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

func parseListProductsInput(r *http.Request) (ListProductsInput, error) {
	offset, limit, err := parsePaginationParams(r)
	if err != nil {
		return ListProductsInput{}, err
	}

	query := r.URL.Query()
	priceLessThan, err := parsePositiveDecimal(query, paramPriceLessThan)
	if err != nil {
		return ListProductsInput{}, err
	}

	return ListProductsInput{
		Offset:        offset,
		Limit:         limit,
		Category:      strings.ToLower(strings.TrimSpace(query.Get(paramCategory))),
		PriceLessThan: priceLessThan,
	}, nil
}

func parsePaginationParams(r *http.Request) (int, int, error) {
	query := r.URL.Query()

	offset, err := parseIntOrDefault(query, paramOffset, defaultOffset)
	if err != nil {
		return 0, 0, err
	}
	if offset < 0 {
		return 0, 0, errors.New("offset must be greater than or equal to 0")
	}

	limit, err := parseIntOrDefault(query, paramLimit, defaultLimit)
	if err != nil {
		return 0, 0, err
	}
	if limit < minLimit || limit > maxLimit {
		return 0, 0, errors.New("limit must be between 1 and 100")
	}

	return offset, limit, nil
}
