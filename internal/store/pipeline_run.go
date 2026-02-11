package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/avvvet/brightpath/internal/domain"
	"github.com/google/uuid"
)

// CreatePipelineRun creates a new pipeline run record
func (db *DB) CreatePipelineRun(state string, config string) (*domain.PipelineRun, error) {
	run := &domain.PipelineRun{
		State:  state,
		Status: "running",
		Config: &config,
	}

	query := `
		INSERT INTO pipeline_runs (state, status, config)
		VALUES ($1, $2, $3)
		RETURNING id, started_at
	`

	err := db.conn.QueryRow(query, run.State, run.Status, run.Config).Scan(&run.ID, &run.StartedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create pipeline run: %w", err)
	}

	return run, nil
}

// UpdatePipelineRun updates the stats of a running pipeline
func (db *DB) UpdatePipelineRun(run *domain.PipelineRun) error {
	query := `
		UPDATE pipeline_runs SET
			ids_found = $1,
			ids_new = $2,
			properties_enriched = $3,
			properties_qualified = $4,
			leads_created = $5,
			skip_trace_hits = $6,
			credits_used_property = $7,
			credits_used_skiptrace = $8,
			updated_at = NOW()
		WHERE id = $9
	`

	_, err := db.conn.Exec(
		query,
		run.IDsFound,
		run.IDsNew,
		run.PropertiesEnriched,
		run.PropertiesQualified,
		run.LeadsCreated,
		run.SkipTraceHits,
		run.CreditsUsedProperty,
		run.CreditsUsedSkiptrace,
		run.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update pipeline run: %w", err)
	}

	return nil
}

// CompletePipelineRun marks a pipeline run as completed or failed
func (db *DB) CompletePipelineRun(id uuid.UUID, status string, errMsg *string) error {
	query := `
		UPDATE pipeline_runs SET
			status = $1,
			completed_at = NOW(),
			error = $2
		WHERE id = $3
	`

	_, err := db.conn.Exec(query, status, errMsg, id)
	if err != nil {
		return fmt.Errorf("failed to complete pipeline run: %w", err)
	}

	return nil
}

// GetPipelineRunByID retrieves a pipeline run by ID
func (db *DB) GetPipelineRunByID(id uuid.UUID) (*domain.PipelineRun, error) {
	var run domain.PipelineRun
	query := `
		SELECT id, state, started_at, completed_at, status,
		       ids_found, ids_new, properties_enriched, properties_qualified,
		       leads_created, skip_trace_hits, credits_used_property, credits_used_skiptrace,
		       error, config
		FROM pipeline_runs
		WHERE id = $1
	`

	err := db.conn.QueryRow(query, id).Scan(
		&run.ID, &run.State, &run.StartedAt, &run.CompletedAt, &run.Status,
		&run.IDsFound, &run.IDsNew, &run.PropertiesEnriched, &run.PropertiesQualified,
		&run.LeadsCreated, &run.SkipTraceHits, &run.CreditsUsedProperty, &run.CreditsUsedSkiptrace,
		&run.Error, &run.Config,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline run: %w", err)
	}

	return &run, nil
}

// GetRecentPipelineRuns retrieves recent pipeline runs, optionally filtered by state
func (db *DB) GetRecentPipelineRuns(state *string, limit int) ([]*domain.PipelineRun, error) {
	var query string
	var args []interface{}

	if state != nil {
		query = `
			SELECT id, state, started_at, completed_at, status,
			       ids_found, ids_new, properties_enriched, properties_qualified,
			       leads_created, skip_trace_hits, credits_used_property, credits_used_skiptrace,
			       error, config
			FROM pipeline_runs
			WHERE state = $1
			ORDER BY started_at DESC
			LIMIT $2
		`
		args = []interface{}{*state, limit}
	} else {
		query = `
			SELECT id, state, started_at, completed_at, status,
			       ids_found, ids_new, properties_enriched, properties_qualified,
			       leads_created, skip_trace_hits, credits_used_property, credits_used_skiptrace,
			       error, config
			FROM pipeline_runs
			ORDER BY started_at DESC
			LIMIT $1
		`
		args = []interface{}{limit}
	}

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query pipeline runs: %w", err)
	}
	defer rows.Close()

	var runs []*domain.PipelineRun
	for rows.Next() {
		var run domain.PipelineRun
		err := rows.Scan(
			&run.ID, &run.State, &run.StartedAt, &run.CompletedAt, &run.Status,
			&run.IDsFound, &run.IDsNew, &run.PropertiesEnriched, &run.PropertiesQualified,
			&run.LeadsCreated, &run.SkipTraceHits, &run.CreditsUsedProperty, &run.CreditsUsedSkiptrace,
			&run.Error, &run.Config,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pipeline run: %w", err)
		}
		runs = append(runs, &run)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pipeline runs: %w", err)
	}

	return runs, nil
}

// GetPipelineRunStats returns aggregate stats for a time period
func (db *DB) GetPipelineRunStats(state *string, since time.Time) (*domain.PipelineRun, error) {
	var query string
	var args []interface{}

	if state != nil {
		query = `
			SELECT 
				COALESCE(SUM(ids_found), 0) as ids_found,
				COALESCE(SUM(ids_new), 0) as ids_new,
				COALESCE(SUM(properties_enriched), 0) as properties_enriched,
				COALESCE(SUM(properties_qualified), 0) as properties_qualified,
				COALESCE(SUM(leads_created), 0) as leads_created,
				COALESCE(SUM(skip_trace_hits), 0) as skip_trace_hits,
				COALESCE(SUM(credits_used_property), 0) as credits_used_property,
				COALESCE(SUM(credits_used_skiptrace), 0) as credits_used_skiptrace
			FROM pipeline_runs
			WHERE state = $1 AND started_at >= $2 AND status = 'completed'
		`
		args = []interface{}{*state, since}
	} else {
		query = `
			SELECT 
				COALESCE(SUM(ids_found), 0) as ids_found,
				COALESCE(SUM(ids_new), 0) as ids_new,
				COALESCE(SUM(properties_enriched), 0) as properties_enriched,
				COALESCE(SUM(properties_qualified), 0) as properties_qualified,
				COALESCE(SUM(leads_created), 0) as leads_created,
				COALESCE(SUM(skip_trace_hits), 0) as skip_trace_hits,
				COALESCE(SUM(credits_used_property), 0) as credits_used_property,
				COALESCE(SUM(credits_used_skiptrace), 0) as credits_used_skiptrace
			FROM pipeline_runs
			WHERE started_at >= $1 AND status = 'completed'
		`
		args = []interface{}{since}
	}

	stats := &domain.PipelineRun{}
	err := db.conn.QueryRow(query, args...).Scan(
		&stats.IDsFound,
		&stats.IDsNew,
		&stats.PropertiesEnriched,
		&stats.PropertiesQualified,
		&stats.LeadsCreated,
		&stats.SkipTraceHits,
		&stats.CreditsUsedProperty,
		&stats.CreditsUsedSkiptrace,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline stats: %w", err)
	}

	return stats, nil
}
