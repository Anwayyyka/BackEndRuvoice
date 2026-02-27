package domain

import "time"

type User struct {
	ID              int       `json:"id"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	FullName        *string   `json:"full_name,omitempty"`
	ArtistName      *string   `json:"artist_name,omitempty"` // добавлено
	Role            string    `json:"role"`
	AvatarURL       *string   `json:"avatar_url,omitempty"`
	BannerURL       *string   `json:"banner_url,omitempty"`
	Bio             *string   `json:"bio,omitempty"`
	Telegram        *string   `json:"telegram,omitempty"`
	Vk              *string   `json:"vk,omitempty"`
	Youtube         *string   `json:"youtube,omitempty"`
	Website         *string   `json:"website,omitempty"`
	ArtistRequested bool      `json:"artist_requested"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
