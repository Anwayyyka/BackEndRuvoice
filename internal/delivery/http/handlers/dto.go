package handlers

import "github.com/Anwayyyka/ruvoice-backend/internal/domain"

type UserDTO struct {
	ID         int     `json:"id"`
	Email      string  `json:"email"`
	FullName   *string `json:"full_name,omitempty"`
	ArtistName *string `json:"artist_name,omitempty"`

	Role      string  `json:"role"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	BannerURL *string `json:"banner_url,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	Telegram  *string `json:"telegram,omitempty"`
	Vk        *string `json:"vk,omitempty"`
	Youtube   *string `json:"youtube,omitempty"`
	Website   *string `json:"website,omitempty"`
}

func ToUserDTO(u *domain.User) *UserDTO {
	if u == nil {
		return nil
	}
	return &UserDTO{
		ID:         u.ID,
		Email:      u.Email,
		FullName:   u.FullName,
		ArtistName: u.ArtistName,
		Role:       u.Role,
		AvatarURL:  u.AvatarURL,
		BannerURL:  u.BannerURL,
		Bio:        u.Bio,
		Telegram:   u.Telegram,
		Vk:         u.Vk,
		Youtube:    u.Youtube,
		Website:    u.Website,
	}
}
