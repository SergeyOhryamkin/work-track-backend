package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sergey/work-track-backend/internal/database"
	"github.com/sergey/work-track-backend/internal/models"
)

type trackItemRow struct {
	ID            int
	Type          models.WorkType
	Subtype       sql.NullString
	InboundRule   sql.NullString
	Emergency     bool
	Holiday       bool
	WorkingHours  float64
	WorkingShifts float64
	Date          time.Time
}

func inboundHoursToShifts(hours float64) float64 {
	switch hours {
	case 13.0:
		return 2.0
	case 11.0:
		return 1.69
	case 10.0:
		return 1.5
	case 6.5:
		return 1.0
	default:
		return hours / 6.5
	}
}

func recalcItem(item trackItemRow) (float64, float64, error) {
	switch item.Type {
	case models.WorkTypeShiftLead:
		hours := 11.0
		return hours, hours / models.HoursPerShiftDefault, nil
	case models.WorkTypeInbound:
		ruleKey := item.InboundRule.String
		rule, ok := models.InboundRules[ruleKey]
		if !ok {
			return 0, 0, fmt.Errorf("invalid inbound rule: %s", ruleKey)
		}
		var hours float64
		if item.Holiday {
			hours = rule.Holiday
		} else {
			hours = rule.Workday
		}
		return hours, inboundHoursToShifts(hours), nil
	case models.WorkTypeOutbound:
		hours := item.WorkingHours
		if hours <= 0 {
			return 0, 0, fmt.Errorf("outbound working hours must be > 0")
		}
		return hours, hours / models.HoursPerShiftDefault, nil
	default:
		return 0, 0, fmt.Errorf("unsupported work type: %s", item.Type)
	}
}

func dbPath() string {
	if value := os.Getenv("DB_PATH"); value != "" {
		return value
	}
	return "./worktrack.db"
}

func main() {
	_ = godotenv.Load()

	dryRun := flag.Bool("dry-run", true, "log changes without updating rows")
	flag.Parse()

	db, err := database.NewSQLiteDB(dbPath())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	rows, err := db.Query(`
		SELECT id, type, subtype, inbound_rule, emergency_call, holiday_call, working_hours, working_shifts, date
		FROM track_items
		ORDER BY id
	`)
	if err != nil {
		log.Fatalf("Failed to query track items: %v", err)
	}
	defer rows.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}

	updateStmt, err := tx.Prepare(`
		UPDATE track_items
		SET working_hours = ?, working_shifts = ?, updated_at = datetime('now')
		WHERE id = ?
	`)
	if err != nil {
		_ = tx.Rollback()
		log.Fatalf("Failed to prepare update statement: %v", err)
	}
	defer updateStmt.Close()

	var updated int
	for rows.Next() {
		var item trackItemRow
		if err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Subtype,
			&item.InboundRule,
			&item.Emergency,
			&item.Holiday,
			&item.WorkingHours,
			&item.WorkingShifts,
			&item.Date,
		); err != nil {
			_ = tx.Rollback()
			log.Fatalf("Failed to scan track item: %v", err)
		}

		newHours, newShifts, err := recalcItem(item)
		if err != nil {
			_ = tx.Rollback()
			log.Fatalf("Failed to recalc item %d: %v", item.ID, err)
		}

		if newHours == item.WorkingHours && newShifts == item.WorkingShifts {
			continue
		}

		log.Printf("item %d: hours %.2f -> %.2f, shifts %.2f -> %.2f", item.ID, item.WorkingHours, newHours, item.WorkingShifts, newShifts)
		updated++

		if *dryRun {
			continue
		}

		if _, err := updateStmt.Exec(newHours, newShifts, item.ID); err != nil {
			_ = tx.Rollback()
			log.Fatalf("Failed to update item %d: %v", item.ID, err)
		}
	}

	if err := rows.Err(); err != nil {
		_ = tx.Rollback()
		log.Fatalf("Row iteration error: %v", err)
	}

	if *dryRun {
		_ = tx.Rollback()
		log.Printf("Dry run complete. %d items would be updated.", updated)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Printf("Manual task complete. %d items updated.", updated)
}
