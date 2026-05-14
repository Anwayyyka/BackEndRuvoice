package service

import (
	"context"
	"errors"
	"strings"

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

func (s *TrackService) GetTracksByIDs(ctx context.Context, ids []int) ([]*domain.Track, error) {
	return s.trackRepo.GetByIDs(ctx, ids)
}

func (s *TrackService) GetOrCreateExternalTrack(ctx context.Context, source string, externalID int, trackData *domain.Track) (*domain.Track, error) {
	existing, err := s.trackRepo.GetByExternal(ctx, source, externalID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	trackData.ExternalID = &externalID
	trackData.ExternalSource = &source
	trackData.IsExternal = true
	if err := s.trackRepo.Create(ctx, trackData); err != nil {
		return nil, err
	}
	return trackData, nil
}

func (s *TrackService) ensureApprovedTrack(ctx context.Context, trackID int) error {
	track, err := s.trackRepo.GetByID(ctx, trackID)
	if err != nil {
		return err
	}
	if track == nil {
		return errors.New("track not found")
	}
	return nil
}

func (s *TrackService) IncrementPlays(ctx context.Context, trackID int) error {
	if err := s.ensureApprovedTrack(ctx, trackID); err != nil {
		return err
	}
	return s.trackRepo.IncrementPlays(ctx, trackID)
}

func (s *TrackService) Like(ctx context.Context, userID, trackID int) error {
	if err := s.ensureApprovedTrack(ctx, trackID); err != nil {
		return err
	}
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
	if err := s.likeRepo.Delete(ctx, userID, trackID); err != nil {
		return err
	}
	return s.trackRepo.UpdateLikesCount(ctx, trackID, -1)
}

func (s *TrackService) AddFavorite(ctx context.Context, userID, trackID int) error {
	if err := s.ensureApprovedTrack(ctx, trackID); err != nil {
		return err
	}
	fav, err := s.favRepo.Create(ctx, userID, trackID)
	if err != nil {
		return err
	}
	if fav == nil {
		return nil // уже в избранном
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
	return s.trackRepo.UpdateStatus(ctx, trackID, "approved")
}

func (s *TrackService) Reject(ctx context.Context, trackID int, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return errors.New("reason is required")
	}
	return s.trackRepo.UpdateStatus(ctx, trackID, "rejected")
}
