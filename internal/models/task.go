package models

import (
	"time"
)

// WorkType represents the type of work shift
type WorkType string

const (
	WorkTypeShiftLead WorkType = "shift_lead"
	WorkTypeInbound   WorkType = "inbound"
	WorkTypeOutbound  WorkType = "outbound"
)

// OutboundSubtype represents the subtype for outbound call shifts
type OutboundSubtype string

const (
	SubtypeExtraShift   OutboundSubtype = "extra"
	SubtypeRegularShift OutboundSubtype = "regular"
)

// InboundRuleHours defines working hours for inbound rules
type InboundRuleHours struct {
	Workday float64
	Holiday float64
}

var InboundRules = map[string]InboundRuleHours{
	"101": {Workday: 13.0, Holiday: 13.0}, // Example values, adjust if needed
	"102": {Workday: 13.0, Holiday: 13.0},
	"103": {Workday: 10.0, Holiday: 10.0},
	"104": {Workday: 6.5, Holiday: 6.5},
	"105": {Workday: 10.0, Holiday: 10.0},
	"106": {Workday: 13.0, Holiday: 11.0},
	"107": {Workday: 10.0, Holiday: 0},
}

const HoursPerShift = 6.5

// TrackItem represents a work tracking entry in the system
type TrackItem struct {
	ID            int             `json:"id"`
	UserID        int             `json:"user_id"`
	Type          WorkType        `json:"type"`
	Subtype       OutboundSubtype `json:"subtype,omitempty"`
	InboundRule   string          `json:"inbound_rule,omitempty"`
	EmergencyCall bool            `json:"emergency_call"`
	HolidayCall   bool            `json:"holiday_call"`
	WorkingHours  float64         `json:"working_hours"`
	WorkingShifts float64         `json:"working_shifts"`
	Date          time.Time       `json:"date"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// CreateTrackItemRequest represents the data needed to create a new track item
type CreateTrackItemRequest struct {
	Type          WorkType        `json:"type"`
	Subtype       OutboundSubtype `json:"subtype,omitempty"`
	InboundRule   string          `json:"inbound_rule,omitempty"`
	EmergencyCall bool            `json:"emergency_call"`
	HolidayCall   bool            `json:"holiday_call"`
	WorkingHours  float64         `json:"working_hours"`
	WorkingShifts float64         `json:"working_shifts"`
	Date          string          `json:"date"` // ISO 8601 format: "2024-01-20T10:00:00Z"
}

// UpdateTrackItemRequest represents the data needed to update a track item
type UpdateTrackItemRequest struct {
	Type          *WorkType        `json:"type,omitempty"`
	Subtype       *OutboundSubtype `json:"subtype,omitempty"`
	InboundRule   *string          `json:"inbound_rule,omitempty"`
	EmergencyCall *bool            `json:"emergency_call,omitempty"`
	HolidayCall   *bool            `json:"holiday_call,omitempty"`
	WorkingHours  *float64         `json:"working_hours,omitempty"`
	WorkingShifts *float64         `json:"working_shifts,omitempty"`
	Date          *string          `json:"date,omitempty"` // ISO 8601 format
}

// DateRangeQuery represents a query for track items within a date range
type DateRangeQuery struct {
	StartDate string `json:"start_date"` // ISO 8601 format: "2024-01-20T00:00:00Z"
	EndDate   string `json:"end_date"`   // ISO 8601 format: "2024-01-25T23:59:59Z"
}
