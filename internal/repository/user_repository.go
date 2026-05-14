package repository

import (
	"context"
	"errors"

	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (*domain.User, error) {
	user := &domain.User{}
	err := scanner.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.ArtistName,
		&user.Role,
		&user.AvatarURL,
		&user.BannerURL,
		&user.Bio,
		&user.Telegram,
		&user.Vk,
		&user.Youtube,
		&user.Website,
		&user.ArtistRequested,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		user.Email, user.PasswordHash, user.FullName, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, artist_name, role,
		       avatar_url, banner_url, bio,
		       telegram, vk, youtube, website, artist_requested,
		       created_at, updated_at
		FROM users WHERE email = $1
	`
	user, err := scanUser(r.db.QueryRow(ctx, query, email))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, artist_name, role,
		       avatar_url, banner_url, bio,
		       telegram, vk, youtube, website, artist_requested,
		       created_at, updated_at
		FROM users WHERE id = $1
	`
	user, err := scanUser(r.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET
			full_name = COALESCE($2, full_name),
			artist_name = COALESCE($3, artist_name),
			avatar_url = COALESCE($4, avatar_url),
			banner_url = COALESCE($5, banner_url),
			bio = COALESCE($6, bio),
			telegram = COALESCE($7, telegram),
			vk = COALESCE($8, vk),
			youtube = COALESCE($9, youtube),
			website = COALESCE($10, website),
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.FullName,
		user.ArtistName,
		user.AvatarURL,
		user.BannerURL,
		user.Bio,
		user.Telegram,
		user.Vk,
		user.Youtube,
		user.Website,
	)
	return err
}

func (r *UserRepository) SetArtistRequested(ctx context.Context, userID int, requested bool) error {
	query := `UPDATE users SET artist_requested = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, requested, userID)
	return err
}

func (r *UserRepository) UpdateRole(ctx context.Context, userID int, role string) error {
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, role, userID)
	return err
}
