package handler_manager

import (
	"avito2/internal/middleware"
	"avito2/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PVZ(t *testing.T) {
	t.Parallel()

	var (
		request = model.CreatePvzRequest{
			City: model.CityMoscow,
		}
		moderatorRole          = string(model.RoleModerator)
		invalidRole            = "test"
		requestWithInvalidCity = model.CreatePvzRequest{
			City: "test",
		}
	)

	t.Run("success create pvz", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().CreatePvz(gomock.Any(), gomock.Any()).Return(nil, nil)
		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("access denied", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, invalidRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("invalid city", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, err := json.Marshal(requestWithInvalidCity)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("failed to create pvz info with internal error", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().CreatePvz(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to create pvz"))
		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("succes get pvz info", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().GetPvzInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
		req := httptest.NewRequest(http.MethodGet, "/pvz", bytes.NewReader(nil))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid startDate fromat", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		params := url.Values{}
		params.Add("startDate", "2025.10.02 15:30:30")
		req := httptest.NewRequest(http.MethodGet, "/pvz?"+params.Encode(), bytes.NewReader(nil))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid endDate fromat", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		params := url.Values{}
		params.Add("endDate", "2025.10.02 15:30:30")
		req := httptest.NewRequest(http.MethodGet, "/pvz?"+params.Encode(), bytes.NewReader(nil))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid page value", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodGet, "/pvz?page=0", bytes.NewReader(nil))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid limit value", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodGet, "/pvz?limit=0", bytes.NewReader(nil))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("endDate less than startDate", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		params := url.Values{}
		params.Add("startDate", "2006-01-03 15:04:05")
		params.Add("endDate", "2006-01-02 15:04:05")
		req := httptest.NewRequest(http.MethodGet, "/pvz?"+params.Encode(), bytes.NewReader(nil))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("invalid page type", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodGet, "/pvz?page=test", bytes.NewReader(nil))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid limit type", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodGet, "/pvz?limit=test", bytes.NewReader(nil))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("failed to get pvz info with internal error", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().GetPvzInfo(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to get pvz info"))
		req := httptest.NewRequest(http.MethodGet, "/pvz", bytes.NewReader(nil))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("invalid http method", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		req := httptest.NewRequest(http.MethodDelete, "/pvz", bytes.NewReader(nil))
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, err := json.Marshal("test")
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
		ctx := context.WithValue(req.Context(), middleware.Role, moderatorRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		s.hm.Pvz(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
