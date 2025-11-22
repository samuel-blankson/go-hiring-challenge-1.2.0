package catalog

import (
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products []ProductDTO `json:"products"`
	Total    int64        `json:"total"`
}

type ProductDTO struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type CatalogHandler struct {
	repo *models.ProductsRepository
}

func NewCatalogHandler(r *models.ProductsRepository) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Parse pagination
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid offset value: "+err.Error())
		return
	}

	if offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid limit value: "+err.Error())
		return
	}

	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 10
	}

	// parse filters
	category := r.URL.Query().Get("category")
	priceStr := r.URL.Query().Get("price_less_than")
	var PriceLessThan float64

	if priceStr != "" {
		PriceLessThan, err = strconv.ParseFloat(priceStr, 64)
		if err != nil {
			api.ErrorResponse(w, http.StatusBadRequest, "invalid price_less_than value: "+err.Error())
			return
		}
	}

	filter := models.ProductFilter{
		CategoryCode:  category,
		PriceLessThan: PriceLessThan,
	}

	products, total, err := h.repo.GetAllProducts(offset, limit, filter)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map response
	response := Response{
		Products: mapProducts(products),
		Total:    total,
	}

	api.OKResponse(w, response)
}

func mapProducts(products []models.Product) []ProductDTO {
	res := make([]ProductDTO, len(products))
	for i, p := range products {
		res[i] = ProductDTO{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	return res
}
