package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"reflect"

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

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login request"
// @Success 200 {object} loginResponse
// @Failure 400 {object} utils.BadRequestError
// @Router /auth/login [post]
func (h *LoginHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteBadRequest(w, "invalid request payload")
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
			case "email":
				fields[jsonField] = "must be a valid email"
			default:
				fields[jsonField] = "is invalid"
			}
		}
		utils.WriteErrorWithFields(w, http.StatusBadRequest, "validation failed", fields)
		return
	}

	user, err := h.Store.GetByEmail(req.Email)
	if err != nil {
		utils.WriteBadRequest(w, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.WriteBadRequest(w, "Invalid credentials")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, h.Config.JwtSecret)
	if err != nil {
		slog.Error("Failed to generate token", "error", err)
		utils.WriteInternalError(w, "Internal server error")
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
