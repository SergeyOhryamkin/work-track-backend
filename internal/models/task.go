package models

import (
	"time"
)

// TrackItem represents a work tracking entry in the system
type TrackItem struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Type          string    `json:"type"`
	EmergencyCall bool      `json:"emergency_call"`
	HolidayCall   bool      `json:"holiday_call"`
	WorkingHours  float64   `json:"working_hours"`
	WorkingShifts float64   `json:"working_shifts"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateTrackItemRequest represents the data needed to create a new track item
type CreateTrackItemRequest struct {
	Type          string  `json:"type"`
	EmergencyCall bool    `json:"emergency_call"`
	HolidayCall   bool    `json:"holiday_call"`
	WorkingHours  float64 `json:"working_hours"`
	WorkingShifts float64 `json:"working_shifts"`
	Date          string  `json:"date"` // ISO 8601 format: "2024-01-20T10:00:00Z"
}

// UpdateTrackItemRequest represents the data needed to update a track item
type UpdateTrackItemRequest struct {
	Type          *string  `json:"type,omitempty"`
	EmergencyCall *bool    `json:"emergency_call,omitempty"`
	HolidayCall   *bool    `json:"holiday_call,omitempty"`
	WorkingHours  *float64 `json:"working_hours,omitempty"`
	WorkingShifts *float64 `json:"working_shifts,omitempty"`
	Date          *string  `json:"date,omitempty"` // ISO 8601 format
}

// DateRangeQuery represents a query for track items within a date range
type DateRangeQuery struct {
	StartDate string `json:"start_date"` // ISO 8601 format: "2024-01-20"
	EndDate   string `json:"end_date"`   // ISO 8601 format: "2024-01-25"
}
