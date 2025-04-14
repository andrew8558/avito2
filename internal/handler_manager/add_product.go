package handler_manager

import (
	"avito2/internal/errors"
	"avito2/internal/middleware"
	"avito2/internal/model"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (hm *HandlerManager) AddProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.ErrInvalidHtppMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	role := r.Context().Value(middleware.Role).(string)
	if role != string(model.RoleEmployee) {
		http.Error(w, errors.ErrAccessDenied.Error(), http.StatusForbidden)
		return
	}

	var req model.AddProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, errors.ErrInvalidJson.Error(), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(req.PvzId)
	if err != nil {
		http.Error(w, errors.ErrInvalidPvzIdFormat.Error(), http.StatusBadRequest)
		return
	}

	if !req.Type.IsValid() {
		http.Error(w, "invalid product type", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	res, err := hm.svc.AddProduct(ctx, uuid, req.Type)

	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
		return
	case errors.ErrPvzDoesNotExist:
		http.Error(w, errors.ErrPvzDoesNotExist.Error(), http.StatusBadRequest)
		return
	case errors.ErrReceptionInProgressDoesNotExist:
		http.Error(w, errors.ErrReceptionInProgressDoesNotExist.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}
}
