package domain

import "time"

type Track struct {
	ID             int        `json:"id"`
	ArtistID       int        `json:"artist_id"`
	ArtistName     string     `json:"artist_name"` // добавлено
	AlbumID        *int       `json:"album_id,omitempty"`
	GenreID        *int       `json:"genre_id,omitempty"` // добавлено
	Title          string     `json:"title"`
	Description    *string    `json:"description,omitempty"`
	AudioURL       string     `json:"audio_url"`
	CoverURL       *string    `json:"cover_url,omitempty"`
	Lyrics         *string    `json:"lyrics,omitempty"`
	AuthorLyrics   *string    `json:"author_lyrics,omitempty"`
	AuthorBeat     *string    `json:"author_beat,omitempty"`
	AuthorComposer *string    `json:"author_composer,omitempty"`
	ReleaseDate    *time.Time `json:"release_date,omitempty"`
	PresaveURL     *string    `json:"presave_url,omitempty"`
	Status         string     `json:"status"`
	PlaysCount     int        `json:"plays_count"`
	LikesCount     int        `json:"likes_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
