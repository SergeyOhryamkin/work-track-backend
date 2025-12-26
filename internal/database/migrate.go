package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// RunMigrations applies all .up.sql migration files in the migrations directory
func RunMigrations(db *sql.DB, migrationsPath string) error {
	log.Println("Running database migrations...")

	// Read migration files
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort .up.sql files
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" && 
		   len(file.Name()) > 7 && file.Name()[len(file.Name())-7:] == ".up.sql" {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	if len(migrationFiles) == 0 {
		log.Println("No migration files found")
		return nil
	}

	// Apply each migration
	for _, filename := range migrationFiles {
		log.Printf("Applying migration: %s", filename)
		
		filePath := filepath.Join(migrationsPath, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Execute the migration SQL
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		log.Printf("✓ Migration applied: %s", filename)
	}

	log.Println("All migrations completed successfully")
	return nil
}
