package domain

import "time"

type Album struct {
	ID          int        `json:"id"`
	ArtistID    int        `json:"artist_id"`
	Title       string     `json:"title"`
	CoverURL    *string    `json:"cover_url,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	PresaveURL  *string    `json:"presave_url,omitempty"`
	Status      string     `json:"status"` // pending, approved, rejected
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
