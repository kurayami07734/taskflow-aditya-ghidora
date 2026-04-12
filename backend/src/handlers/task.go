package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/middleware"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/models"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
)

type TaskHandler struct {
	TaskStore    *models.TaskStore
	ProjectStore *models.ProjectStore
}

type taskResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Priority    string  `json:"priority"`
	ProjectID   string  `json:"project_id"`
	AssigneeID  *string `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type createTaskRequest struct {
	Title       string  `json:"title" validate:"required,min=3"`
	Description string  `json:"description"`
	Priority    string  `json:"priority" validate:"omitempty,oneof=low medium high"`
	AssigneeID  *string `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
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

	project, err := h.ProjectStore.GetByID(projectID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "project not found")
		return
	}

	if project.OwnerID != userID {
		utils.WriteError(w, http.StatusForbidden, "you don't have permission to access this project")
		return
	}

	var req createTaskRequest
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
			case "oneof":
				fields[jsonField] = "must be one of: low, medium, high"
			default:
				fields[jsonField] = "is invalid"
			}
		}
		utils.WriteErrorWithFields(w, http.StatusBadRequest, "validation failed", fields)
		return
	}

	priority := models.PriorityMedium
	if req.Priority != "" {
		priority = models.TaskPriority(req.Priority)
	}

	var assigneeID *uuid.UUID
	if req.AssigneeID != nil && *req.AssigneeID != "" {
		parsed, err := uuid.Parse(*req.AssigneeID)
		if err == nil {
			assigneeID = &parsed
		}
	}

	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err == nil {
			dueDate = &parsed
		}
	}

	task, err := h.TaskStore.Create(req.Title, req.Description, projectID, models.StatusTodo, models.TaskPriority(priority), assigneeID, dueDate)
	if err != nil {
		slog.Error("Failed to create task", "error", err, "project_id", projectID)
		utils.WriteError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	var assigneeIDStr, dueDateStr *string
	if task.AssigneeID != nil {
		s := task.AssigneeID.String()
		assigneeIDStr = &s
	}
	if task.DueDate != nil {
		s := task.DueDate.Format("2006-01-02")
		dueDateStr = &s
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(taskResponse{
		ID:          task.ID.String(),
		Title:       task.Title,
		Description: stringOrNil(task.Description),
		Status:      string(task.Status),
		Priority:    string(task.Priority),
		ProjectID:   task.ProjectID.String(),
		AssigneeID:  assigneeIDStr,
		DueDate:     dueDateStr,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

type listTasksResponse struct {
	Tasks      []taskResponse   `json:"tasks"`
	Pagination utils.Pagination `json:"pagination"`
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
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

	project, err := h.ProjectStore.GetByID(projectID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "project not found")
		return
	}

	if project.OwnerID != userID {
		utils.WriteError(w, http.StatusForbidden, "you don't have permission to access this project")
		return
	}

	page, limit, offset := utils.GetPagination(r)

	statusFilter := r.URL.Query().Get("status")
	assigneeFilter := r.URL.Query().Get("assignee")

	var status *models.TaskStatus
	if statusFilter != "" {
		s := models.TaskStatus(statusFilter)
		status = &s
	}

	var assigneeID *uuid.UUID
	if assigneeFilter != "" {
		parsed, err := uuid.Parse(assigneeFilter)
		if err == nil {
			assigneeID = &parsed
		}
	}

	tasks, err := h.TaskStore.GetByProjectIDPaginated(projectID, status, assigneeID, limit, offset)
	if err != nil {
		slog.Error("Failed to fetch tasks", "error", err, "project_id", projectID)
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	total, err := h.TaskStore.CountByProjectID(projectID, status, assigneeID)
	if err != nil {
		slog.Error("Failed to count tasks", "error", err, "project_id", projectID)
		utils.WriteError(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	resp := listTasksResponse{
		Tasks:      make([]taskResponse, len(tasks)),
		Pagination: utils.CalculatePagination(page, limit, total),
	}

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

		resp.Tasks[i] = taskResponse{
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
	json.NewEncoder(w).Encode(resp)
}

func stringOrNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
