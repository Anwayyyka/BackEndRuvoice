package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
)

type TrackLikeRepository struct {
	db *pgxpool.Pool
}

func NewTrackLikeRepository(db *pgxpool.Pool) *TrackLikeRepository {
	return &TrackLikeRepository{db: db}
}

func (r *TrackLikeRepository) Create(ctx context.Context, userID, trackID int) error {
	query := `INSERT INTO track_likes (user_id, track_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, userID, trackID)
	return err
}

func (r *TrackLikeRepository) Delete(ctx context.Context, userID, trackID int) error {
	query := `DELETE FROM track_likes WHERE user_id = $1 AND track_id = $2`
	_, err := r.db.Exec(ctx, query, userID, trackID)
	return err
}

func (r *TrackLikeRepository) CountByTrack(ctx context.Context, trackID int) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM track_likes WHERE track_id = $1`, trackID).Scan(&count)
	return count, err
}

func (r *TrackLikeRepository) GetByUserAndTrack(ctx context.Context, userID, trackID int) (*domain.TrackLike, error) {
	var like domain.TrackLike
	err := r.db.QueryRow(ctx, `SELECT user_id, track_id, created_at FROM track_likes WHERE user_id = $1 AND track_id = $2`, userID, trackID).
		Scan(&like.UserID, &like.TrackID, &like.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &like, err
}
