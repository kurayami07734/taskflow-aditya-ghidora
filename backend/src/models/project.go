package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Project struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	OwnerID     uuid.UUID `db:"owner_id" json:"owner_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type ProjectStore struct {
	DB *sqlx.DB
}

func (s *ProjectStore) Create(name, description string, ownerID uuid.UUID) (*Project, error) {
	var p Project

	query := `
		INSERT INTO projects (name, description, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, owner_id, created_at`

	err := s.DB.QueryRowx(query, name, description, ownerID).StructScan(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *ProjectStore) GetByID(id uuid.UUID) (*Project, error) {
	var p Project
	query := `SELECT id, name, description, owner_id, created_at FROM projects WHERE id = $1`

	err := s.DB.Get(&p, query, id)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *ProjectStore) GetByOwnerID(ownerID uuid.UUID) ([]Project, error) {
	var projects []Project
	query := `SELECT id, name, description, owner_id, created_at FROM projects WHERE owner_id = $1 ORDER BY created_at DESC`

	err := s.DB.Select(&projects, query, ownerID)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
