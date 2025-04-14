package handler_manager

import (
	customErrors "avito2/internal/errors"
	"avito2/internal/middleware"
	"avito2/internal/model"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_CloseLastReception(t *testing.T) {
	t.Parallel()

	var (
		pvzId        = uuid.New().String()
		invalidPvzId = "test"
		employeeRole = string(model.RoleEmployee)
		invalidRole  = "test"
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().CloseLastReception(gomock.Any(), gomock.Any()).Return(nil, nil)
		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/close_last_reception", http.HandlerFunc(s.hm.CloseLastReception))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/close_last_reception", nil)
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("access denied", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/close_last_reception", http.HandlerFunc(s.hm.CloseLastReception))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/close_last_reception", nil)
		ctx := context.WithValue(req.Context(), middleware.Role, invalidRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
	t.Run("invalid pvz_id format", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/close_last_reception", http.HandlerFunc(s.hm.CloseLastReception))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+invalidPvzId+"/close_last_reception", nil)
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("invalid http method", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/close_last_reception", http.HandlerFunc(s.hm.CloseLastReception))
		req := httptest.NewRequest(http.MethodGet, "/pvz/"+pvzId+"/close_last_reception", nil)
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().CloseLastReception(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to close last reception"))
		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/close_last_reception", http.HandlerFunc(s.hm.CloseLastReception))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/close_last_reception", nil)
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
	t.Run("pvz does not exist", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().CloseLastReception(gomock.Any(), gomock.Any()).Return(nil, customErrors.ErrPvzDoesNotExist)
		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/close_last_reception", http.HandlerFunc(s.hm.CloseLastReception))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/close_last_reception", nil)
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("reception in progress does not exist", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().CloseLastReception(gomock.Any(), gomock.Any()).Return(nil, customErrors.ErrReceptionInProgressDoesNotExist)
		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/close_last_reception", http.HandlerFunc(s.hm.CloseLastReception))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/close_last_reception", nil)
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
