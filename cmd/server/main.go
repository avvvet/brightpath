package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/avvvet/brightpath/internal/config"
	"github.com/avvvet/brightpath/internal/pipeline"
	"github.com/avvvet/brightpath/internal/providers/reapi"
	"github.com/avvvet/brightpath/internal/store"
)

type Server struct {
	cfg         *config.Config
	db          *store.DB
	reapiClient *reapi.Client
}

func main() {
	log.Println("Starting BrightPath API server...")

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

	// Create server
	server := &Server{
		cfg:         cfg,
		db:          db,
		reapiClient: reapiClient,
	}

	// Setup routes
	http.HandleFunc("/health", server.handleHealth)
	http.HandleFunc("/api/pipeline/run", server.handlePipelineRun)
	http.HandleFunc("/api/pipeline/runs", server.handleGetPipelineRuns)
	http.HandleFunc("/api/leads", server.handleGetLeads)

	// Start server
	port := ":8080"
	log.Printf("Server listening on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// handlePipelineRun triggers a pipeline run for a specific state
// POST /api/pipeline/run?state=TX
func (s *Server) handlePipelineRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("state")))
	if state == "" {
		http.Error(w, "state parameter is required", http.StatusBadRequest)
		return
	}

	log.Printf("API: Triggering pipeline run for state: %s", state)

	// Run pipeline in a goroutine (async)
	go func() {
		p := pipeline.New(s.cfg, s.db, s.reapiClient, state)
		if err := p.Run(); err != nil {
			log.Printf("Pipeline run failed for %s: %v", state, err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "started",
		"state":   state,
		"message": "Pipeline run started in background",
	})
}

// handleGetPipelineRuns returns recent pipeline runs
// GET /api/pipeline/runs?state=TX&limit=10
func (s *Server) handleGetPipelineRuns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query params
	stateParam := r.URL.Query().Get("state")
	var state *string
	if stateParam != "" {
		stateUpper := strings.ToUpper(strings.TrimSpace(stateParam))
		state = &stateUpper
	}

	limitParam := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Query database
	runs, err := s.db.GetRecentPipelineRuns(state, limit)
	if err != nil {
		log.Printf("Failed to get pipeline runs: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(runs)
}

// handleGetLeads returns leads by status
// GET /api/leads?status=new&limit=50
func (s *Server) handleGetLeads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query params
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "new" // default
	}

	limitParam := r.URL.Query().Get("limit")
	limit := 50 // default
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 500 {
			limit = l
		}
	}

	// Query database
	leads, err := s.db.GetLeadsByStatus(status, limit)
	if err != nil {
		log.Printf("Failed to get leads: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leads)
}
