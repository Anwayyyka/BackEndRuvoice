// internal/repository/favorite_repository.go
package repository

import (
	"context"
	"errors"

	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FavoriteRepository struct {
	db *pgxpool.Pool
}

func NewFavoriteRepository(db *pgxpool.Pool) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Create(ctx context.Context, userID, trackID int) (*domain.Favorite, error) {
	query := `
		INSERT INTO favorites (user_id, track_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id, track_id) DO NOTHING
		RETURNING user_id, track_id, created_at
	`
	var fav domain.Favorite
	err := r.db.QueryRow(ctx, query, userID, trackID).Scan(&fav.UserID, &fav.TrackID, &fav.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil // уже существует
	}
	if err != nil {
		return nil, err
	}
	return &fav, nil
}

func (r *FavoriteRepository) Delete(ctx context.Context, userID, trackID int) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND track_id = $2`
	_, err := r.db.Exec(ctx, query, userID, trackID)
	return err
}

func (r *FavoriteRepository) GetByUser(ctx context.Context, userID int) ([]*domain.Favorite, error) {
	query := `SELECT user_id, track_id, created_at FROM favorites WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favs []*domain.Favorite
	for rows.Next() {
		var f domain.Favorite
		err := rows.Scan(&f.UserID, &f.TrackID, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		favs = append(favs, &f)
	}
	return favs, nil
}

// GetUserFavorites – алиас для совместимости
func (r *FavoriteRepository) GetUserFavorites(ctx context.Context, userID int) ([]*domain.Favorite, error) {
	return r.GetByUser(ctx, userID)
}
