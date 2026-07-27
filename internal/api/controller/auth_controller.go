package controller

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"gosync/internal/api/apperror"
	"gosync/internal/api/middleware"
	"gosync/internal/api/service"
	"gosync/internal/models"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) error {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if req.Email == "" || req.Password == "" {
		return apperror.BadRequest("email and password are required")
	}
	if len(req.Password) < 8 {
		return apperror.BadRequest("password must be at least 8 characters")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.Internal("failed to hash password", err)
	}

	user := &models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashed),
	}

	if err := c.authService.Register(user); err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(map[string]string{"id": user.ID.String()})
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) error {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if req.Email == "" || req.Password == "" {
		return apperror.BadRequest("email and password are required")
	}

	token, err := c.authService.Login(req.Email, req.Password)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(authResponse{Token: token})
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	if err := c.authService.Logout(userID); err != nil {
		return apperror.Internal("failed to log out", err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
