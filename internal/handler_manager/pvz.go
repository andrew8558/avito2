package handler_manager

import (
	"avito2/internal/errors"
	"avito2/internal/middleware"
	"avito2/internal/model"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (hm *HandlerManager) Pvz(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		role := r.Context().Value(middleware.Role).(string)
		if role != string(model.RoleModerator) {
			http.Error(w, errors.ErrAccessDenied.Error(), http.StatusForbidden)
			return
		}

		var req model.CreatePvzRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, errors.ErrInvalidJson.Error(), http.StatusBadRequest)
			return
		}

		if !req.City.IsValid() {
			http.Error(w, "invalid city", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		res, err := hm.svc.CreatePvz(ctx, req.City)

		if err != nil {
			http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
		return
	case http.MethodGet:
		queryParams := r.URL.Query()

		startDateStr := queryParams.Get("startDate")
		startDate := time.Time{}
		if startDateStr != "" {
			var err error
			startDate, err = time.Parse(time.DateTime, startDateStr)
			if err != nil {
				http.Error(w, "invalid start date format", http.StatusBadRequest)
				return
			}
		}

		endDateStr := queryParams.Get("endDate")
		endDate := time.Now()
		if endDateStr != "" {
			var err error
			endDate, err = time.Parse(time.DateTime, endDateStr)
			if err != nil {
				http.Error(w, "invalid end date format", http.StatusBadRequest)
				return
			}
		}

		if endDate.Before(startDate) {
			http.Error(w, "end date must be later than start date", http.StatusBadRequest)
			return
		}

		page := queryParams.Get("page")
		pageNumber := 1
		if page != "" {
			var err error
			pageNumber, err = strconv.Atoi(page)
			if err != nil {
				http.Error(w, "page value must be integer", http.StatusBadRequest)
				return
			}

			if pageNumber < 1 {
				http.Error(w, "page value must be greater than 0", http.StatusBadRequest)
				return
			}
		}

		limit := queryParams.Get("limit")
		lim := 10
		if limit != "" {
			var err error
			lim, err = strconv.Atoi(limit)
			if err != nil {
				http.Error(w, "limit value must be integer", http.StatusBadRequest)
				return
			}

			if lim < 1 || lim > 30 {
				http.Error(w, "limit value must be greater than 0 or less than or equal to 30", http.StatusBadRequest)
				return
			}
		}

		ctx := r.Context()
		res, err := hm.svc.GetPvzInfo(ctx, startDate, endDate, int32(pageNumber), int32(lim))

		if err != nil {
			http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return

	default:
		http.Error(w, errors.ErrInvalidHtppMethod.Error(), http.StatusMethodNotAllowed)
		return
	}
}
