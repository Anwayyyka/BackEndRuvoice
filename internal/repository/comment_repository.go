package repository

import (
	"context"

	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	query := `
		INSERT INTO comments (user_id, track_id, text)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, comment.UserID, comment.TrackID, comment.Text).
		Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	return err
}

func (r *CommentRepository) GetByTrack(ctx context.Context, trackID int) ([]*domain.Comment, error) {
	query := `
		SELECT id, user_id, track_id, text, created_at, updated_at
		FROM comments
		WHERE track_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.TrackID, &c.Text, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
