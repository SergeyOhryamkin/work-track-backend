package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sergey/work-track-backend/internal/models"
)

var (
	ErrTrackItemNotFound = errors.New("track item not found")
)

// TrackItemRepository handles database operations for track items
type TrackItemRepository struct {
	db *sql.DB
}

// NewTrackItemRepository creates a new track item repository
func NewTrackItemRepository(db *sql.DB) *TrackItemRepository {
	return &TrackItemRepository{db: db}
}

// Create inserts a new track item into the database
func (r *TrackItemRepository) Create(ctx context.Context, item *models.TrackItem) error {
	query := `
		INSERT INTO track_items (user_id, type, emergency_call, holiday_call, working_hours, working_shifts, date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`

	result, err := r.db.ExecContext(ctx, query, item.UserID, item.Type, item.EmergencyCall, item.HolidayCall, item.WorkingHours, item.WorkingShifts, item.Date)
	if err != nil {
		return fmt.Errorf("failed to create track item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	item.ID = int(id)

	err = r.db.QueryRowContext(ctx, "SELECT created_at, updated_at FROM track_items WHERE id = ?", item.ID).
		Scan(&item.CreatedAt, &item.UpdatedAt)

	return err
}

// FindByUserID retrieves all track items for a specific user
func (r *TrackItemRepository) FindByUserID(ctx context.Context, userID int) ([]models.TrackItem, error) {
	query := `
		SELECT id, user_id, type, emergency_call, holiday_call, working_hours, working_shifts, date, created_at, updated_at
		FROM track_items
		WHERE user_id = ?
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
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
		WHERE user_id = ? AND date >= ? AND date <= ?
		ORDER BY date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, startDate, endDate)
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
		WHERE id = ?
	`

	var item models.TrackItem
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&item.ID, &item.UserID, &item.Type, &item.EmergencyCall, &item.HolidayCall, &item.WorkingHours, &item.WorkingShifts, &item.Date, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
		SET type = ?, emergency_call = ?, holiday_call = ?, working_hours = ?, working_shifts = ?, date = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, item.Type, item.EmergencyCall, item.HolidayCall, item.WorkingHours, item.WorkingShifts, item.Date, item.ID)
	if err != nil {
		return fmt.Errorf("failed to update track item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrTrackItemNotFound
	}

	err = r.db.QueryRowContext(ctx, "SELECT updated_at FROM track_items WHERE id = ?", item.ID).
		Scan(&item.UpdatedAt)

	return nil
}

// Delete removes a track item from the database
func (r *TrackItemRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM track_items WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete track item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return ErrTrackItemNotFound
	}

	return nil
}
