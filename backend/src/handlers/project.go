package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/middleware"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/models"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
)

type ProjectHandler struct {
	Store     *models.ProjectStore
	TaskStore *models.TaskStore
}

type projectResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	OwnerID     string            `json:"owner_id"`
	CreatedAt   string            `json:"created_at"`
	Tasks       []taskResponse    `json:"tasks"`
	Pagination  *utils.Pagination `json:"pagination,omitempty"`
}

type listProjectsResponse struct {
	Projects   []projectResponse `json:"projects"`
	Pagination utils.Pagination  `json:"pagination"`
}

type createProjectRequest struct {
	Name        string `json:"name" validate:"required,min=3"`
	Description string `json:"description"`
}

type updateProjectRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3"`
	Description string `json:"description"`
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	page, limit, offset := utils.GetPagination(r)

	projects, err := h.Store.GetByOwnerIDPaginated(userID, limit, offset)
	if err != nil {
		slog.Error("Failed to fetch projects", "error", err, "user_id", userID)
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch projects")
		return
	}

	total, err := h.Store.CountByOwnerID(userID)
	if err != nil {
		slog.Error("Failed to count projects", "error", err, "user_id", userID)
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch projects")
		return
	}

	resp := listProjectsResponse{
		Projects:   make([]projectResponse, len(projects)),
		Pagination: utils.CalculatePagination(page, limit, total),
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
		validationErrors := err.(validator.ValidationErrors)
		fields := make(map[string]string)
		structType := reflect.TypeOf(req)
		for _, e := range validationErrors {
			jsonField := getJSONFieldName(structType, e.Field())
			switch e.Tag() {
			case "required":
				fields[jsonField] = "is required"
			case "min":
				fields[jsonField] = "must be at least " + e.Param() + " characters"
			default:
				fields[jsonField] = "is invalid"
			}
		}
		utils.WriteErrorWithFields(w, http.StatusBadRequest, "validation failed", fields)
		return
	}

	project, err := h.Store.Create(req.Name, req.Description, userID)
	if err != nil {
		slog.Error("Failed to create project", "error", err, "user_id", userID)
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

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	project, err := h.Store.GetByID(projectID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "project not found")
		return
	}

	if project.OwnerID != userID {
		utils.WriteError(w, http.StatusForbidden, "you don't have permission to access this project")
		return
	}

	page, limit, offset := utils.GetPagination(r)

	tasks, err := h.TaskStore.GetByProjectIDPaginated(projectID, nil, nil, limit, offset)
	if err != nil {
		slog.Error("Failed to fetch tasks for project", "error", err, "project_id", projectID)
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	total, err := h.TaskStore.CountByProjectID(projectID, nil, nil)
	if err != nil {
		slog.Error("Failed to count tasks for project", "error", err, "project_id", projectID)
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	taskResponses := make([]taskResponse, len(tasks))
	for i, t := range tasks {
		var assigneeIDStr, dueDateStr *string
		if t.AssigneeID != nil {
			s := t.AssigneeID.String()
			assigneeIDStr = &s
		}
		if t.DueDate != nil {
			s := t.DueDate.Format("2006-01-02")
			dueDateStr = &s
		}

		taskResponses[i] = taskResponse{
			ID:          t.ID.String(),
			Title:       t.Title,
			Description: stringOrNil(t.Description),
			Status:      string(t.Status),
			Priority:    string(t.Priority),
			ProjectID:   t.ProjectID.String(),
			AssigneeID:  assigneeIDStr,
			DueDate:     dueDateStr,
			CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   t.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	pagination := utils.CalculatePagination(page, limit, total)
	json.NewEncoder(w).Encode(projectResponse{
		ID:          project.ID.String(),
		Name:        project.Name,
		Description: project.Description,
		OwnerID:     project.OwnerID.String(),
		CreatedAt:   project.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Tasks:       taskResponses,
		Pagination:  &pagination,
	})
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	project, err := h.Store.GetByID(projectID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "project not found")
		return
	}

	if project.OwnerID != userID {
		utils.WriteError(w, http.StatusForbidden, "you don't have permission to access this project")
		return
	}

	var req updateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		fields := make(map[string]string)
		structType := reflect.TypeOf(req)
		for _, e := range validationErrors {
			jsonField := getJSONFieldName(structType, e.Field())
			switch e.Tag() {
			case "required":
				fields[jsonField] = "is required"
			case "min":
				fields[jsonField] = "must be at least " + e.Param() + " characters"
			default:
				fields[jsonField] = "is invalid"
			}
		}
		utils.WriteErrorWithFields(w, http.StatusBadRequest, "validation failed", fields)
		return
	}

	name := project.Name
	description := project.Description
	if req.Name != "" {
		name = req.Name
	}
	if req.Description != "" {
		description = req.Description
	}

	updatedProject, err := h.Store.Update(projectID, name, description)
	if err != nil {
		slog.Error("Failed to update project", "error", err, "project_id", projectID)
		utils.WriteError(w, http.StatusInternalServerError, "failed to update project")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projectResponse{
		ID:          updatedProject.ID.String(),
		Name:        updatedProject.Name,
		Description: updatedProject.Description,
		OwnerID:     updatedProject.OwnerID.String(),
		CreatedAt:   updatedProject.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	projectIDStr := chi.URLParam(r, "id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid project id")
		return
	}

	project, err := h.Store.GetByID(projectID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "project not found")
		return
	}

	if project.OwnerID != userID {
		utils.WriteError(w, http.StatusForbidden, "you don't have permission to access this project")
		return
	}

	if err := h.Store.Delete(projectID); err != nil {
		slog.Error("Failed to delete project", "error", err, "project_id", projectID)
		utils.WriteError(w, http.StatusInternalServerError, "failed to delete project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
