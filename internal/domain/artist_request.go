package domain

import "time"

type ArtistRequest struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	ArtistName       string    `json:"artist_name"`
	Bio              string    `json:"bio"`
	Status           string    `json:"status"`
	ModeratorID      *int      `json:"moderator_id,omitempty"`
	ModeratorComment *string   `json:"moderator_comment,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
