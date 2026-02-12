package pipeline

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/avvvet/brightpath/internal/config"
	"github.com/avvvet/brightpath/internal/domain"
	"github.com/avvvet/brightpath/internal/providers/reapi"
	"github.com/avvvet/brightpath/internal/store"
)

type Pipeline struct {
	cfg   *config.Config
	store *store.DB
	reapi *reapi.Client
	state string
}

// New creates a new pipeline instance for a specific state
func New(cfg *config.Config, db *store.DB, reapiClient *reapi.Client, state string) *Pipeline {
	return &Pipeline{
		cfg:   cfg,
		store: db,
		reapi: reapiClient,
		state: state,
	}
}

// Run executes the complete pipeline for the configured state
func (p *Pipeline) Run() error {
	log.Printf("[%s] Starting pipeline run", p.state)

	// Create pipeline run record
	run, err := p.store.CreatePipelineRun(p.state, fmt.Sprintf(`{"auction_days_min": %d, "auction_days_max": %d, "equity_min": %d}`,
		p.cfg.AuctionDaysMin, p.cfg.AuctionDaysMax, p.cfg.EquityMin))
	if err != nil {
		return fmt.Errorf("failed to create pipeline run: %w", err)
	}

	// Defer completion handling
	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("panic: %v", r)
			_ = p.store.CompletePipelineRun(run.ID, "failed", &errMsg)
			panic(r)
		}
	}()

	// Step 1: Discovery (FREE - ids_only search)
	log.Printf("[%s] Step 1: Discovery - searching for property IDs", p.state)

	auctionDateMin := time.Now().AddDate(0, 0, p.cfg.AuctionDaysMin)
	auctionDateMax := time.Now().AddDate(0, 0, p.cfg.AuctionDaysMax)

	propertyIDs, err := p.reapi.SearchPropertyIDs(p.state, auctionDateMin, auctionDateMax, p.cfg.EquityMin)
	if err != nil {
		errMsg := fmt.Sprintf("property search failed: %v", err)
		_ = p.store.CompletePipelineRun(run.ID, "failed", &errMsg)
		return fmt.Errorf("property search failed: %w", err)
	}

	run.IDsFound = len(propertyIDs)
	log.Printf("[%s] Found %d property IDs", p.state, run.IDsFound)

	// Step 2: Enrich NEW properties AND load existing qualified properties
	log.Printf("[%s] Step 2: Enriching new properties and loading existing qualified properties", p.state)

	qualifiedProperties, err := p.enrichAndLoadProperties(propertyIDs, run)
	if err != nil {
		errMsg := fmt.Sprintf("enrichment failed: %v", err)
		_ = p.store.CompletePipelineRun(run.ID, "failed", &errMsg)
		return fmt.Errorf("enrichment failed: %w", err)
	}

	log.Printf("[%s] Enriched %d new properties, %d total qualified for skip tracing", p.state, run.PropertiesEnriched, len(qualifiedProperties))

	// Step 3: Skip trace qualified leads (1 credit each)
	log.Printf("[%s] Step 3: Skip tracing qualified leads", p.state)

	err = p.skipTraceQualifiedLeads(qualifiedProperties, run)
	if err != nil {
		errMsg := fmt.Sprintf("skip trace failed: %v", err)
		_ = p.store.CompletePipelineRun(run.ID, "failed", &errMsg)
		return fmt.Errorf("skip trace failed: %w", err)
	}

	log.Printf("[%s] Created %d leads (%d with phones)", p.state, run.LeadsCreated, run.SkipTraceHits)

	// Update final stats and complete
	if err := p.store.UpdatePipelineRun(run); err != nil {
		log.Printf("[%s] Warning: failed to update pipeline run stats: %v", p.state, err)
	}

	if err := p.store.CompletePipelineRun(run.ID, "completed", nil); err != nil {
		log.Printf("[%s] Warning: failed to mark pipeline run complete: %v", p.state, err)
	}

	log.Printf("[%s] Pipeline run completed successfully", p.state)
	log.Printf("[%s] Stats: %d IDs found, %d new, %d enriched, %d qualified, %d leads, %d with phones",
		p.state, run.IDsFound, run.IDsNew, run.PropertiesEnriched, run.PropertiesQualified, run.LeadsCreated, run.SkipTraceHits)
	log.Printf("[%s] Credits used: %d property, %d skip trace, %d total",
		p.state, run.CreditsUsedProperty, run.CreditsUsedSkiptrace, run.CreditsUsedProperty+run.CreditsUsedSkiptrace)

	return nil
}

