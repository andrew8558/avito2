package handler_manager

import (
	"avito2/internal/errors"
	"avito2/internal/model"
	"encoding/json"
	"net/http"
)

func (hm *HandlerManager) DummyLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errors.ErrInvalidHtppMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	var req model.DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, errors.ErrInvalidJson.Error(), http.StatusBadRequest)
		return
	}

	if !req.Role.IsValid() {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	token, err := hm.jwtGen.GenerateJWT(string(req.Role))
	if err != nil {
		http.Error(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
