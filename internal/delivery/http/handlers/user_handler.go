package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Anwayyyka/ruvoice-backend/internal/delivery/http/middleware"
	"github.com/Anwayyyka/ruvoice-backend/internal/service"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type updateProfileRequest struct {
	FullName *string `json:"full_name,omitempty"`
	Bio      *string `json:"bio,omitempty"`
	Avatar   *string `json:"avatar_url,omitempty"`
	Banner   *string `json:"banner_url,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
	Vk       *string `json:"vk,omitempty"`
	Youtube  *string `json:"youtube,omitempty"`
	Website  *string `json:"website,omitempty"`
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	updates := make(map[string]interface{})
	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.Avatar != nil {
		updates["avatar_url"] = *req.Avatar
	}
	if req.Banner != nil {
		updates["banner_url"] = *req.Banner
	}
	if req.Telegram != nil {
		updates["telegram"] = *req.Telegram
	}
	if req.Vk != nil {
		updates["vk"] = *req.Vk
	}
	if req.Youtube != nil {
		updates["youtube"] = *req.Youtube
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}
	user, err := h.userService.UpdateProfile(r.Context(), userID, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToUserDTO(user))
}

type requestArtistRequest struct {
	ArtistName string `json:"artist_name"`
	Bio        string `json:"bio"`
}

func (h *UserHandler) RequestArtist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req requestArtistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := h.userService.RequestArtist(r.Context(), userID, req.ArtistName, req.Bio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "request sent"})
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	user, err := h.userService.GetByEmail(r.Context(), email)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ToUserDTO(user))
}
