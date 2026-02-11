package store

import (
	"database/sql"
	"fmt"

	"github.com/avvvet/brightpath/internal/domain"
	"github.com/google/uuid"
)

// LeadExistsByPropertyID checks if a lead exists for a given property
// This is the CREDIT GUARD for skip tracing - never skip trace the same property twice
func (db *DB) LeadExistsByPropertyID(propertyID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM leads WHERE property_id = $1)`
	err := db.conn.QueryRow(query, propertyID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check lead existence: %w", err)
	}
	return exists, nil
}

// GetLeadByPropertyID retrieves a lead by property_id
func (db *DB) GetLeadByPropertyID(propertyID uuid.UUID) (*domain.Lead, error) {
	var lead domain.Lead
	query := `
		SELECT id, property_id, status,
		       phone, phone_type, phone_connected, phone_dnc, alt_phones,
		       email, alt_emails,
		       owner_age, owner_gender,
		       lead_score, score_reasons, days_to_auction,
		       closer_id, forwarded_at,
		       created_at, updated_at
		FROM leads
		WHERE property_id = $1
	`
	err := db.conn.QueryRow(query, propertyID).Scan(
		&lead.ID, &lead.PropertyID, &lead.Status,
		&lead.Phone, &lead.PhoneType, &lead.PhoneConnected, &lead.PhoneDNC, &lead.AltPhones,
		&lead.Email, &lead.AltEmails,
		&lead.OwnerAge, &lead.OwnerGender,
		&lead.LeadScore, &lead.ScoreReasons, &lead.DaysToAuction,
		&lead.CloserID, &lead.ForwardedAt,
		&lead.CreatedAt, &lead.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get lead: %w", err)
	}
	return &lead, nil
}

// GetLeadByID retrieves a lead by UUID
func (db *DB) GetLeadByID(id uuid.UUID) (*domain.Lead, error) {
	var lead domain.Lead
	query := `
		SELECT id, property_id, status,
		       phone, phone_type, phone_connected, phone_dnc, alt_phones,
		       email, alt_emails,
		       owner_age, owner_gender,
		       lead_score, score_reasons, days_to_auction,
		       closer_id, forwarded_at,
		       created_at, updated_at
		FROM leads
		WHERE id = $1
	`
	err := db.conn.QueryRow(query, id).Scan(
		&lead.ID, &lead.PropertyID, &lead.Status,
		&lead.Phone, &lead.PhoneType, &lead.PhoneConnected, &lead.PhoneDNC, &lead.AltPhones,
		&lead.Email, &lead.AltEmails,
		&lead.OwnerAge, &lead.OwnerGender,
		&lead.LeadScore, &lead.ScoreReasons, &lead.DaysToAuction,
		&lead.CloserID, &lead.ForwardedAt,
		&lead.CreatedAt, &lead.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get lead: %w", err)
	}
	return &lead, nil
}

// CreateLead inserts a new lead
func (db *DB) CreateLead(lead *domain.Lead) error {
	query := `
		INSERT INTO leads (
			property_id, status,
			phone, phone_type, phone_connected, phone_dnc, alt_phones,
			email, alt_emails,
			owner_age, owner_gender,
			lead_score, score_reasons, days_to_auction,
			closer_id, forwarded_at
		) VALUES (
			$1, $2,
			$3, $4, $5, $6, $7,
			$8, $9,
			$10, $11,
			$12, $13, $14,
			$15, $16
		)
		RETURNING id, created_at, updated_at
	`

	err := db.conn.QueryRow(
		query,
		lead.PropertyID, lead.Status,
		lead.Phone, lead.PhoneType, lead.PhoneConnected, lead.PhoneDNC, lead.AltPhones,
		lead.Email, lead.AltEmails,
		lead.OwnerAge, lead.OwnerGender,
		lead.LeadScore, lead.ScoreReasons, lead.DaysToAuction,
		lead.CloserID, lead.ForwardedAt,
	).Scan(&lead.ID, &lead.CreatedAt, &lead.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create lead: %w", err)
	}

	return nil
}

// UpdateLeadStatus updates the status of a lead
func (db *DB) UpdateLeadStatus(id uuid.UUID, status string) error {
	query := `UPDATE leads SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := db.conn.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update lead status: %w", err)
	}
	return nil
}

// GetLeadsByStatus retrieves leads by status
func (db *DB) GetLeadsByStatus(status string, limit int) ([]*domain.Lead, error) {
	query := `
		SELECT id, property_id, status,
		       phone, phone_type, phone_connected, phone_dnc, alt_phones,
		       email, alt_emails,
		       owner_age, owner_gender,
		       lead_score, score_reasons, days_to_auction,
		       closer_id, forwarded_at,
		       created_at, updated_at
		FROM leads
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := db.conn.Query(query, status, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query leads: %w", err)
	}
	defer rows.Close()

	var leads []*domain.Lead
	for rows.Next() {
		var lead domain.Lead
		err := rows.Scan(
			&lead.ID, &lead.PropertyID, &lead.Status,
			&lead.Phone, &lead.PhoneType, &lead.PhoneConnected, &lead.PhoneDNC, &lead.AltPhones,
			&lead.Email, &lead.AltEmails,
			&lead.OwnerAge, &lead.OwnerGender,
			&lead.LeadScore, &lead.ScoreReasons, &lead.DaysToAuction,
			&lead.CloserID, &lead.ForwardedAt,
			&lead.CreatedAt, &lead.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lead: %w", err)
		}
		leads = append(leads, &lead)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating leads: %w", err)
	}

	return leads, nil
}
