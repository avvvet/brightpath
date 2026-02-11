package domain

import (
	"time"

	"github.com/google/uuid"
)

type Lead struct {
	ID         uuid.UUID `db:"id"`
	PropertyID uuid.UUID `db:"property_id"`

	// Status flow: new → qualified → sms_sent → replied → hot → forwarded → contracted → closed → dead
	Status string `db:"status"`

	// Contact (from skip trace)
	Phone          *string `db:"phone"`
	PhoneType      *string `db:"phone_type"`
	PhoneConnected *bool   `db:"phone_connected"`
	PhoneDNC       *bool   `db:"phone_dnc"`
	AltPhones      *string `db:"alt_phones"` // JSONB stored as string
	Email          *string `db:"email"`
	AltEmails      *string `db:"alt_emails"` // JSONB stored as string

	// Demographics (from skip trace)
	OwnerAge    *int    `db:"owner_age"`
	OwnerGender *string `db:"owner_gender"`

	// Scoring (Phase 2)
	LeadScore     *int    `db:"lead_score"`
	ScoreReasons  *string `db:"score_reasons"` // TEXT[] stored as string, parse as needed
	DaysToAuction *int    `db:"days_to_auction"`

	// Deal tracking (Phase 3)
	CloserID    *uuid.UUID `db:"closer_id"`
	ForwardedAt *time.Time `db:"forwarded_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
