package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"` // Hashed; excluded from JSON
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UserStore struct {
	DB *sqlx.DB
}

func (s *UserStore) Create(name, email, hashedPwd string) (*User, error) {
	var u User

	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email`

	err := s.DB.QueryRowx(query, name, email, hashedPwd).StructScan(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) GetByEmail(email string) (*User, error) {
	var u User
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`

	err := s.DB.Get(&u, query, email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) GetByID(id uuid.UUID) (*User, error) {
	var u User
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`

	err := s.DB.Get(&u, query, id)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) Search(query string) ([]User, error) {
	var users []User
	sqlQuery := `SELECT id, name, email, created_at FROM users WHERE name ILIKE $1 OR email ILIKE $1 LIMIT 20`
	searchPattern := "%" + query + "%"

	err := s.DB.Select(&users, sqlQuery, searchPattern)
	if err != nil {
		return nil, err
	}

	return users, nil
}
