package category

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

// Mock repository implementing the CategoryRepository interface
type MockCategoryRepo struct {
	GetFn    func() ([]models.Category, error)
	CreateFn func(cat *models.Category) error
}

func (m *MockCategoryRepo) GetAllCategories() ([]models.Category, error) {
	return m.GetFn()
}

func (m *MockCategoryRepo) CreateCategory(cat *models.Category) error {
	return m.CreateFn(cat)
}

// --------------------------
// Tests for HandleGetAll
// --------------------------
func TestHandleGetAll(t *testing.T) {
	t.Run("returns all categories", func(t *testing.T) {
		mock := &MockCategoryRepo{
			GetFn: func() ([]models.Category, error) {
				return []models.Category{
					{Code: "CLOTHING", Name: "Clothing"},
					{Code: "SHOES", Name: "Shoes"},
				}, nil
			},
		}

		h := NewCategoryHandler(mock)
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		w := httptest.NewRecorder()

		h.HandleGet(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res []models.Category
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Len(t, res, 2)
	})

	t.Run("empty list returns 200", func(t *testing.T) {
		mock := &MockCategoryRepo{
			GetFn: func() ([]models.Category, error) {
				return []models.Category{}, nil
			},
		}

		h := NewCategoryHandler(mock)
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		w := httptest.NewRecorder()

		h.HandleGet(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res []models.Category
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Len(t, res, 0)
	})

	t.Run("repo error returns 500", func(t *testing.T) {
		mock := &MockCategoryRepo{
			GetFn: func() ([]models.Category, error) {
				return nil, errors.New("db error")
			},
		}

		h := NewCategoryHandler(mock)
		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		w := httptest.NewRecorder()

		h.HandleGet(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"db error"}`, w.Body.String())
	})

}

// --------------------------
// Tests for HandleCreate
// --------------------------
func TestHandleCreate(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mock := &MockCategoryRepo{
			CreateFn: func(cat *models.Category) error { return nil },
		}
		h := NewCategoryHandler(mock)

		payload := map[string]string{
			"code": "BAGS",
			"name": "Bags & Backpacks",
		}

		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		h.HandleCreate(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var res models.Category
		json.Unmarshal(w.Body.Bytes(), &res)
		assert.Equal(t, "BAGS", res.Code)
		assert.Equal(t, "Bags & Backpacks", res.Name)
	})

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		mock := &MockCategoryRepo{}
		h := NewCategoryHandler(mock)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer([]byte("{ bad json}")))
		w := httptest.NewRecorder()

		h.HandleCreate(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("duplicate category returns 409", func(t *testing.T) {
		mock := &MockCategoryRepo{
			CreateFn: func(cat *models.Category) error { return models.ErrDuplicateCategory },
		}
		h := NewCategoryHandler(mock)

		payload := map[string]string{
			"code": "BAGS",
			"name": "Bags & Backpacks",
		}

		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		h.HandleCreate(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, `{"error":"category already exists"}`, w.Body.String())
	})
}
