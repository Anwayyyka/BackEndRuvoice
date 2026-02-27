package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Anwayyyka/ruvoice-backend/internal/service"
)

type ModerationHandler struct {
	trackService *service.TrackService
}

func NewModerationHandler(trackService *service.TrackService) *ModerationHandler {
	return &ModerationHandler{trackService: trackService}
}

func (h *ModerationHandler) ListPending(w http.ResponseWriter, r *http.Request) {
	tracks, err := h.trackService.GetPending(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

type approveRequest struct {
	TrackID int `json:"track_id"`
}

func (h *ModerationHandler) Approve(w http.ResponseWriter, r *http.Request) {
	var req approveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := h.trackService.Approve(r.Context(), req.TrackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type rejectRequest struct {
	TrackID int    `json:"track_id"`
	Reason  string `json:"reason"`
}

func (h *ModerationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	var req rejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err := h.trackService.Reject(r.Context(), req.TrackID, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
