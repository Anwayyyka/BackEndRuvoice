// internal/service/user_service.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/Anwayyyka/ruvoice-backend/internal/pkg/hash"
	"github.com/Anwayyyka/ruvoice-backend/internal/pkg/jwt"
	"github.com/Anwayyyka/ruvoice-backend/internal/repository"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

type RegisterInput struct {
	Email    string
	Password string
	FullName string
}

func (s *UserService) Register(ctx context.Context, input RegisterInput) (*domain.User, string, error) {
	existing, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", errors.New("Пользователь с таким email уже существует")
	}

	hashed, err := hash.HashPassword(input.Password)
	if err != nil {
		return nil, "", err
	}

	user := &domain.User{
		Email:        input.Email,
		PasswordHash: hashed,
		FullName:     &input.FullName,
		Role:         "user",
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, "", err
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret, 7*24*time.Hour)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *UserService) Login(ctx context.Context, input LoginInput) (*domain.User, string, error) {
	user, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return nil, "", errors.New("Неверный email или пароль")
	}
	if !hash.CheckPasswordHash(input.Password, user.PasswordHash) {
		return nil, "", errors.New("Неверный email или пароль")
	}
	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret, 7*24*time.Hour)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *UserService) GetProfile(ctx context.Context, userID int) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("Пользователь не найден")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID int, updates map[string]interface{}) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("Пользователь не найден")
	}

	if name, ok := updates["full_name"].(string); ok {
		user.FullName = &name
	}
	if bio, ok := updates["bio"].(string); ok {
		user.Bio = &bio
	}
	if avatar, ok := updates["avatar_url"].(string); ok {
		user.AvatarURL = &avatar
	}
	if banner, ok := updates["banner_url"].(string); ok {
		user.BannerURL = &banner
	}
	if telegram, ok := updates["telegram"].(string); ok {
		user.Telegram = &telegram
	}
	if vk, ok := updates["vk"].(string); ok {
		user.Vk = &vk
	}
	if youtube, ok := updates["youtube"].(string); ok {
		user.Youtube = &youtube
	}
	if website, ok := updates["website"].(string); ok {
		user.Website = &website
	}

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) RequestArtist(ctx context.Context, userID int, artistName, bio string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("Пользователь не найден")
	}
	user.ArtistName = &artistName
	user.Bio = &bio
	user.Role = "artist"

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
