package domain

import (
	"time"

	"github.com/google/uuid"
)

type Property struct {
	ID       uuid.UUID `db:"id"`
	Source   string    `db:"source"`
	SourceID string    `db:"source_id"`

	// Address (from propertyInfo.address)
	Address  string  `db:"address"`
	Street   *string `db:"street"`
	City     *string `db:"city"`
	State    string  `db:"state"`
	Zip      *string `db:"zip"`
	County   *string `db:"county"`
	FIPS     *string `db:"fips"`
	Location *string `db:"location"` // PostGIS GEOGRAPHY stored as string

	// Property details (from propertyInfo)
	PropertyType *string  `db:"property_type"`
	Bedrooms     *int     `db:"bedrooms"`
	Bathrooms    *float64 `db:"bathrooms"`
	Sqft         *int     `db:"sqft"`
	YearBuilt    *int     `db:"year_built"`
	LotSqft      *int     `db:"lot_sqft"`

	// Valuation
	AVM           *int `db:"avm"`
	AssessedValue *int `db:"assessed_value"`

	// Mortgage & equity (from root + currentMortgages)
	OpenMortgageBalance *int    `db:"open_mortgage_balance"`
	EquityAmount        *int    `db:"equity_amount"`
	EquityPercent       *int    `db:"equity_percent"`
	LoanType            *string `db:"loan_type"`
	LenderName          *string `db:"lender_name"`

	// Foreclosure (from root + auctionInfo)
	PreForeclosure  *bool      `db:"pre_foreclosure"`
	Auction         *bool      `db:"auction"`
	AuctionDate     *time.Time `db:"auction_date"`
	AuctionTime     *string    `db:"auction_time"`
	AuctionLocation *string    `db:"auction_location"`
	NoticeType      *string    `db:"notice_type"`
	RecordingDate   *time.Time `db:"recording_date"`
	Foreclosure     *bool      `db:"foreclosure"`

	// Auction enrichment (from auctionInfo)
	BankEstimatedValue      *int    `db:"bank_estimated_value"`
	TrusteeName             *string `db:"trustee_name"`
	TrusteePhone            *string `db:"trustee_phone"`
	ForeclosureDocumentType *string `db:"foreclosure_document_type"`

	// Owner (from ownerInfo)
	Owner1FirstName       *string `db:"owner1_first_name"`
	Owner1LastName        *string `db:"owner1_last_name"`
	Owner1FullName        *string `db:"owner1_full_name"`
	Owner1Type            *string `db:"owner1_type"`
	Owner2FirstName       *string `db:"owner2_first_name"`
	Owner2LastName        *string `db:"owner2_last_name"`
	OwnerOccupied         *bool   `db:"owner_occupied"`
	AbsenteeOwner         *bool   `db:"absentee_owner"`
	InvestorBuyer         *bool   `db:"investor_buyer"`
	CorporateOwned        *bool   `db:"corporate_owned"`
	OwnershipLengthMonths *int    `db:"ownership_length_months"`
	MailStreet            *string `db:"mail_street"`
	MailCity              *string `db:"mail_city"`
	MailState             *string `db:"mail_state"`
	MailZip               *string `db:"mail_zip"`

	// Extras for future LLM scoring
	SuggestedRent     *string    `db:"suggested_rent"`
	MedianIncome      *string    `db:"median_income"`
	FloodZone         *bool      `db:"flood_zone"`
	HOA               *bool      `db:"hoa"`
	TaxLien           *bool      `db:"tax_lien"`
	TaxAmount         *string    `db:"tax_amount"`
	TaxDelinquentYear *string    `db:"tax_delinquent_year"`
	YearsOwned        *int       `db:"years_owned"`
	LastSaleDate      *time.Time `db:"last_sale_date"`

	// Skip trace tracking
	SkipTraceStatus      *string    `db:"skip_trace_status"`
	SkipTraceAttemptedAt *time.Time `db:"skip_trace_attempted_at"`
	SkipTraceError       *string    `db:"skip_trace_error"`

	// Metadata
	RawData      *string    `db:"raw_data"` // JSONB stored as string, parse as needed
	DiscoveredAt time.Time  `db:"discovered_at"`
	EnrichedAt   *time.Time `db:"enriched_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}
