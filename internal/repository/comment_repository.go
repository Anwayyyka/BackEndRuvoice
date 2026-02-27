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
		SELECT c.id, c.user_id, c.track_id, c.text, c.created_at, c.updated_at, u.full_name
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.track_id = $1
		ORDER BY c.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		var userName string
		err := rows.Scan(&c.ID, &c.UserID, &c.TrackID, &c.Text, &c.CreatedAt, &c.UpdatedAt, &userName)
		if err != nil {
			return nil, err
		}
		// можно добавить поле UserName в домен, если нужно, но пока оставим так
		comments = append(comments, &c)
	}
	return comments, nil
}
