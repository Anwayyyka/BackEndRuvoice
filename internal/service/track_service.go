package service

import (
	"context"
	"errors"

	"github.com/Anwayyyka/ruvoice-backend/internal/domain"
	"github.com/Anwayyyka/ruvoice-backend/internal/repository"
)

type TrackService struct {
	trackRepo *repository.TrackRepository
	userRepo  *repository.UserRepository
	likeRepo  *repository.LikeRepository
	favRepo   *repository.FavoriteRepository
}

func NewTrackService(
	trackRepo *repository.TrackRepository,
	userRepo *repository.UserRepository,
	likeRepo *repository.LikeRepository,
	favRepo *repository.FavoriteRepository,
) *TrackService {
	return &TrackService{
		trackRepo: trackRepo,
		userRepo:  userRepo,
		likeRepo:  likeRepo,
		favRepo:   favRepo,
	}
}

func (s *TrackService) Create(ctx context.Context, track *domain.Track) error {
	return s.trackRepo.Create(ctx, track)
}

func (s *TrackService) GetApproved(ctx context.Context) ([]*domain.Track, error) {
	return s.trackRepo.ListApproved(ctx)
}

func (s *TrackService) GetByArtist(ctx context.Context, artistID int) ([]*domain.Track, error) {
	return s.trackRepo.ListByArtist(ctx, artistID)
}

func (s *TrackService) GetByID(ctx context.Context, id int) (*domain.Track, error) {
	return s.trackRepo.GetByID(ctx, id)
}

func (s *TrackService) IncrementPlays(ctx context.Context, trackID int) error {
	return s.trackRepo.IncrementPlays(ctx, trackID)
}

func (s *TrackService) Like(ctx context.Context, userID, trackID int) error {
	like, err := s.likeRepo.Create(ctx, userID, trackID)
	if err != nil {
		return err
	}
	if like != nil {
		return s.trackRepo.UpdateLikesCount(ctx, trackID, 1)
	}
	return nil
}

func (s *TrackService) Unlike(ctx context.Context, userID, trackID int) error {
	like, err := s.likeRepo.GetByUserAndTrack(ctx, userID, trackID)
	if err != nil {
		return err
	}
	if like == nil {
		return errors.New("like not found")
	}
	err = s.likeRepo.Delete(ctx, userID, trackID)
	if err != nil {
		return err
	}
	return s.trackRepo.UpdateLikesCount(ctx, trackID, -1)
}

func (s *TrackService) AddFavorite(ctx context.Context, userID, trackID int) error {
	fav, err := s.favRepo.Create(ctx, userID, trackID)
	if err != nil {
		return err
	}
	if fav == nil {
		return errors.New("already in favorites")
	}
	return nil
}

func (s *TrackService) RemoveFavorite(ctx context.Context, userID, trackID int) error {
	return s.favRepo.Delete(ctx, userID, trackID)
}

func (s *TrackService) GetUserFavorites(ctx context.Context, userID int) ([]*domain.Favorite, error) {
	return s.favRepo.GetByUser(ctx, userID)
}

func (s *TrackService) GetPending(ctx context.Context) ([]*domain.Track, error) {
	return s.trackRepo.ListPending(ctx)
}

func (s *TrackService) Approve(ctx context.Context, trackID int) error {
	return s.trackRepo.Approve(ctx, trackID)
}

func (s *TrackService) Reject(ctx context.Context, trackID int, reason string) error {
	return s.trackRepo.Reject(ctx, trackID, reason)
}
