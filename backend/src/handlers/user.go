package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/models"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
)

type UserHandler struct {
	Store *models.UserStore
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.WriteBadRequest(w, "Invalid user ID")
		return
	}

	user, err := h.Store.GetByID(id)
	if err != nil {
		utils.WriteNotFound(w, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.WriteBadRequest(w, "Query parameter 'q' is required")
		return
	}

	users, err := h.Store.Search(query)
	if err != nil {
		utils.WriteInternalError(w, "Failed to search users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
