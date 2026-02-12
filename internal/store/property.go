package store

import (
	"database/sql"
	"fmt"

	"github.com/avvvet/brightpath/internal/domain"
	"github.com/google/uuid"
)

// PropertyExistsBySourceID checks if a property exists by source and source_id
// This is the CREDIT GUARD - must be called before expensive API calls
func (db *DB) PropertyExistsBySourceID(source, sourceID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM properties WHERE source = $1 AND source_id = $2)`
	err := db.conn.QueryRow(query, source, sourceID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check property existence: %w", err)
	}
	return exists, nil
}

// GetPropertyBySourceID retrieves a property by source and source_id
func (db *DB) GetPropertyBySourceID(source, sourceID string) (*domain.Property, error) {
	var prop domain.Property
	query := `
		SELECT id, source, source_id, address, street, city, state, zip, county, fips,
			ST_AsText(location) as location,
			property_type, bedrooms, bathrooms, sqft, year_built, lot_sqft,
			avm, assessed_value,
			open_mortgage_balance, equity_amount, equity_percent, loan_type, lender_name,
			pre_foreclosure, auction, auction_date, auction_time, auction_location,
			notice_type, recording_date, foreclosure,
			bank_estimated_value, trustee_name, trustee_phone, foreclosure_document_type,
			owner1_first_name, owner1_last_name, owner1_full_name, owner1_type,
			owner2_first_name, owner2_last_name,
			owner_occupied, absentee_owner, investor_buyer, corporate_owned,
			ownership_length_months, mail_street, mail_city, mail_state, mail_zip,
			suggested_rent, median_income, flood_zone, hoa, tax_lien,
			tax_amount, tax_delinquent_year, years_owned, last_sale_date,
			skip_trace_status, skip_trace_attempted_at, skip_trace_error,
			raw_data, discovered_at, enriched_at, created_at, updated_at
		FROM properties
		WHERE source = $1 AND source_id = $2
	`

	err := db.conn.QueryRow(query, source, sourceID).Scan(
		&prop.ID, &prop.Source, &prop.SourceID, &prop.Address, &prop.Street, &prop.City,
		&prop.State, &prop.Zip, &prop.County, &prop.FIPS, &prop.Location,
		&prop.PropertyType, &prop.Bedrooms, &prop.Bathrooms, &prop.Sqft, &prop.YearBuilt, &prop.LotSqft,
		&prop.AVM, &prop.AssessedValue,
		&prop.OpenMortgageBalance, &prop.EquityAmount, &prop.EquityPercent, &prop.LoanType, &prop.LenderName,
		&prop.PreForeclosure, &prop.Auction, &prop.AuctionDate, &prop.AuctionTime, &prop.AuctionLocation,
		&prop.NoticeType, &prop.RecordingDate, &prop.Foreclosure,
		&prop.BankEstimatedValue, &prop.TrusteeName, &prop.TrusteePhone, &prop.ForeclosureDocumentType,
		&prop.Owner1FirstName, &prop.Owner1LastName, &prop.Owner1FullName, &prop.Owner1Type,
		&prop.Owner2FirstName, &prop.Owner2LastName,
		&prop.OwnerOccupied, &prop.AbsenteeOwner, &prop.InvestorBuyer, &prop.CorporateOwned,
		&prop.OwnershipLengthMonths, &prop.MailStreet, &prop.MailCity, &prop.MailState, &prop.MailZip,
		&prop.SuggestedRent, &prop.MedianIncome, &prop.FloodZone, &prop.HOA, &prop.TaxLien,
		&prop.TaxAmount, &prop.TaxDelinquentYear, &prop.YearsOwned, &prop.LastSaleDate,
		&prop.SkipTraceStatus, &prop.SkipTraceAttemptedAt, &prop.SkipTraceError,
		&prop.RawData, &prop.DiscoveredAt, &prop.EnrichedAt, &prop.CreatedAt, &prop.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get property: %w", err)
	}
	return &prop, nil
}

// GetPropertyByID retrieves a property by UUID
func (db *DB) GetPropertyByID(id uuid.UUID) (*domain.Property, error) {
	var prop domain.Property
	query := `
		SELECT id, source, source_id, address, street, city, state, zip, county, fips,
			ST_AsText(location) as location,
			property_type, bedrooms, bathrooms, sqft, year_built, lot_sqft,
			avm, assessed_value,
			open_mortgage_balance, equity_amount, equity_percent, loan_type, lender_name,
			pre_foreclosure, auction, auction_date, auction_time, auction_location,
			notice_type, recording_date, foreclosure,
			bank_estimated_value, trustee_name, trustee_phone, foreclosure_document_type,
			owner1_first_name, owner1_last_name, owner1_full_name, owner1_type,
			owner2_first_name, owner2_last_name,
			owner_occupied, absentee_owner, investor_buyer, corporate_owned,
			ownership_length_months, mail_street, mail_city, mail_state, mail_zip,
			suggested_rent, median_income, flood_zone, hoa, tax_lien,
			tax_amount, tax_delinquent_year, years_owned, last_sale_date,
			skip_trace_status, skip_trace_attempted_at, skip_trace_error,
			raw_data, discovered_at, enriched_at, created_at, updated_at
		FROM properties
		WHERE id = $1
	`
	err := db.conn.QueryRow(query, id).Scan(
		&prop.ID, &prop.Source, &prop.SourceID, &prop.Address, &prop.Street, &prop.City,
		&prop.State, &prop.Zip, &prop.County, &prop.FIPS, &prop.Location,
		&prop.PropertyType, &prop.Bedrooms, &prop.Bathrooms, &prop.Sqft, &prop.YearBuilt, &prop.LotSqft,
		&prop.AVM, &prop.AssessedValue,
		&prop.OpenMortgageBalance, &prop.EquityAmount, &prop.EquityPercent, &prop.LoanType, &prop.LenderName,
		&prop.PreForeclosure, &prop.Auction, &prop.AuctionDate, &prop.AuctionTime, &prop.AuctionLocation,
		&prop.NoticeType, &prop.RecordingDate, &prop.Foreclosure,
		&prop.BankEstimatedValue, &prop.TrusteeName, &prop.TrusteePhone, &prop.ForeclosureDocumentType,
		&prop.Owner1FirstName, &prop.Owner1LastName, &prop.Owner1FullName, &prop.Owner1Type,
		&prop.Owner2FirstName, &prop.Owner2LastName,
		&prop.OwnerOccupied, &prop.AbsenteeOwner, &prop.InvestorBuyer, &prop.CorporateOwned,
		&prop.OwnershipLengthMonths, &prop.MailStreet, &prop.MailCity, &prop.MailState, &prop.MailZip,
		&prop.SuggestedRent, &prop.MedianIncome, &prop.FloodZone, &prop.HOA, &prop.TaxLien,
		&prop.TaxAmount, &prop.TaxDelinquentYear, &prop.YearsOwned, &prop.LastSaleDate,
		&prop.SkipTraceStatus, &prop.SkipTraceAttemptedAt, &prop.SkipTraceError,
		&prop.RawData, &prop.DiscoveredAt, &prop.EnrichedAt, &prop.CreatedAt, &prop.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get property: %w", err)
	}
	return &prop, nil
}

// UpsertProperty inserts or updates a property (ON CONFLICT UPDATE)
func (db *DB) UpsertProperty(prop *domain.Property) error {
	query := `
		INSERT INTO properties (
			source, source_id, address, street, city, state, zip, county, fips,
			location,
			property_type, bedrooms, bathrooms, sqft, year_built, lot_sqft,
			avm, assessed_value,
			open_mortgage_balance, equity_amount, equity_percent, loan_type, lender_name,
			pre_foreclosure, auction, auction_date, auction_time, auction_location,
			notice_type, recording_date, foreclosure,
			bank_estimated_value, trustee_name, trustee_phone, foreclosure_document_type,
			owner1_first_name, owner1_last_name, owner1_full_name, owner1_type,
			owner2_first_name, owner2_last_name,
			owner_occupied, absentee_owner, investor_buyer, corporate_owned,
			ownership_length_months, mail_street, mail_city, mail_state, mail_zip,
			suggested_rent, median_income, flood_zone, hoa, tax_lien,
			tax_amount, tax_delinquent_year, years_owned, last_sale_date,
			raw_data, discovered_at, enriched_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			ST_GeogFromText($10),
			$11, $12, $13, $14, $15, $16,
			$17, $18,
			$19, $20, $21, $22, $23,
			$24, $25, $26, $27, $28,
			$29, $30, $31,
			$32, $33, $34, $35,
			$36, $37, $38, $39,
			$40, $41,
			$42, $43, $44, $45,
			$46, $47, $48, $49, $50,
			$51, $52, $53, $54, $55,
			$56, $57, $58, $59,
			$60, $61, $62
		)
		ON CONFLICT (source, source_id) DO UPDATE SET
			address = EXCLUDED.address,
			street = EXCLUDED.street,
			city = EXCLUDED.city,
			state = EXCLUDED.state,
			zip = EXCLUDED.zip,
			county = EXCLUDED.county,
			fips = EXCLUDED.fips,
			location = EXCLUDED.location,
			property_type = EXCLUDED.property_type,
			bedrooms = EXCLUDED.bedrooms,
			bathrooms = EXCLUDED.bathrooms,
			sqft = EXCLUDED.sqft,
			year_built = EXCLUDED.year_built,
			lot_sqft = EXCLUDED.lot_sqft,
			avm = EXCLUDED.avm,
			assessed_value = EXCLUDED.assessed_value,
			open_mortgage_balance = EXCLUDED.open_mortgage_balance,
			equity_amount = EXCLUDED.equity_amount,
			equity_percent = EXCLUDED.equity_percent,
			loan_type = EXCLUDED.loan_type,
			lender_name = EXCLUDED.lender_name,
			pre_foreclosure = EXCLUDED.pre_foreclosure,
			auction = EXCLUDED.auction,
			auction_date = EXCLUDED.auction_date,
			auction_time = EXCLUDED.auction_time,
			auction_location = EXCLUDED.auction_location,
			notice_type = EXCLUDED.notice_type,
			recording_date = EXCLUDED.recording_date,
			foreclosure = EXCLUDED.foreclosure,
			bank_estimated_value = EXCLUDED.bank_estimated_value,
			trustee_name = EXCLUDED.trustee_name,
			trustee_phone = EXCLUDED.trustee_phone,
			foreclosure_document_type = EXCLUDED.foreclosure_document_type,
			owner1_first_name = EXCLUDED.owner1_first_name,
			owner1_last_name = EXCLUDED.owner1_last_name,
			owner1_full_name = EXCLUDED.owner1_full_name,
			owner1_type = EXCLUDED.owner1_type,
			owner2_first_name = EXCLUDED.owner2_first_name,
			owner2_last_name = EXCLUDED.owner2_last_name,
			owner_occupied = EXCLUDED.owner_occupied,
			absentee_owner = EXCLUDED.absentee_owner,
			investor_buyer = EXCLUDED.investor_buyer,
			corporate_owned = EXCLUDED.corporate_owned,
			ownership_length_months = EXCLUDED.ownership_length_months,
			mail_street = EXCLUDED.mail_street,
			mail_city = EXCLUDED.mail_city,
			mail_state = EXCLUDED.mail_state,
			mail_zip = EXCLUDED.mail_zip,
			suggested_rent = EXCLUDED.suggested_rent,
			median_income = EXCLUDED.median_income,
			flood_zone = EXCLUDED.flood_zone,
			hoa = EXCLUDED.hoa,
			tax_lien = EXCLUDED.tax_lien,
			tax_amount = EXCLUDED.tax_amount,
			tax_delinquent_year = EXCLUDED.tax_delinquent_year,
			years_owned = EXCLUDED.years_owned,
			last_sale_date = EXCLUDED.last_sale_date,
			raw_data = EXCLUDED.raw_data,
			enriched_at = EXCLUDED.enriched_at,
			updated_at = NOW()
		RETURNING id
	`

	// Build location string for PostGIS if coordinates exist
	var locationStr *string
	if prop.Location != nil {
		locationStr = prop.Location
	}

	err := db.conn.QueryRow(
		query,
		prop.Source, prop.SourceID, prop.Address, prop.Street, prop.City, prop.State, prop.Zip, prop.County, prop.FIPS,
		locationStr,
		prop.PropertyType, prop.Bedrooms, prop.Bathrooms, prop.Sqft, prop.YearBuilt, prop.LotSqft,
		prop.AVM, prop.AssessedValue,
		prop.OpenMortgageBalance, prop.EquityAmount, prop.EquityPercent, prop.LoanType, prop.LenderName,
		prop.PreForeclosure, prop.Auction, prop.AuctionDate, prop.AuctionTime, prop.AuctionLocation,
		prop.NoticeType, prop.RecordingDate, prop.Foreclosure,
		prop.BankEstimatedValue, prop.TrusteeName, prop.TrusteePhone, prop.ForeclosureDocumentType,
		prop.Owner1FirstName, prop.Owner1LastName, prop.Owner1FullName, prop.Owner1Type,
		prop.Owner2FirstName, prop.Owner2LastName,
		prop.OwnerOccupied, prop.AbsenteeOwner, prop.InvestorBuyer, prop.CorporateOwned,
		prop.OwnershipLengthMonths, prop.MailStreet, prop.MailCity, prop.MailState, prop.MailZip,
		prop.SuggestedRent, prop.MedianIncome, prop.FloodZone, prop.HOA, prop.TaxLien,
		prop.TaxAmount, prop.TaxDelinquentYear, prop.YearsOwned, prop.LastSaleDate,
		prop.RawData, prop.DiscoveredAt, prop.EnrichedAt,
	).Scan(&prop.ID)

	if err != nil {
		return fmt.Errorf("failed to upsert property: %w", err)
	}

	return nil
}

// UpdatePropertySkipTraceStatus updates skip trace tracking fields
func (db *DB) UpdatePropertySkipTraceStatus(propertyID uuid.UUID, status string, errorMsg string) error {
	var errorMsgPtr *string
	if errorMsg != "" {
		errorMsgPtr = &errorMsg
	}

	query := `
		UPDATE properties 
		SET 
			skip_trace_status = $1,
			skip_trace_attempted_at = NOW(),
			skip_trace_error = $2,
			updated_at = NOW()
		WHERE id = $3
	`

	_, err := db.conn.Exec(query, status, errorMsgPtr, propertyID)
	if err != nil {
		return fmt.Errorf("failed to update skip trace status: %w", err)
	}

	return nil
}
