package repository

import (
	"context"
	"errors"

	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtistRequestRepository struct {
	db *pgxpool.Pool
}

func NewArtistRequestRepository(db *pgxpool.Pool) *ArtistRequestRepository {
	return &ArtistRequestRepository{db: db}
}

func (r *ArtistRequestRepository) Create(ctx context.Context, req *domain.ArtistRequest) error {
	query := `
		INSERT INTO artist_requests (user_id, artist_name, bio, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		req.UserID, req.ArtistName, req.Bio, req.Status,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
	return err
}

func (r *ArtistRequestRepository) GetByID(ctx context.Context, id int) (*domain.ArtistRequest, error) {
	query := `
		SELECT id, user_id, artist_name, bio, status, moderator_id, moderator_comment, created_at, updated_at
		FROM artist_requests WHERE id = $1
	`
	var req domain.ArtistRequest
	err := r.db.QueryRow(ctx, query, id).Scan(
		&req.ID, &req.UserID, &req.ArtistName, &req.Bio, &req.Status,
		&req.ModeratorID, &req.ModeratorComment, &req.CreatedAt, &req.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &req, err
}

func (r *ArtistRequestRepository) GetByUserID(ctx context.Context, userID int) (*domain.ArtistRequest, error) {
	query := `
		SELECT id, user_id, artist_name, bio, status, moderator_id, moderator_comment, created_at, updated_at
		FROM artist_requests WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1
	`
	var req domain.ArtistRequest
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&req.ID, &req.UserID, &req.ArtistName, &req.Bio, &req.Status,
		&req.ModeratorID, &req.ModeratorComment, &req.CreatedAt, &req.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &req, err
}

func (r *ArtistRequestRepository) ListByStatus(ctx context.Context, status string) ([]*domain.ArtistRequest, error) {
	query := `
		SELECT id, user_id, artist_name, bio, status, moderator_id, moderator_comment, created_at, updated_at
		FROM artist_requests WHERE status = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*domain.ArtistRequest
	for rows.Next() {
		var req domain.ArtistRequest
		err := rows.Scan(
			&req.ID, &req.UserID, &req.ArtistName, &req.Bio, &req.Status,
			&req.ModeratorID, &req.ModeratorComment, &req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, &req)
	}
	return requests, nil
}

func (r *ArtistRequestRepository) UpdateStatus(ctx context.Context, id int, status string, moderatorID int, comment string) error {
	query := `
		UPDATE artist_requests
		SET status = $1, moderator_id = $2, moderator_comment = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.Exec(ctx, query, status, moderatorID, comment, id)
	return err
}
