package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/models"
	"github.com/kurayami07734/taskflow-aditya-ghidora/src/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	Store  *models.UserStore
	Config utils.Config
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

func (h *LoginHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.Store.GetByEmail(req.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid credentials")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, h.Config.JwtSecret)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	res := loginResponse{
		Token: token,
		User: userResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
