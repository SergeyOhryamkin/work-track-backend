package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sergey/work-track-backend/internal/models"
)

var (
	ErrTrackItemNotFound = errors.New("track item not found")
)

// TrackItemRepository handles database operations for track items
type TrackItemRepository struct {
	db *pgxpool.Pool
}

// NewTrackItemRepository creates a new track item repository
func NewTrackItemRepository(db *pgxpool.Pool) *TrackItemRepository {
	return &TrackItemRepository{db: db}
}

// Create inserts a new track item into the database
func (r *TrackItemRepository) Create(ctx context.Context, item *models.TrackItem) error {
	query := `
		INSERT INTO track_items (user_id, type, emergency_call, holiday_call, working_hours, working_shifts, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, item.UserID, item.Type, item.EmergencyCall, item.HolidayCall, item.WorkingHours, item.WorkingShifts, item.Date).
		Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create track item: %w", err)
	}

	return nil
}

// FindByUserID retrieves all track items for a specific user
func (r *TrackItemRepository) FindByUserID(ctx context.Context, userID int) ([]models.TrackItem, error) {
	query := `
		SELECT id, user_id, type, emergency_call, holiday_call, working_hours, working_shifts, date, created_at, updated_at
		FROM track_items
		WHERE user_id = $1
		ORDER BY date DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query track items: %w", err)
	}
	defer rows.Close()

	var items []models.TrackItem
	for rows.Next() {
		var item models.TrackItem
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Type,
			&item.EmergencyCall,
			&item.HolidayCall,
			&item.WorkingHours,
			&item.WorkingShifts,
			&item.Date,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating track items: %w", err)
	}

	return items, nil
}

// FindByDateRange retrieves track items for a user within a date range
func (r *TrackItemRepository) FindByDateRange(ctx context.Context, userID int, startDate, endDate time.Time) ([]models.TrackItem, error) {
	query := `
		SELECT id, user_id, type, emergency_call, holiday_call, working_hours, working_shifts, date, created_at, updated_at
		FROM track_items
		WHERE user_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date DESC
	`

	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query track items by date range: %w", err)
	}
	defer rows.Close()

	var items []models.TrackItem
	for rows.Next() {
		var item models.TrackItem
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Type,
			&item.EmergencyCall,
			&item.HolidayCall,
			&item.WorkingHours,
			&item.WorkingShifts,
			&item.Date,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating track items: %w", err)
	}

	return items, nil
}

// FindByID retrieves a specific track item by ID
func (r *TrackItemRepository) FindByID(ctx context.Context, id int) (*models.TrackItem, error) {
	query := `
		SELECT id, user_id, type, emergency_call, holiday_call, working_hours, working_shifts, date, created_at, updated_at
		FROM track_items
		WHERE id = $1
	`

	var item models.TrackItem
	err := r.db.QueryRow(ctx, query, id).
		Scan(&item.ID, &item.UserID, &item.Type, &item.EmergencyCall, &item.HolidayCall, &item.WorkingHours, &item.WorkingShifts, &item.Date, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTrackItemNotFound
		}
		return nil, fmt.Errorf("failed to find track item: %w", err)
	}

	return &item, nil
}

// Update updates an existing track item
func (r *TrackItemRepository) Update(ctx context.Context, item *models.TrackItem) error {
	query := `
		UPDATE track_items
		SET type = $1, emergency_call = $2, holiday_call = $3, working_hours = $4, working_shifts = $5, date = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query, item.Type, item.EmergencyCall, item.HolidayCall, item.WorkingHours, item.WorkingShifts, item.Date, item.ID).
		Scan(&item.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTrackItemNotFound
		}
		return fmt.Errorf("failed to update track item: %w", err)
	}

	return nil
}

// Delete removes a track item from the database
func (r *TrackItemRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM track_items WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete track item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrTrackItemNotFound
	}

	return nil
}
