package domain

import (
	"time"

	"github.com/google/uuid"
)

type PipelineRun struct {
	ID          uuid.UUID  `db:"id"`
	State       string     `db:"state"`
	StartedAt   time.Time  `db:"started_at"`
	CompletedAt *time.Time `db:"completed_at"`
	Status      string     `db:"status"` // running, completed, failed

	// Stats
	IDsFound             int `db:"ids_found"`
	IDsNew               int `db:"ids_new"`
	PropertiesEnriched   int `db:"properties_enriched"`
	PropertiesQualified  int `db:"properties_qualified"`
	LeadsCreated         int `db:"leads_created"`
	SkipTraceHits        int `db:"skip_trace_hits"`
	CreditsUsedProperty  int `db:"credits_used_property"`
	CreditsUsedSkiptrace int `db:"credits_used_skiptrace"`

	Error  *string `db:"error"`
	Config *string `db:"config"` // JSONB stored as string
}
