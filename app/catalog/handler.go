package catalog

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products []ProductDTO `json:"products"`
}

type ProductDTO struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
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
	products, err := h.repo.GetAllProducts()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map response
	response := mapProducts(products)

	api.OKResponse(w, response)
}

func mapProducts(products []models.Product) Response {
	res := make([]ProductDTO, len(products))
	for i, p := range products {
		res[i] = ProductDTO{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
		}
	}

	return Response{Products: res}
}
