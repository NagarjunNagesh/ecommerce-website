package categories

import (
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