// enrichAndLoadProperties fetches PropertyDetail for new IDs AND loads existing qualified properties
// This ensures we can resume skip tracing for properties that were enriched but not skip traced
func (p *Pipeline) enrichAndLoadProperties(propertyIDs []int64, run *domain.PipelineRun) ([]*domain.Property, error) {
	var qualifiedProperties []*domain.Property

	for _, id := range propertyIDs {
		sourceID := strconv.FormatInt(id, 10)

		// CREDIT GUARD: Check if property already exists
		exists, err := p.store.PropertyExistsBySourceID("reapi", sourceID)
		if err != nil {
			return nil, fmt.Errorf("failed to check property existence for ID %d: %w", id, err)
		}

		var property *domain.Property

		if exists {
			// Property already in DB - LOAD IT instead of skipping
			property, err = p.store.GetPropertyBySourceID("reapi", sourceID)
			if err != nil {
				log.Printf("[%s] Warning: failed to load existing property ID %d: %v", p.state, id, err)
				continue
			}

			// Don't count as "new" or "enriched" (already done)
			// But DO check if it qualifies for skip tracing
		} else {
			// Property is NEW - fetch detail (costs 1 credit)
			run.IDsNew++

			detail, err := p.reapi.GetPropertyDetail(id)
			if err != nil {
				log.Printf("[%s] Warning: failed to get property detail for ID %d: %v", p.state, id, err)
				continue
			}

			run.CreditsUsedProperty++
			run.PropertiesEnriched++

			// Map to domain property
			property, err = MapPropertyDetailToProperty(detail)
			if err != nil {
				log.Printf("[%s] Warning: failed to map property ID %d: %v", p.state, id, err)
				continue
			}

			// Save new property to DB
			if err := p.store.UpsertProperty(property); err != nil {
				log.Printf("[%s] Warning: failed to save property ID %d: %v", p.state, id, err)
				continue
			}
		}

		// Apply Go post-filters (no API cost) - applies to BOTH new and existing
		if !PassesAllFilters(property) {
			// Property doesn't qualify, skip
			continue
		}

		// Property passed all filters - add to qualified list
		run.PropertiesQualified++
		qualifiedProperties = append(qualifiedProperties, property)

		// Update run stats periodically (every 10 properties)
		if run.PropertiesQualified%10 == 0 {
			if err := p.store.UpdatePipelineRun(run); err != nil {
				log.Printf("[%s] Warning: failed to update pipeline run: %v", p.state, err)
			}
		}
	}

	return qualifiedProperties, nil
}

// skipTraceQualifiedLeads performs skip tracing on qualified properties (credit guard)
func (p *Pipeline) skipTraceQualifiedLeads(properties []*domain.Property, run *domain.PipelineRun) error {
	for _, prop := range properties {
		// CREDIT GUARD: Check if lead already exists for this property
		exists, err := p.store.LeadExistsByPropertyID(prop.ID)
		if err != nil {
			return fmt.Errorf("failed to check lead existence for property %s: %w", prop.ID, err)
		}

		if exists {
			// Already skip traced, skip (saves 1 credit)
			continue
		}

		// Validate required fields for skip trace
		if prop.Owner1FirstName == nil || prop.Owner1LastName == nil ||
			prop.MailStreet == nil || prop.MailCity == nil ||
			prop.MailState == nil || prop.MailZip == nil {
			log.Printf("[%s] Warning: property %s missing required skip trace fields", p.state, prop.ID)
			continue
		}

		// Perform skip trace (costs 1 credit)
		skipTraceReq := reapi.SkipTraceRequest{
			FirstName:   *prop.Owner1FirstName,
			LastName:    *prop.Owner1LastName,
			MailAddress: *prop.MailStreet,
			MailCity:    *prop.MailCity,
			MailState:   *prop.MailState,
			MailZip:     *prop.MailZip,
			MatchRequirements: reapi.MatchRequirements{
				Phones: true,
			},
		}

		skipTraceResp, err := p.reapi.SkipTrace(skipTraceReq)
		if err != nil {
			log.Printf("[%s] Warning: skip trace failed for property %s: %v", p.state, prop.ID, err)
			continue
		}

		run.CreditsUsedSkiptrace++

		// Calculate days to auction
		daysToAuction := 0
		if prop.AuctionDate != nil {
			daysToAuction = CalculateDaysToAuction(*prop.AuctionDate)
		}

		// Map to lead (returns nil if no valid mobile phone found)
		lead, err := MapSkipTraceToLead(prop.ID, *prop.Owner1FirstName, *prop.Owner1LastName, skipTraceResp, daysToAuction)
		if err != nil {
			log.Printf("[%s] Warning: failed to map skip trace for property %s: %v", p.state, prop.ID, err)
			continue
		}

		if lead == nil {
			// No valid phone found, skip lead creation
			continue
		}

		// Create lead
		if err := p.store.CreateLead(lead); err != nil {
			log.Printf("[%s] Warning: failed to create lead for property %s: %v", p.state, prop.ID, err)
			continue
		}

		run.LeadsCreated++
		run.SkipTraceHits++

		// Update run stats periodically (every 5 leads)
		if run.LeadsCreated%5 == 0 {
			if err := p.store.UpdatePipelineRun(run); err != nil {
				log.Printf("[%s] Warning: failed to update pipeline run: %v", p.state, err)
			}
		}
	}

	return nil
}
