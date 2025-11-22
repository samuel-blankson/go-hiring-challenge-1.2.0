package catalog_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Mock repository implementing ProductsRepository interface
type MockProductsRepo struct{}

func (m *MockProductsRepo) GetAllProducts(offset, limit int, filter models.ProductFilter) ([]models.Product, int64, error) {
	category := models.Category{ID: 1, Code: "CLOTHING", Name: "Clothing"}
	price := decimal.NewFromFloat(10.99)

	// Filtering logic for testing
	if filter.CategoryCode != "" && filter.CategoryCode != category.Code {
		return []models.Product{}, 0, nil
	}

	if filter.PriceLessThan > 0 && price.GreaterThan(decimal.NewFromFloat(filter.PriceLessThan)) {
		return []models.Product{}, 0, nil
	}

	products := []models.Product{
		{
			Code:     "PROD001",
			Price:    price,
			Category: category,
		},
	}

	return products, int64(len(products)), nil
}

func (m *MockProductsRepo) GetProductByCode(code string) (*models.Product, error) {
	if code != "PROD001" {
		return nil, nil // simulate non-existence product
	}

	category := models.Category{ID: 1, Code: "CLOTHING", Name: "Clothing"}
	price := decimal.NewFromFloat(10.99)
	variantPrice := decimal.NewFromFloat(11.99)

	product := &models.Product{
		Code:     "PROD001",
		Price:    price,
		Category: category,
		Variants: []models.Variant{
			{Name: "Variant A", SKU: "SKU001A", Price: variantPrice},
			{Name: "Variant B", SKU: "SKU001B", Price: decimal.Zero}, // inherits product price
		},
	}

	return product, nil
}

func TestHandleGet(t *testing.T) {
	handler := catalog.NewCatalogHandler(&MockProductsRepo{})

	t.Run("valid request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?offset=0&limit=10", nil)
		w := httptest.NewRecorder()
		handler.HandleGet(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	})

	t.Run("invalid offset", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?offset=abc&limit=10", nil)
		w := httptest.NewRecorder()
		handler.HandleGet(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?offset=0&limit=abc", nil)
		w := httptest.NewRecorder()
		handler.HandleGet(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("filter by non-matching category", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?offset=0&limit=10&category=SHOES", nil)
		w := httptest.NewRecorder()
		handler.HandleGet(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("filter by price less than", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog?offset=0&limit=10&price_less_than=5", nil)
		w := httptest.NewRecorder()
		handler.HandleGet(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

}

func TestHandleGetByCode(t *testing.T) {
	handler := catalog.NewCatalogHandler(&MockProductsRepo{})

	t.Run("existing product", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
		w := httptest.NewRecorder()
		handler.HandleGetByCode(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("non-existing product", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/catalog/PROD999", nil)
		w := httptest.NewRecorder()
		handler.HandleGetByCode(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
