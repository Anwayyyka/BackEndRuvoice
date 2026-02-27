package domain

import "time"

type Like struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TrackID   int       `json:"track_id"`
	CreatedAt time.Time `json:"created_at"`
}
