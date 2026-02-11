package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/avvvet/brightpath/internal/config"
	"github.com/avvvet/brightpath/internal/store"
)

func main() {
	log.Println("Starting database migrations...")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := store.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	// Get migration files
	migrationsDir := "migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Sort files to ensure order
	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	if len(migrationFiles) == 0 {
		log.Println("No migration files found")
		return
	}

	log.Printf("Found %d migration files", len(migrationFiles))

	// Execute each migration
	for _, filename := range migrationFiles {
		log.Printf("Running migration: %s", filename)

		sqlPath := filepath.Join(migrationsDir, filename)
		sqlBytes, err := os.ReadFile(sqlPath)
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", filename, err)
		}

		sql := string(sqlBytes)

		// Execute migration
		_, err = db.Conn().Exec(sql)
		if err != nil {
			log.Fatalf("Failed to execute migration %s: %v", filename, err)
		}

		log.Printf("✓ Completed: %s", filename)
	}

	log.Println("All migrations completed successfully!")
}
