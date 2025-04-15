package tests

import (
	"avito2/internal/handler_manager"
	"avito2/internal/middleware"
	"avito2/internal/model"
	"avito2/internal/repository"
	"avito2/internal/service"
	"avito2/internal/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReceptionPipline(t *testing.T) {
	moderatorDummyLoginRequest := model.DummyLoginRequest{
		Role: "moderator",
	}
	employeeDummyLoginRequest := model.DummyLoginRequest{
		Role: "employee",
	}

	t.Run("reception pipline", func(t *testing.T) {
		database.SetUp(t, "pvz", "products", "receptions")
		repo := repository.NewRepository(database.DB)
		svc := service.NewService(repo)
		jwrGen := &utils.JWTGen{}
		hm := handler_manager.NewHandlerManager(svc, jwrGen)

		body, err := json.Marshal(moderatorDummyLoginRequest)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		hm.DummyLogin(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		res := map[string]string{}
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)

		moderatorToken, ok := res["token"]
		assert.True(t, ok)

		body, err = json.Marshal(employeeDummyLoginRequest)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
		rec = httptest.NewRecorder()

		hm.DummyLogin(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		res = map[string]string{}
		err = json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)

		employeeToken, ok := res["token"]
		assert.True(t, ok)

		body, err = json.Marshal(model.CreatePvzRequest{City: model.CityMoscow})
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+moderatorToken)
		rec = httptest.NewRecorder()

		handler := middleware.AuthMiddleware(http.HandlerFunc(hm.Pvz))
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var pvz model.Pvz
		err = json.Unmarshal(rec.Body.Bytes(), &pvz)
		require.NoError(t, err)

		body, err = json.Marshal(model.CreateReceptionRequest{PvzId: pvz.Id.String()})
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+employeeToken)
		rec = httptest.NewRecorder()

		handler = middleware.AuthMiddleware(http.HandlerFunc(hm.CreateReception))
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		for range 50 {
			body, err = json.Marshal(model.AddProductRequest{Type: model.ProductTypeClothes, PvzId: pvz.Id.String()})
			require.NoError(t, err)

			req = httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
			req.Header.Set("Authorization", "Bearer "+employeeToken)
			rec = httptest.NewRecorder()

			handler = middleware.AuthMiddleware(http.HandlerFunc(hm.AddProduct))
			handler.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusCreated, rec.Code)
		}

		router := mux.NewRouter()
		router.Handle("/pvz/{pvzId}/close_last_reception", middleware.AuthMiddleware(http.HandlerFunc(hm.CloseLastReception)))

		req = httptest.NewRequest(http.MethodPost, "/pvz/"+pvz.Id.String()+"/close_last_reception", nil)
		req.Header.Set("Authorization", "Bearer "+employeeToken)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
