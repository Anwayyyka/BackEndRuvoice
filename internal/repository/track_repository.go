package repository

import (
	"context"
	"errors"
	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackRepository struct {
	db *pgxpool.Pool
}

func NewTrackRepository(db *pgxpool.Pool) *TrackRepository {
	return &TrackRepository{db: db}
}

func (r *TrackRepository) Create(ctx context.Context, track *domain.Track) error {
	query := `
		INSERT INTO tracks (artist_id, album_id, title, description, audio_url, cover_url,
		                    lyrics, author_lyrics, author_beat, author_composer,
		                    release_date, presave_url, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		track.ArtistID, track.AlbumID, track.Title, track.Description,
		track.AudioURL, track.CoverURL, track.Lyrics, track.AuthorLyrics,
		track.AuthorBeat, track.AuthorComposer, track.ReleaseDate,
		track.PresaveURL, track.Status,
	).Scan(&track.ID, &track.CreatedAt, &track.UpdatedAt)
	return err
}

func (r *TrackRepository) GetByID(ctx context.Context, id int) (*domain.Track, error) {
	query := `
		SELECT id, artist_id, album_id, title, description, audio_url, cover_url,
		       lyrics, author_lyrics, author_beat, author_composer,
		       release_date, presave_url, status, plays_count, likes_count,
		       created_at, updated_at
		FROM tracks WHERE id = $1
	`
	track := &domain.Track{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&track.ID, &track.ArtistID, &track.AlbumID, &track.Title, &track.Description,
		&track.AudioURL, &track.CoverURL, &track.Lyrics, &track.AuthorLyrics,
		&track.AuthorBeat, &track.AuthorComposer, &track.ReleaseDate,
		&track.PresaveURL, &track.Status, &track.PlaysCount, &track.LikesCount,
		&track.CreatedAt, &track.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return track, err
}

func (r *TrackRepository) ListApproved(ctx context.Context) ([]*domain.Track, error) {
	rows, err := r.db.Query(ctx, `SELECT id, artist_id, album_id, title, description, audio_url, cover_url, lyrics, author_lyrics, author_beat, author_composer, release_date, presave_url, status, plays_count, likes_count, created_at, updated_at FROM tracks WHERE status = 'approved' ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tracks []*domain.Track
	for rows.Next() {
		var t domain.Track
		err := rows.Scan(
			&t.ID, &t.ArtistID, &t.AlbumID, &t.Title, &t.Description,
			&t.AudioURL, &t.CoverURL, &t.Lyrics, &t.AuthorLyrics,
			&t.AuthorBeat, &t.AuthorComposer, &t.ReleaseDate,
			&t.PresaveURL, &t.Status, &t.PlaysCount, &t.LikesCount,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &t)
	}
	return tracks, nil
}

func (r *TrackRepository) ListByArtist(ctx context.Context, artistID int) ([]*domain.Track, error) {
	rows, err := r.db.Query(ctx, `SELECT id, artist_id, album_id, title, description, audio_url, cover_url, lyrics, author_lyrics, author_beat, author_composer, release_date, presave_url, status, plays_count, likes_count, created_at, updated_at FROM tracks WHERE artist_id = $1 ORDER BY created_at DESC`, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tracks []*domain.Track
	for rows.Next() {
		var t domain.Track
		err := rows.Scan(
			&t.ID, &t.ArtistID, &t.AlbumID, &t.Title, &t.Description,
			&t.AudioURL, &t.CoverURL, &t.Lyrics, &t.AuthorLyrics,
			&t.AuthorBeat, &t.AuthorComposer, &t.ReleaseDate,
			&t.PresaveURL, &t.Status, &t.PlaysCount, &t.LikesCount,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &t)
	}
	return tracks, nil
}

func (r *TrackRepository) ListPending(ctx context.Context) ([]*domain.Track, error) {
	rows, err := r.db.Query(ctx, `SELECT id, artist_id, album_id, title, description, audio_url, cover_url, lyrics, author_lyrics, author_beat, author_composer, release_date, presave_url, status, plays_count, likes_count, created_at, updated_at FROM tracks WHERE status = 'pending' ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tracks []*domain.Track
	for rows.Next() {
		var t domain.Track
		err := rows.Scan(
			&t.ID, &t.ArtistID, &t.AlbumID, &t.Title, &t.Description,
			&t.AudioURL, &t.CoverURL, &t.Lyrics, &t.AuthorLyrics,
			&t.AuthorBeat, &t.AuthorComposer, &t.ReleaseDate,
			&t.PresaveURL, &t.Status, &t.PlaysCount, &t.LikesCount,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &t)
	}
	return tracks, nil
}

func (r *TrackRepository) Update(ctx context.Context, track *domain.Track) error {
	query := `
		UPDATE tracks SET
			title = COALESCE($2, title),
			description = COALESCE($3, description),
			cover_url = COALESCE($4, cover_url),
			lyrics = COALESCE($5, lyrics),
			author_lyrics = COALESCE($6, author_lyrics),
			author_beat = COALESCE($7, author_beat),
			author_composer = COALESCE($8, author_composer),
			release_date = $9,
			presave_url = COALESCE($10, presave_url),
			status = $11,
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		track.ID, track.Title, track.Description, track.CoverURL,
		track.Lyrics, track.AuthorLyrics, track.AuthorBeat, track.AuthorComposer,
		track.ReleaseDate, track.PresaveURL, track.Status,
	)
	return err
}

func (r *TrackRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE tracks SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *TrackRepository) IncrementPlays(ctx context.Context, id int) error {
	query := `UPDATE tracks SET plays_count = plays_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
func (r *TrackRepository) UpdateLikesCount(ctx context.Context, trackID int, delta int) error {
	query := `UPDATE tracks SET likes_count = likes_count + $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, delta, trackID)
	return err
}
