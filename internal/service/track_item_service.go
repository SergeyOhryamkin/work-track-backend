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

	// Parse date first as we need it for holiday check and validation
	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, use ISO 8601 (RFC3339): %w", err)
	}

	// Work type specific validation and defaults
	var workingHours = req.WorkingHours
	var workingShifts = req.WorkingShifts

	switch req.Type {
	case models.WorkTypeShiftLead:
		// Non-call shift: no required inputs, lasts 8 hours
		workingHours = 8.0
		workingShifts = workingHours / models.HoursPerShift
	case models.WorkTypeInbound:
		// Inbound: requires inbound rule
		if req.InboundRule == "" {
			return nil, errors.New("inbound rule is required for inbound shifts")
		}

		rule, ok := models.InboundRules[req.InboundRule]
		if !ok {
			return nil, fmt.Errorf("invalid inbound rule: %s", req.InboundRule)
		}

		// Determine hours based on holiday status
		if req.HolidayCall {
			workingHours = rule.Holiday
		} else {
			workingHours = rule.Workday
		}
		workingShifts = workingHours / models.HoursPerShift

	case models.WorkTypeOutbound:
		// Outbound: requires hours and subtype
		if workingHours <= 0 {
			return nil, errors.New("working hours are required for outbound shifts")
		}
		if req.Subtype == "" {
			return nil, errors.New("subtype (regular/extra) is required for outbound shifts")
		}
		workingShifts = workingHours / models.HoursPerShift
	default:
		return nil, fmt.Errorf("invalid work type: %s", req.Type)
	}

	item := &models.TrackItem{
		UserID:        userID,
		Type:          req.Type,
		Subtype:       req.Subtype,
		InboundRule:   req.InboundRule,
		EmergencyCall: req.EmergencyCall,
		HolidayCall:   req.HolidayCall,
		WorkingHours:  workingHours,
		WorkingShifts: workingShifts,
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

// GetTrackItemsSummaryByDateRange aggregates shift totals for a user within a date range
func (s *TrackItemService) GetTrackItemsSummaryByDateRange(ctx context.Context, userID int, startDateStr, endDateStr string) (*models.TrackItemSummary, error) {
	items, err := s.GetTrackItemsByDateRange(ctx, userID, startDateStr, endDateStr)
	if err != nil {
		return nil, err
	}

	summary := &models.TrackItemSummary{}
	for _, item := range items {
		switch item.Type {
		case models.WorkTypeShiftLead:
			summary.ShiftLeadShifts += item.WorkingShifts
		case models.WorkTypeInbound:
			summary.InboundShifts += item.WorkingShifts
		case models.WorkTypeOutbound:
			summary.OutboundShifts += item.WorkingShifts
		}

		if item.EmergencyCall {
			summary.EmergencyCallShifts += item.WorkingShifts
		}
		if item.HolidayCall {
			summary.HolidayCallShifts += item.WorkingShifts
		}
	}

	summary.TotalShifts = summary.ShiftLeadShifts + summary.InboundShifts + summary.OutboundShifts
	return summary, nil
}

// GetTrackItemsByDateRange retrieves track items for a user within a date range
func (s *TrackItemService) GetTrackItemsByDateRange(ctx context.Context, userID int, startDateStr, endDateStr string) ([]models.TrackItem, error) {
	// Parse dates using ISO 8601 (RFC3339)
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format, use ISO 8601 (RFC3339): %w", err)
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format, use ISO 8601 (RFC3339): %w", err)
	}

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
	if req.Subtype != nil {
		item.Subtype = *req.Subtype
	}
	if req.InboundRule != nil {
		item.InboundRule = *req.InboundRule
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

	// Recalculate derived fields based on updated Type and rules
	switch item.Type {
	case models.WorkTypeShiftLead:
		item.WorkingHours = 8.0
		item.WorkingShifts = item.WorkingHours / models.HoursPerShift
	case models.WorkTypeInbound:
		if rule, ok := models.InboundRules[item.InboundRule]; ok {
			if item.HolidayCall {
				item.WorkingHours = rule.Holiday
			} else {
				item.WorkingHours = rule.Workday
			}
		}
		item.WorkingShifts = item.WorkingHours / models.HoursPerShift
	case models.WorkTypeOutbound:
		if req.WorkingHours != nil {
			item.WorkingHours = *req.WorkingHours
		}
		item.WorkingShifts = item.WorkingHours / models.HoursPerShift
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
