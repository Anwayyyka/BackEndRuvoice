package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/Anwayyyka/ruvoice-backend/internal/pkg/hash"
	"github.com/Anwayyyka/ruvoice-backend/internal/pkg/jwt"
	"github.com/Anwayyyka/ruvoice-backend/internal/repository"
)

type UserService struct {
	repo          *repository.UserRepository
	artistReqRepo *repository.ArtistRequestRepository
	jwtSecret     string
}

func NewUserService(repo *repository.UserRepository, artistReqRepo *repository.ArtistRequestRepository, jwtSecret string) *UserService {
	return &UserService{
		repo:          repo,
		artistReqRepo: artistReqRepo,
		jwtSecret:     jwtSecret,
	}
}

type RegisterInput struct {
	Email    string
	Password string
	FullName string
}

func (s *UserService) Register(ctx context.Context, input RegisterInput) (*domain.User, string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" || strings.TrimSpace(input.Password) == "" {
		return nil, "", errors.New("email и пароль обязательны")
	}

	existing, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", errors.New("Пользователь с таким email уже существует")
	}

	hashed, err := hash.HashPassword(input.Password)
	if err != nil {
		return nil, "", err
	}

	var fullName *string
	if name := strings.TrimSpace(input.FullName); name != "" {
		fullName = &name
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: hashed,
		FullName:     fullName,
		Role:         "user",
	}

	if err := s.repo.Create(ctx, user); err != nil {
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
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" || strings.TrimSpace(input.Password) == "" {
		return nil, "", errors.New("Неверный email или пароль")
	}

	user, err := s.repo.GetByEmail(ctx, email)
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

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) RequestArtist(ctx context.Context, userID int, artistName, bio string) (*domain.User, error) {
	artistName = strings.TrimSpace(artistName)
	bio = strings.TrimSpace(bio)

	if artistName == "" {
		return nil, errors.New("artist_name обязателен")
	}

	if s.artistReqRepo == nil {
		return nil, errors.New("artist request repository not configured")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("Пользователь не найден")
	}

	if user.Role == "artist" || user.Role == "admin" {
		return nil, errors.New("Пользователь уже является артистом")
	}

	existingReq, err := s.artistReqRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingReq != nil {
		switch existingReq.Status {
		case "pending":
			return nil, errors.New("Заявка уже на модерации")
		case "approved":
			return nil, errors.New("Заявка уже одобрена")
		}
	}

	req := &domain.ArtistRequest{
		UserID:     userID,
		ArtistName: artistName,
		Bio:        bio,
		Status:     "pending",
	}
	if err := s.artistReqRepo.Create(ctx, req); err != nil {
		return nil, err
	}

	user.ArtistName = &artistName
	user.Bio = &bio
	user.ArtistRequested = true

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	if err := s.repo.SetArtistRequested(ctx, userID, true); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.repo.GetByEmail(ctx, email)
}
