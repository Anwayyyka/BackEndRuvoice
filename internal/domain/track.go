package domain

import "time"

type Track struct {
	ID             int        `json:"id"`
	ArtistID       *int       `json:"artist_id,omitempty"`
	ArtistName     string     `json:"artist_name"`
	AlbumID        *int       `json:"album_id,omitempty"`
	GenreID        *int       `json:"genre_id,omitempty"`
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
	Duration       int        `json:"duration"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Внешние треки (Jamendo и др.)
	ExternalID     *int    `json:"external_id,omitempty"`
	ExternalSource *string `json:"external_source,omitempty"`
	IsExternal     bool    `json:"is_external"`
}
