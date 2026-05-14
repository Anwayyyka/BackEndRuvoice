package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/Anwayyyka/ruvoice-backend/internal/service"
)

type JamendoHandler struct {
	clientID     string
	trackService *service.TrackService
}

func NewJamendoHandler(trackService *service.TrackService) *JamendoHandler {
	return &JamendoHandler{
		clientID:     os.Getenv("JAMENDO_CLIENT_ID"),
		trackService: trackService,
	}
}

type jamendoTrack struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Duration   int    `json:"duration"`
	ArtistName string `json:"artist_name"`
	Image      string `json:"image"`
	Audio      string `json:"audio"`
}

type jamendoResponse struct {
	Headers struct {
		Status string `json:"status"`
		Code   int    `json:"code"`
	} `json:"headers"`
	Results []jamendoTrack `json:"results"`
}

func (h *JamendoHandler) GetTracks(w http.ResponseWriter, r *http.Request) {
	log.Println("JamendoHandler: started")
	if h.clientID == "" {
		log.Println("Jamendo client ID is empty")
		http.Error(w, "Jamendo client ID not configured", http.StatusInternalServerError)
		return
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "20"
	}
	url := fmt.Sprintf("https://api.jamendo.com/v3.0/tracks/?client_id=%s&format=json&limit=%s&include=musicinfo&order=popularity_total", h.clientID, limit)
	log.Printf("Requesting URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP request error: %v", err)
		http.Error(w, "Failed to fetch from Jamendo", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	log.Printf("Jamendo response status: %s", resp.Status)

	bodyBytes, _ := io.ReadAll(resp.Body)
	// log.Printf("Jamendo response body (first 500 chars): %s", string(bodyBytes[:min(500, len(bodyBytes))]))

	var jamResp jamendoResponse
	if err := json.Unmarshal(bodyBytes, &jamResp); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Failed to parse Jamendo response", http.StatusInternalServerError)
		return
	}
	if jamResp.Headers.Status != "success" {
		log.Printf("Jamendo API error: %+v", jamResp.Headers)
		http.Error(w, "Jamendo API error", http.StatusInternalServerError)
		return
	}

	resultTracks := make([]*domain.Track, 0, len(jamResp.Results))
	for _, jt := range jamResp.Results {
		trackID, err := strconv.Atoi(jt.ID)
		if err != nil {
			log.Printf("Failed to convert track ID %s: %v", jt.ID, err)
			continue
		}
		coverURL := jt.Image
		newTrack := &domain.Track{
			Title:      jt.Name,
			ArtistName: jt.ArtistName,
			Duration:   jt.Duration,
			CoverURL:   &coverURL,
			AudioURL:   jt.Audio,
			Status:     "approved",
		}
		saved, err := h.trackService.GetOrCreateExternalTrack(r.Context(), "jamendo", trackID, newTrack)
		if err != nil {
			log.Printf("Failed to save track %s: %v", jt.Name, err)
			continue
		}
		resultTracks = append(resultTracks, saved)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultTracks)
}
