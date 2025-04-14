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

func Test_DeleteLastProduct(t *testing.T) {
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

		s.mockSvc.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).Return(nil)
		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/delete_last_product", http.HandlerFunc(s.hm.DeleteLastProduct))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/delete_last_product", nil)
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
		r.Handle("/pvz/{pvzId}/delete_last_product", http.HandlerFunc(s.hm.DeleteLastProduct))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/delete_last_product", nil)
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
		r.Handle("/pvz/{pvzId}/delete_last_product", http.HandlerFunc(s.hm.DeleteLastProduct))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+invalidPvzId+"/delete_last_product", nil)
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
		r.Handle("/pvz/{pvzId}/delete_last_product", http.HandlerFunc(s.hm.DeleteLastProduct))
		req := httptest.NewRequest(http.MethodGet, "/pvz/"+pvzId+"/delete_last_product", nil)
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

		s.mockSvc.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).Return(errors.New("failed to delete last product"))
		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/delete_last_product", http.HandlerFunc(s.hm.DeleteLastProduct))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/delete_last_product", nil)
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

		s.mockSvc.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).Return(customErrors.ErrPvzDoesNotExist)
		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/delete_last_product", http.HandlerFunc(s.hm.DeleteLastProduct))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/delete_last_product", nil)
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

		s.mockSvc.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).Return(customErrors.ErrReceptionInProgressDoesNotExist)
		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/delete_last_product", http.HandlerFunc(s.hm.DeleteLastProduct))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/delete_last_product", nil)
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("no product to delete", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockSvc.EXPECT().DeleteLastProduct(gomock.Any(), gomock.Any()).Return(customErrors.ErrNoProductToDelete)
		r := mux.NewRouter()
		r.Handle("/pvz/{pvzId}/delete_last_product", http.HandlerFunc(s.hm.DeleteLastProduct))
		req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzId+"/delete_last_product", nil)
		ctx := context.WithValue(req.Context(), middleware.Role, employeeRole)
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
