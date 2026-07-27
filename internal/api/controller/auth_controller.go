package controller

import (
	"encoding/json"
	"net/http"

	"gosync/internal/api/apperror"
	"gosync/internal/api/dto"
	"gosync/internal/api/httputil"
	"gosync/internal/api/middleware"
	"gosync/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Registration details"
// @Success      201   {object}  dto.RegisterResponse
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /auth/register [post]
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) error {
	var req dto.RegisterRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	id, err := c.authService.Register(req.Email, req.Name, req.Password)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(dto.RegisterResponse{ID: id.String()})
}

// Login godoc
// @Summary      Log in
// @Description  Authenticates a user and returns a short-lived access token plus a long-lived refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest  true  "Login credentials"
// @Success      200   {object}  dto.AuthResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /auth/login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) error {
	var req dto.LoginRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		return err
	}

	accessToken, refreshToken, err := c.authService.Login(req.Email, req.Password)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// // Refresh godoc
// // @Summary      Refresh an access token
// // @Description  Exchanges a valid refresh token for a new access token + refresh token pair (rotation)
// // @Tags         auth
// // @Accept       json
// // @Produce      json
// // @Param        body  body      dto.RefreshRequest  true  "Refresh token"
// // @Success      200   {object}  dto.AuthResponse
// // @Failure      400   {object}  map[string]string
// // @Failure      401   {object}  map[string]string
// // @Router       /auth/refresh [post]
// func (c *AuthController) Refresh(w http.ResponseWriter, r *http.Request) error {
// 	var req dto.RefreshRequest
// 	if err := httputil.DecodeAndValidate(r, &req); err != nil {
// 		return err
// 	}

// 	accessToken, refreshToken, err := c.authService.Refresh(req.RefreshToken)
// 	if err != nil {
// 		return err
// 	}

// 	return json.NewEncoder(w).Encode(dto.AuthResponse{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	})
// }

// Logout godoc
// @Summary      Log out
// @Description  Revokes all of the authenticated user's active refresh tokens
// @Tags         auth
// @Success      200
// @Failure      401  {object}  map[string]string
// @Security     BearerAuth
// @Router       /auth/logout [post]
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.UserIDFromContext(r)
	if !ok {
		return apperror.Unauthorized("not authenticated")
	}

	if err := c.authService.Logout(userID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
