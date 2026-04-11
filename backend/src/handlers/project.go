package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kurayami07734/taskflow-aditya-ghidora/src/middleware"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/models"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
)

type ProjectHandler struct {
	Store *models.ProjectStore
}

type projectResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
	CreatedAt   string `json:"created_at"`
}

type listProjectsResponse struct {
	Projects []projectResponse `json:"projects"`
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	projects, err := h.Store.GetByOwnerID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch projects")
		return
	}

	resp := listProjectsResponse{
		Projects: make([]projectResponse, len(projects)),
	}

	for i, p := range projects {
		resp.Projects[i] = projectResponse{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			OwnerID:     p.OwnerID.String(),
			CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
