package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sergey/work-track-backend/internal/models"
	"github.com/sergey/work-track-backend/internal/repository"
)

// TrackItemService handles track item business logic
type TrackItemService struct {
	trackItemRepo *repository.TrackItemRepository
}

// NewTrackItemService creates a new track item service
func NewTrackItemService(trackItemRepo *repository.TrackItemRepository) *TrackItemService {
	return &TrackItemService{
		trackItemRepo: trackItemRepo,
	}
}

// CreateTrackItem creates a new track item for a user
func (s *TrackItemService) CreateTrackItem(ctx context.Context, userID int, req *models.CreateTrackItemRequest) (*models.TrackItem, error) {
	// Validate input
	if req.Type == "" {
		return nil, errors.New("type is required")
	}

	// Parse date
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, use ISO 8601 (RFC3339): %w", err)
	}

	item := &models.TrackItem{
		UserID:        userID,
		Type:          req.Type,
		EmergencyCall: req.EmergencyCall,
		HolidayCall:   req.HolidayCall,
		WorkingHours:  req.WorkingHours,
		WorkingShifts: req.WorkingShifts,
		Date:          date,
	}

	err = s.trackItemRepo.Create(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to create track item: %w", err)
	}

	return item, nil
}

// GetUserTrackItems retrieves all track items for a user
func (s *TrackItemService) GetUserTrackItems(ctx context.Context, userID int) ([]models.TrackItem, error) {
	items, err := s.trackItemRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get track items: %w", err)
	}

	return items, nil
}

// GetTrackItemsByDateRange retrieves track items for a user within a date range
func (s *TrackItemService) GetTrackItemsByDateRange(ctx context.Context, userID int, startDateStr, endDateStr string) ([]models.TrackItem, error) {
	// Parse dates (accept date-only format)
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format, use YYYY-MM-DD: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format, use YYYY-MM-DD: %w", err)
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	items, err := s.trackItemRepo.FindByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get track items by date range: %w", err)
	}

	return items, nil
}

// GetTrackItem retrieves a specific track item, ensuring it belongs to the user
func (s *TrackItemService) GetTrackItem(ctx context.Context, userID, itemID int) (*models.TrackItem, error) {
	item, err := s.trackItemRepo.FindByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if item.UserID != userID {
		return nil, ErrUnauthorized
	}

	return item, nil
}

// UpdateTrackItem updates a track item, ensuring it belongs to the user
func (s *TrackItemService) UpdateTrackItem(ctx context.Context, userID, itemID int, req *models.UpdateTrackItemRequest) (*models.TrackItem, error) {
	// Get existing item
	item, err := s.trackItemRepo.FindByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if item.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Update fields if provided
	if req.Type != nil {
		item.Type = *req.Type
	}
	if req.EmergencyCall != nil {
		item.EmergencyCall = *req.EmergencyCall
	}
	if req.HolidayCall != nil {
		item.HolidayCall = *req.HolidayCall
	}
	if req.WorkingHours != nil {
		item.WorkingHours = *req.WorkingHours
	}
	if req.WorkingShifts != nil {
		item.WorkingShifts = *req.WorkingShifts
	}
	if req.Date != nil {
		date, err := time.Parse(time.RFC3339, *req.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, use ISO 8601 (RFC3339): %w", err)
		}
		item.Date = date
	}

	err = s.trackItemRepo.Update(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to update track item: %w", err)
	}

	return item, nil
}

// DeleteTrackItem deletes a track item, ensuring it belongs to the user
func (s *TrackItemService) DeleteTrackItem(ctx context.Context, userID, itemID int) error {
	// Get existing item
	item, err := s.trackItemRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}

	// Verify ownership
	if item.UserID != userID {
		return ErrUnauthorized
	}

	err = s.trackItemRepo.Delete(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to delete track item: %w", err)
	}

	return nil
}
