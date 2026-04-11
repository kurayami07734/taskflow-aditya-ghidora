package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
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

type createProjectRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description"`
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

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := h.Store.Create(req.Name, req.Description, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to create project")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(projectResponse{
		ID:          project.ID.String(),
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID.String(),
		CreatedAt:   project.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}
