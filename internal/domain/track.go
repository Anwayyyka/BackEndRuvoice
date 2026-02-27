package domain

import "time"

type Track struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	ArtistID        int       `json:"artist_id"`
	ArtistName      string    `json:"artist_name"`
	GenreID         *int      `json:"genre_id,omitempty"`
	CoverURL        *string   `json:"cover_url,omitempty"`
	AudioURL        string    `json:"audio_url"`
	Duration        *int      `json:"duration,omitempty"`
	Description     *string   `json:"description,omitempty"`
	PlaysCount      int       `json:"plays_count"`
	LikesCount      int       `json:"likes_count"`
	Status          string    `json:"status"`
	RejectionReason *string   `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
