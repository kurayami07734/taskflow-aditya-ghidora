package handlers

import (
	"github.com/google/uuid"
)

type userResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type authResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}
