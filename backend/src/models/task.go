package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TaskStatus string
type TaskPriority string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"

	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          uuid.UUID    `db:"id" json:"id"`
	Title       string       `db:"title" json:"title"`
	Description *string      `db:"description" json:"description"`
	Status      TaskStatus   `db:"status" json:"status"`
	Priority    TaskPriority `db:"priority" json:"priority"`
	ProjectID   uuid.UUID    `db:"project_id" json:"project_id"`
	AssigneeID  *uuid.UUID   `db:"assignee_id" json:"assignee_id"`
	DueDate     *time.Time   `db:"due_date" json:"due_date"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updated_at"`
}

type TaskStore struct {
	DB *sqlx.DB
}

func (s *TaskStore) Create(title, description string, projectID uuid.UUID, status, priority TaskStatus, assigneeID *uuid.UUID, dueDate *time.Time) (*Task, error) {
	var t Task

	query := `
		INSERT INTO tasks (title, description, project_id, status, priority, assignee_id, due_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at`

	err := s.DB.QueryRowx(query, title, description, projectID, status, priority, assigneeID, dueDate).StructScan(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *TaskStore) GetByID(id uuid.UUID) (*Task, error) {
	var t Task
	query := `SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at FROM tasks WHERE id = $1`

	err := s.DB.Get(&t, query, id)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *TaskStore) GetByProjectID(projectID uuid.UUID) ([]Task, error) {
	var tasks []Task
	query := `SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at FROM tasks WHERE project_id = $1 ORDER BY created_at DESC`

	err := s.DB.Select(&tasks, query, projectID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskStore) Update(id uuid.UUID, title, description *string, status, priority *TaskStatus, assigneeID *uuid.UUID, dueDate *time.Time) (*Task, error) {
	var t Task

	query := `
		UPDATE tasks 
		SET title = COALESCE($1, title),
			description = COALESCE($2, description),
			status = COALESCE($3, status),
			priority = COALESCE($4, priority),
			assignee_id = $5,
			due_date = $6
		WHERE id = $7
		RETURNING id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at`

	err := s.DB.QueryRowx(query, title, description, status, priority, assigneeID, dueDate, id).StructScan(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *TaskStore) Delete(id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}

func (s *TaskStore) GetByIDAndProjectID(id, projectID uuid.UUID) (*Task, error) {
	var t Task
	query := `SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at FROM tasks WHERE id = $1 AND project_id = $2`

	err := s.DB.Get(&t, query, id, projectID)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
