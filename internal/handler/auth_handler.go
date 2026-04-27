package handler

import (
	"encoding/json"
	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/dto"
	"go-loan-management-api/internal/response"
	"go-loan-management-api/internal/service"
	"net/http"
)

// AuthHandler handles auth-related HTTP requests.
type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req dto.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Register(r.Context(), service.RegisterInput{
		FullName:   req.FullName,
		Email:      req.Email,
		Password:   req.Password,
		Role:       req.Role,
		CustomerID: req.CustomerID,
	})
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "user registered successfully", result)
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.Login(r.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "login successful", result)
}

// Me handles GET /auth/me.
// func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := auth.GetClaims(r.Context())
// 	if !ok {
// 		response.Error(w, http.StatusUnauthorized, "unauthenticated")
// 		return
// 	}

// 	user, err := h.service.GetMe(r.Context(), claims.UserID)
// 	if err != nil {
// 		response.HandleError(w, err)
// 		return
// 	}

//		response.Success(w, http.StatusOK, "authenticated user fetched successfully", user)
//	}
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	authenticatedUser, ok := auth.GetAuthenticatedUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	user, err := h.service.GetMe(r.Context(), authenticatedUser.UserID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "authenticated user fetched successfully", user)
}
