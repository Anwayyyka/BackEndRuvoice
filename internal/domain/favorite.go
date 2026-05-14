package domain

import "time"

type Favorite struct {
	UserID    int       `json:"user_id"`
	TrackID   int       `json:"track_id"`
	CreatedAt time.Time `json:"created_at"`
}
