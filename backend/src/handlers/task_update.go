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

type updateTaskRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=3"`
	Description *string `json:"description"`
	Status      *string `json:"status" validate:"omitempty,oneof=todo in_progress done"`
	Priority    *string `json:"priority" validate:"omitempty,oneof=low medium high"`
	AssigneeID  *string `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}

// @Summary Update a task
// @Description Update an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param request body updateTaskRequest true "Update task request"
// @Success 200 {object} taskResponse
// @Failure 400 {object} utils.BadRequestError
// @Failure 401 {object} utils.UnauthorizedError
// @Failure 403 {object} utils.ForbiddenError
// @Failure 404 {object} utils.NotFoundError
// @Failure 500 {object} utils.InternalError
// @Router /tasks/{id} [patch]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteUnauthorized(w, "unauthorized")
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.WriteBadRequest(w, "invalid task id")
		return
	}

	task, err := h.TaskStore.GetByID(taskID)
	if err != nil {
		utils.WriteNotFound(w, "task not found")
		return
	}

	project, err := h.ProjectStore.GetByID(task.ProjectID)
	if err != nil {
		utils.WriteNotFound(w, "task not found")
		return
	}

	if project.OwnerID != userID {
		utils.WriteForbidden(w, "you don't have permission to access this task")
		return
	}

	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteBadRequest(w, "invalid request body")
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
			case "min":
				fields[jsonField] = "must be at least " + e.Param() + " characters"
			case "oneof":
				fields[jsonField] = "must be one of the allowed values"
			default:
				fields[jsonField] = "is invalid"
			}
		}
		utils.WriteErrorWithFields(w, http.StatusBadRequest, "validation failed", fields)
		return
	}

	var title *string
	if req.Title != nil {
		title = req.Title
	}

	var description *string
	if req.Description != nil {
		description = req.Description
	}

	var status *models.TaskStatus
	if req.Status != nil {
		s := models.TaskStatus(*req.Status)
		status = &s
	}

	var priority *models.TaskPriority
	if req.Priority != nil {
		p := models.TaskPriority(*req.Priority)
		priority = &p
	}

	var assigneeID *uuid.UUID
	if req.AssigneeID != nil {
		if *req.AssigneeID == "" {
			assigneeID = nil
		} else {
			parsed, err := uuid.Parse(*req.AssigneeID)
			if err == nil {
				assigneeID = &parsed
			}
		}
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		if *req.DueDate == "" {
			dueDate = nil
		} else {
			parsed, err := time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				parsed, err = time.Parse("2006-01-02", *req.DueDate)
			}
			if err == nil {
				dueDate = &parsed
			}
		}
	}

	updatedTask, err := h.TaskStore.Update(taskID, title, description, status, priority, assigneeID, dueDate)
	if err != nil {
		slog.Error("Failed to update task", "error", err, "task_id", taskID)
		utils.WriteInternalError(w, "failed to update task")
		return
	}

	var assigneeIDStr, dueDateStr *string
	if updatedTask.AssigneeID != nil {
		s := updatedTask.AssigneeID.String()
		assigneeIDStr = &s
	}
	if updatedTask.DueDate != nil {
		s := updatedTask.DueDate.Format("2006-01-02")
		dueDateStr = &s
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskResponse{
		ID:          updatedTask.ID.String(),
		Title:       updatedTask.Title,
		Description: stringOrNil(updatedTask.Description),
		Status:      string(updatedTask.Status),
		Priority:    string(updatedTask.Priority),
		ProjectID:   updatedTask.ProjectID.String(),
		AssigneeID:  assigneeIDStr,
		DueDate:     dueDateStr,
		CreatedAt:   updatedTask.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   updatedTask.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 204
// @Failure 401 {object} utils.UnauthorizedError
// @Failure 403 {object} utils.ForbiddenError
// @Failure 404 {object} utils.NotFoundError
// @Failure 500 {object} utils.InternalError
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.WriteUnauthorized(w, "unauthorized")
		return
	}

	taskIDStr := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		utils.WriteBadRequest(w, "invalid task id")
		return
	}

	task, err := h.TaskStore.GetByID(taskID)
	if err != nil {
		utils.WriteNotFound(w, "task not found")
		return
	}

	project, err := h.ProjectStore.GetByID(task.ProjectID)
	if err != nil {
		utils.WriteNotFound(w, "task not found")
		return
	}

	if project.OwnerID != userID {
		utils.WriteForbidden(w, "you don't have permission to access this task")
		return
	}

	if err := h.TaskStore.Delete(taskID); err != nil {
		utils.WriteInternalError(w, "failed to delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
