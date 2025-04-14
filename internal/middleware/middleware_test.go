package middleware

import (
	"avito2/internal/model"
	"avito2/internal/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Middleware(t *testing.T) {
	t.Parallel()

	var (
		secret = "test"
		role   = string(model.RoleEmployee)
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		os.Setenv("JWT_SECRET", secret)
		jwtGen := utils.JWTGen{}
		token, err := jwtGen.GenerateJWT(role)
		require.NoError(t, err)

		called := false
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			val := r.Context().Value(Role)
			assert.Equal(t, role, val)
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		rec := httptest.NewRecorder()

		mw := AuthMiddleware(nextHandler)
		mw.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, called)
	})
	t.Run("empty token", func(t *testing.T) {
		t.Parallel()

		token := ""

		called := false
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		rec := httptest.NewRecorder()

		mw := AuthMiddleware(nextHandler)
		mw.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.False(t, called)
	})
	t.Run("invalid token", func(t *testing.T) {
		t.Parallel()

		token := "test"

		called := false
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		rec := httptest.NewRecorder()

		mw := AuthMiddleware(nextHandler)
		mw.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.False(t, called)
	})
}
