package repository

import (
	"context"
	"errors"
	"log"

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

type trackScanner interface {
	Scan(dest ...any) error
}

func scanTrack(scanner trackScanner) (*domain.Track, error) {
	track := &domain.Track{}
	var artistID *int
	err := scanner.Scan(
		&track.ID,
		&artistID, // сканируем во временную переменную
		&track.AlbumID,
		&track.Title,
		&track.Description,
		&track.AudioURL,
		&track.CoverURL,
		&track.Lyrics,
		&track.AuthorLyrics,
		&track.AuthorBeat,
		&track.AuthorComposer,
		&track.ReleaseDate,
		&track.PresaveURL,
		&track.Status,
		&track.PlaysCount,
		&track.LikesCount,
		&track.CreatedAt,
		&track.UpdatedAt,
		&track.ArtistName,
		&track.Duration,
		&track.ExternalID,
		&track.ExternalSource,
		&track.IsExternal,
	)
	track.ArtistID = artistID
	if err != nil {
		return nil, err
	}
	return track, nil
}

func (r *TrackRepository) Create(ctx context.Context, track *domain.Track) error {
	query := `
		INSERT INTO tracks (
			artist_id, album_id, title, description, audio_url, cover_url,
			lyrics, author_lyrics, author_beat, author_composer,
			release_date, presave_url, status, artist_name, duration,
			plays_count, likes_count, external_id, external_source, is_external
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		track.ArtistID, track.AlbumID, track.Title, track.Description,
		track.AudioURL, track.CoverURL, track.Lyrics, track.AuthorLyrics,
		track.AuthorBeat, track.AuthorComposer, track.ReleaseDate,
		track.PresaveURL, track.Status, track.ArtistName, track.Duration,
		track.PlaysCount, track.LikesCount,
		track.ExternalID, track.ExternalSource, track.IsExternal,
	).Scan(&track.ID, &track.CreatedAt, &track.UpdatedAt)
	return err
}

func (r *TrackRepository) GetByID(ctx context.Context, id int) (*domain.Track, error) {
	query := `
		SELECT id, artist_id, album_id, title, description, audio_url, cover_url,
		       lyrics, author_lyrics, author_beat, author_composer,
		       release_date, presave_url, status, plays_count, likes_count,
		       created_at, updated_at, artist_name, duration,
		       external_id, external_source, is_external
		FROM tracks
		WHERE id = $1 AND status = 'approved'
	`
	track, err := scanTrack(r.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return track, err
}

func (r *TrackRepository) GetByExternal(ctx context.Context, source string, externalID int) (*domain.Track, error) {
	query := `
		SELECT id, artist_id, album_id, title, description, audio_url, cover_url,
		       lyrics, author_lyrics, author_beat, author_composer,
		       release_date, presave_url, status, plays_count, likes_count,
		       created_at, updated_at, artist_name, duration,
		       external_id, external_source, is_external
		FROM tracks
		WHERE external_source = $1 AND external_id = $2
	`
	track, err := scanTrack(r.db.QueryRow(ctx, query, source, externalID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return track, err
}

func (r *TrackRepository) GetByIDs(ctx context.Context, ids []int) ([]*domain.Track, error) {
	if len(ids) == 0 {
		return []*domain.Track{}, nil
	}
	query := `
		SELECT id, artist_id, album_id, title, description, audio_url, cover_url,
		       lyrics, author_lyrics, author_beat, author_composer,
		       release_date, presave_url, status, plays_count, likes_count,
		       created_at, updated_at, artist_name, duration,
		       external_id, external_source, is_external
		FROM tracks
		WHERE id = ANY($1) AND status = 'approved'
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		log.Printf("GetByIDs query error: %v", err)
		return nil, err
	}
	defer rows.Close()
	var tracks []*domain.Track
	for rows.Next() {
		t, err := scanTrack(rows)
		if err != nil {
			log.Printf("GetByIDs scan error: %v", err)
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func (r *TrackRepository) ListApproved(ctx context.Context) ([]*domain.Track, error) {
	query := `
		SELECT id, artist_id, album_id, title, description, audio_url, cover_url,
		       lyrics, author_lyrics, author_beat, author_composer,
		       release_date, presave_url, status, plays_count, likes_count,
		       created_at, updated_at, artist_name, duration,
		       external_id, external_source, is_external
		FROM tracks
		WHERE status = 'approved'
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tracks []*domain.Track
	for rows.Next() {
		t, err := scanTrack(rows)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func (r *TrackRepository) ListByArtist(ctx context.Context, artistID int) ([]*domain.Track, error) {
	query := `
		SELECT id, artist_id, album_id, title, description, audio_url, cover_url,
		       lyrics, author_lyrics, author_beat, author_composer,
		       release_date, presave_url, status, plays_count, likes_count,
		       created_at, updated_at, artist_name, duration,
		       external_id, external_source, is_external
		FROM tracks
		WHERE artist_id = $1 AND status = 'approved'
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tracks []*domain.Track
	for rows.Next() {
		t, err := scanTrack(rows)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func (r *TrackRepository) ListPending(ctx context.Context) ([]*domain.Track, error) {
	query := `
		SELECT id, artist_id, album_id, title, description, audio_url, cover_url,
		       lyrics, author_lyrics, author_beat, author_composer,
		       release_date, presave_url, status, plays_count, likes_count,
		       created_at, updated_at, artist_name, duration,
		       external_id, external_source, is_external
		FROM tracks
		WHERE status = 'pending'
		ORDER BY created_at
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tracks []*domain.Track
	for rows.Next() {
		t, err := scanTrack(rows)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
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
	query := `
		UPDATE tracks
		SET likes_count = GREATEST(likes_count + $1, 0)
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, delta, trackID)
	return err
}
