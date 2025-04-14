package handler_manager

import (
	"avito2/internal/errors"
	"avito2/internal/middleware"
	"avito2/internal/model"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (hm *HandlerManager) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.ErrInvalidHtppMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	role := r.Context().Value(middleware.Role).(string)
	if role != string(model.RoleEmployee) {
		http.Error(w, errors.ErrAccessDenied.Error(), http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	pvzId := vars["pvzId"]

	uuid, err := uuid.Parse(pvzId)
	if err != nil {
		http.Error(w, errors.ErrInvalidPvzIdFormat.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	res, err := hm.svc.CloseLastReception(ctx, uuid)

	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
