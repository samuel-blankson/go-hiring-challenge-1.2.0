package catalog

import (
	"net/http"
	"strconv"
	"strings"

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

type ProductDetailDTO struct {
	Code     string       `json:"code"`
	Price    float64      `json:"price"`
	Category string       `json:"category"`
	Variants []VariantDTO `json:"variants"`
}

type VariantDTO struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type CatalogHandler struct {
	repo models.IProductsRepository
}

func NewCatalogHandler(r models.IProductsRepository) *CatalogHandler {
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

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code, err := extractCodeFromPath(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Product code missing: "+err.Error())
		return
	}

	product, err := h.repo.GetProductByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if product == nil {
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// Map product to DTO
	detail := mapProductDetail(*product)
	api.OKResponse(w, detail)
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

func mapProductDetail(product models.Product) ProductDetailDTO {
	var variants []VariantDTO

	for _, v := range product.Variants {
		price := v.Price
		if price.IsZero() {
			price = product.Price // inherit price from parent product if variant is zero
		}
		variants = append(variants, VariantDTO{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price.InexactFloat64(),
		})
	}

	return ProductDetailDTO{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: product.Category.Name,
		Variants: variants,
	}
}

func extractCodeFromPath(r *http.Request) (string, error) {
	path := r.URL.Path                // "/catalog/PROD001"
	parts := strings.Split(path, "/") // ["", "catalog", "PROD001"]

	if len(parts) < 3 || parts[2] == "" {
		return "", http.ErrNoLocation
	}

	code := parts[2]
	return code, nil
}
