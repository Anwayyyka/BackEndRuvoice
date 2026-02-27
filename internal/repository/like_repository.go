package repository

import (
	"context"
	"errors"

	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LikeRepository struct {
	db *pgxpool.Pool
}

func NewLikeRepository(db *pgxpool.Pool) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) Create(ctx context.Context, userID, trackID int) (*domain.Like, error) {
	var like domain.Like
	query := `
		INSERT INTO likes (user_id, track_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, track_id) DO NOTHING
		RETURNING id, user_id, track_id, created_at
	`
	err := r.db.QueryRow(ctx, query, userID, trackID).Scan(&like.ID, &like.UserID, &like.TrackID, &like.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &like, nil
}

func (r *LikeRepository) Delete(ctx context.Context, userID, trackID int) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND track_id = $2`
	_, err := r.db.Exec(ctx, query, userID, trackID)
	return err
}

func (r *LikeRepository) GetByUserAndTrack(ctx context.Context, userID, trackID int) (*domain.Like, error) {
	var like domain.Like
	query := `SELECT id, user_id, track_id, created_at FROM likes WHERE user_id = $1 AND track_id = $2`
	err := r.db.QueryRow(ctx, query, userID, trackID).Scan(&like.ID, &like.UserID, &like.TrackID, &like.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &like, err
}

func (r *LikeRepository) CountByTrack(ctx context.Context, trackID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM likes WHERE track_id = $1`
	err := r.db.QueryRow(ctx, query, trackID).Scan(&count)
	return count, err
}
