package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Anwayyyka/ruvoice-backend/internal/delivery/http/middleware"
	"github.com/Anwayyyka/ruvoice-backend/internal/service"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	User  *UserDTO `json:"user"`
	Token string   `json:"token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Register decode error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user, token, err := h.userService.Register(r.Context(), service.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		log.Printf("Register service error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest) // здесь пробрасывается русское сообщение
		return
	}
	resp := authResponse{
		User:  ToUserDTO(user),
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Login decode error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user, token, err := h.userService.Login(r.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Login service error: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized) // здесь пробрасывается русское сообщение
		return
	}
	resp := authResponse{
		User:  ToUserDTO(user),
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		log.Printf("GetProfile error: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToUserDTO(user))
}
