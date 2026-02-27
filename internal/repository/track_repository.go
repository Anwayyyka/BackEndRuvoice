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
		INSERT INTO tracks (title, artist_id, artist_name, genre_id, cover_url, audio_url, duration, description, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		track.Title, track.ArtistID, track.ArtistName, track.GenreID,
		track.CoverURL, track.AudioURL, track.Duration, track.Description, track.Status,
	).Scan(&track.ID, &track.CreatedAt, &track.UpdatedAt)
	return err
}

func (r *TrackRepository) GetByID(ctx context.Context, id int) (*domain.Track, error) {
	query := `
		SELECT id, title, artist_id, artist_name, genre_id, cover_url, audio_url, duration,
		       description, plays_count, likes_count, status, rejection_reason, created_at, updated_at
		FROM tracks WHERE id = $1
	`
	track := &domain.Track{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&track.ID, &track.Title, &track.ArtistID, &track.ArtistName, &track.GenreID,
		&track.CoverURL, &track.AudioURL, &track.Duration, &track.Description,
		&track.PlaysCount, &track.LikesCount, &track.Status, &track.RejectionReason,
		&track.CreatedAt, &track.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return track, err
}

func (r *TrackRepository) ListApproved(ctx context.Context) ([]*domain.Track, error) {
	query := `
        SELECT id, title, artist_id, artist_name, genre_id, cover_url, audio_url, duration,
               description, plays_count, likes_count, status, rejection_reason, created_at, updated_at
        FROM tracks WHERE status = 'approved' ORDER BY created_at DESC
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*domain.Track
	for rows.Next() {
		var t domain.Track
		err := rows.Scan(
			&t.ID, &t.Title, &t.ArtistID, &t.ArtistName, &t.GenreID,
			&t.CoverURL, &t.AudioURL, &t.Duration, &t.Description,
			&t.PlaysCount, &t.LikesCount, &t.Status, &t.RejectionReason,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tracks, nil
}

func (r *TrackRepository) ListByArtist(ctx context.Context, artistID int) ([]*domain.Track, error) {
	query := `
		SELECT id, title, artist_id, artist_name, genre_id, cover_url, audio_url, duration,
		       description, plays_count, likes_count, status, rejection_reason, created_at, updated_at
		FROM tracks WHERE artist_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*domain.Track
	for rows.Next() {
		var t domain.Track
		err := rows.Scan(
			&t.ID, &t.Title, &t.ArtistID, &t.ArtistName, &t.GenreID,
			&t.CoverURL, &t.AudioURL, &t.Duration, &t.Description,
			&t.PlaysCount, &t.LikesCount, &t.Status, &t.RejectionReason,
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
	query := `
		SELECT id, title, artist_id, artist_name, genre_id, cover_url, audio_url, duration,
		       description, plays_count, likes_count, status, rejection_reason, created_at, updated_at
		FROM tracks WHERE status = 'pending' ORDER BY created_at
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []*domain.Track
	for rows.Next() {
		var t domain.Track
		err := rows.Scan(
			&t.ID, &t.Title, &t.ArtistID, &t.ArtistName, &t.GenreID,
			&t.CoverURL, &t.AudioURL, &t.Duration, &t.Description,
			&t.PlaysCount, &t.LikesCount, &t.Status, &t.RejectionReason,
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
			artist_name = COALESCE($3, artist_name),
			genre_id = $4,
			cover_url = COALESCE($5, cover_url),
			description = COALESCE($6, description),
			status = $7,
			rejection_reason = $8,
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		track.ID, track.Title, track.ArtistName, track.GenreID,
		track.CoverURL, track.Description, track.Status, track.RejectionReason,
	)
	return err
}

func (r *TrackRepository) IncrementPlays(ctx context.Context, trackID int) error {
	query := `UPDATE tracks SET plays_count = plays_count + 1 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, trackID)
	return err
}

func (r *TrackRepository) UpdateLikesCount(ctx context.Context, trackID int, delta int) error {
	query := `UPDATE tracks SET likes_count = likes_count + $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, delta, trackID)
	return err
}

func (r *TrackRepository) Approve(ctx context.Context, trackID int) error {
	query := `UPDATE tracks SET status = 'approved', updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, trackID)
	return err
}

func (r *TrackRepository) Reject(ctx context.Context, trackID int, reason string) error {
	query := `UPDATE tracks SET status = 'rejected', rejection_reason = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, reason, trackID)
	return err
}
