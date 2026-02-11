package main

import (
	"flag"
	"log"
	"strings"

	"github.com/avvvet/brightpath/internal/config"
	"github.com/avvvet/brightpath/internal/pipeline"
	"github.com/avvvet/brightpath/internal/providers/reapi"
	"github.com/avvvet/brightpath/internal/store"
)

func main() {
	// Parse CLI flags
	stateFlag := flag.String("state", "", "State to process (e.g., TX, FL, GA)")
	flag.Parse()

	if *stateFlag == "" {
		log.Fatal("--state flag is required (e.g., --state=TX)")
	}

	state := strings.ToUpper(strings.TrimSpace(*stateFlag))
	log.Printf("Worker starting for state: %s", state)

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

	// Initialize REAPI client
	reapiClient := reapi.NewClient(cfg.REAPIKey)
	log.Println("REAPI client initialized")

	// Create pipeline for this state
	p := pipeline.New(cfg, db, reapiClient, state)

	// Run pipeline
	log.Printf("Starting pipeline run for %s", state)
	if err := p.Run(); err != nil {
		log.Fatalf("Pipeline run failed: %v", err)
	}

	log.Printf("Pipeline run completed successfully for %s", state)
}
