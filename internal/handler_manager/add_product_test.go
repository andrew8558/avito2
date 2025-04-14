package handler_manager

import (
	customErrors "avito2/internal/errors"
	"avito2/internal/middleware"
	"avito2/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AddProduct(t *testing.T) {
	t.Parallel()

	var (
		request = model.AddProductRequest{
			Type:  model.ProductTypeClothes,
			PvzId: uuid.New().String(),
		}
		employeeRole                  = string(model.RoleEmployee)
		invalidRole                   = "test"
		requestWithInvalidProductType = model.AddProductRequest{
			Type:  "test",
			PvzId: uuid.New().String(),
		}
		requestWithInvalidPvzId = model.AddProductRequest{
			Type:  model.ProductTypeClothes,
			PvzId: "test",
		}
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().AddProduct(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.AddProduct(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("access denied", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, invalidRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.AddProduct(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("invalid productType", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, err := json.Marshal(requestWithInvalidProductType)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.AddProduct(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("invalid pvz_id format", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, err := json.Marshal(requestWithInvalidPvzId)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.AddProduct(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, err := json.Marshal("test")
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.AddProduct(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid http method", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodDelete, "/products", bytes.NewReader(nil))
		rec := httptest.NewRecorder()

		s.hm.AddProduct(rec, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().AddProduct(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to add product"))
		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.AddProduct(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
	t.Run("pvz does not exist", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().AddProduct(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, customErrors.ErrPvzDoesNotExist)
		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.AddProduct(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("no reception in progress", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().AddProduct(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, customErrors.ErrReceptionInProgressDoesNotExist)
		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.AddProduct(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
