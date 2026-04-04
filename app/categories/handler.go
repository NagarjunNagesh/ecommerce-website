package categories

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
)

type CategoriesHandler struct {
	service ICategoriesService
}

func NewCategoriesHandler(r ICategoryRepository) *CategoriesHandler {
	return &CategoriesHandler{
		service: NewCategoriesService(r),
	}
}

func (h *CategoriesHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("categories handler: %s %s", r.Method, r.URL.Path)

	categories, err := h.service.ListCategories()
	if err != nil {
		log.Printf("categories handler: failed to list categories: %v", err)
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("categories handler: returning %d categories", len(categories))

	api.OKResponse(w, categories)
}

func (h *CategoriesHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	log.Printf("categories handler: %s %s", r.Method, r.URL.Path)

	var req CreateCategoryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("categories handler: failed to decode request: %v", err)
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	category, err := h.service.CreateCategory(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCategoryInput) || errors.Is(err, ErrCategoryAlreadyExists) {
			log.Printf("categories handler: validation/duplicate error: %v", err)
			api.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		log.Printf("categories handler: failed to create category: %v", err)
		api.ErrorResponse(w, http.StatusInternalServerError, "failed to create category")
		return
	}

	api.JSONResponse(w, http.StatusCreated, category)
}
