package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Anwayyyka/ruvoice-backend/internal/delivery/http/middleware"
	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/Anwayyyka/ruvoice-backend/internal/service"
	"github.com/Anwayyyka/ruvoice-backend/internal/upload"
	"github.com/go-chi/chi/v5"
)

type TrackHandler struct {
	trackService *service.TrackService
	userService  *service.UserService
}

func NewTrackHandler(trackService *service.TrackService, userService *service.UserService) *TrackHandler {
	return &TrackHandler{trackService: trackService, userService: userService}
}

func (h *TrackHandler) ListApproved(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.trackService.GetApproved(r.Context())
	if err != nil {
		log.Printf("ERROR in ListApproved: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

func (h *TrackHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	track, err := h.trackService.GetByID(r.Context(), id)
	if err != nil || track == nil {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(track)
}

func (h *TrackHandler) GetByArtist(w http.ResponseWriter, r *http.Request) {
	artistIDStr := chi.URLParam(r, "artistId")
	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}
	tracks, err := h.trackService.GetByArtist(r.Context(), artistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

func (h *TrackHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	artistName := r.FormValue("artist_name")
	genreIDStr := r.FormValue("genre_id")
	description := r.FormValue("description")

	if title == "" || artistName == "" {
		http.Error(w, "Title and artist name are required", http.StatusBadRequest)
		return
	}

	var genreID *int
	if genreIDStr != "" {
		if gid, err := strconv.Atoi(genreIDStr); err == nil {
			genreID = &gid
		}
	}

	audioFile, audioHeader, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Audio file is required", http.StatusBadRequest)
		return
	}
	defer audioFile.Close()
	audioURL, err := upload.SaveFile(audioFile, audioHeader, "audio")
	if err != nil {
		http.Error(w, "Failed to save audio", http.StatusInternalServerError)
		return
	}

	var coverURL *string
	coverFile, coverHeader, err := r.FormFile("cover")
	if err == nil {
		defer coverFile.Close()
		url, err := upload.SaveFile(coverFile, coverHeader, "image")
		if err == nil {
			coverURL = &url
		}
	}

	track := &domain.Track{
		Title:       title,
		ArtistID:    &userID,
		ArtistName:  artistName,
		GenreID:     genreID,
		CoverURL:    coverURL,
		AudioURL:    audioURL,
		Description: &description,
		Status:      "pending",
	}

	err = h.trackService.Create(r.Context(), track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(track)
}

func (h *TrackHandler) Play(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = h.trackService.IncrementPlays(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *TrackHandler) Like(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	trackID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = h.trackService.Like(r.Context(), userID, trackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "liked"})
}

func (h *TrackHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	trackID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = h.trackService.Unlike(r.Context(), userID, trackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "unliked"})
}

func (h *TrackHandler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	trackID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = h.trackService.AddFavorite(r.Context(), userID, trackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"track_id": trackID,
		"status":   "added",
	})
}

func (h *TrackHandler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := chi.URLParam(r, "id")
	trackID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = h.trackService.RemoveFavorite(r.Context(), userID, trackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "removed"})
}

func (h *TrackHandler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	favs, err := h.trackService.GetUserFavorites(r.Context(), userID)
	if err != nil {
		log.Printf("GetUserFavorites service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(favs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*domain.Track{})
		return
	}
	trackIDs := make([]int, len(favs))
	for i, f := range favs {
		trackIDs[i] = f.TrackID
	}
	tracks, err := h.trackService.GetTracksByIDs(r.Context(), trackIDs)
	if err != nil {
		log.Printf("GetTracksByIDs error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}
