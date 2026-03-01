package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
)

type AlbumRepository struct {
	db *pgxpool.Pool
}

func NewAlbumRepository(db *pgxpool.Pool) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (r *AlbumRepository) Create(ctx context.Context, album *domain.Album) error {
	query := `
		INSERT INTO albums (artist_id, title, cover_url, release_date, presave_url, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		album.ArtistID, album.Title, album.CoverURL, album.ReleaseDate, album.PresaveURL, album.Status,
	).Scan(&album.ID, &album.CreatedAt, &album.UpdatedAt)
	return err
}

func (r *AlbumRepository) GetByID(ctx context.Context, id int) (*domain.Album, error) {
	query := `
		SELECT id, artist_id, title, cover_url, release_date, presave_url, status, created_at, updated_at
		FROM albums WHERE id = $1
	`
	album := &domain.Album{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&album.ID, &album.ArtistID, &album.Title, &album.CoverURL, &album.ReleaseDate,
		&album.PresaveURL, &album.Status, &album.CreatedAt, &album.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (r *AlbumRepository) ListByArtist(ctx context.Context, artistID int) ([]*domain.Album, error) {
	rows, err := r.db.Query(ctx, `SELECT id, artist_id, title, cover_url, release_date, presave_url, status, created_at, updated_at FROM albums WHERE artist_id = $1 ORDER BY created_at DESC`, artistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var albums []*domain.Album
	for rows.Next() {
		var a domain.Album
		err := rows.Scan(&a.ID, &a.ArtistID, &a.Title, &a.CoverURL, &a.ReleaseDate, &a.PresaveURL, &a.Status, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		albums = append(albums, &a)
	}
	return albums, nil
}

func (r *AlbumRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE albums SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}
