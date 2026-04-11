package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/models"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
	"golang.org/x/crypto/bcrypt"
)

type RegisterHandler struct {
	Store  *models.UserStore
	Config utils.Config
}
type registerRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type userResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type registerResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

func (h *RegisterHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	existingUser, err := h.Store.GetByEmail(req.Email)
	if existingUser != nil {
		utils.WriteError(w, http.StatusBadRequest, "User already registered!")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		log.Printf("Failed to generate password: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	user, err := h.Store.Create(req.Name, req.Email, string(hashed))
	if err != nil {
		log.Printf("Failed to save user: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, h.Config.JwtSecret)

	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	res := registerResponse{
		Token: token,
		User: userResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
