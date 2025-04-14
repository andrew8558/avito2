package handler_manager

import (
	"avito2/internal/model"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DummyLogin(t *testing.T) {
	t.Parallel()
	var (
		request                = model.DummyLoginRequest{Role: model.RoleEmployee}
		requestWithInvalidRole = model.DummyLoginRequest{Role: "invalid"}
		generateJwtErr         = errors.New("failed to geenrate jwt")
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		s.mockJWTGen.EXPECT().GenerateJWT(gomock.Any()).Return("token", nil)
		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.DummyLogin(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("invalid role", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, err := json.Marshal(requestWithInvalidRole)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.DummyLogin(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("invalid http method", func(t *testing.T) {
		t.Parallel()

		s := setUp(t)
		defer s.tearDown()

		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodGet, "/dummyLogin", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.DummyLogin(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})
	t.Run("invalid json", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		body, err := json.Marshal("bad request")
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.DummyLogin(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
	t.Run("failed to generate jwt", func(t *testing.T) {
		s := setUp(t)
		defer s.tearDown()

		s.mockJWTGen.EXPECT().GenerateJWT(gomock.Any()).Return("", generateJwtErr)
		body, err := json.Marshal(request)
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		s.hm.DummyLogin(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
