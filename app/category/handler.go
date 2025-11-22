package category

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoryHandler struct {
	repo models.ICategoriesRepository
}

func NewCategoryHandler(r models.ICategoriesRepository) *CategoryHandler {
	return &CategoryHandler{
		repo: r,
	}
}

func (h *CategoryHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.OKResponse(w, categories)
}

func (h *CategoryHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if category.Code == "" || category.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "category code and name are required")
		return
	}

	if err := h.repo.CreateCategory(&category); err != nil {
		if errors.Is(err, models.ErrDuplicateCategory) {
			api.ErrorResponse(w, http.StatusConflict, "category already exists")
			return
		}
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, category)

}
