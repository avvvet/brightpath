package pipeline

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/avvvet/brightpath/internal/domain"
	"github.com/avvvet/brightpath/internal/providers/reapi"
	"github.com/google/uuid"
)

// MapPropertyDetailToProperty converts REAPI PropertyDetail response to domain.Property
func MapPropertyDetailToProperty(detail *reapi.PropertyDetailResponse) (*domain.Property, error) {
	data := detail.Data

	// Marshal raw data to JSONB
	rawData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal raw data: %w", err)
	}
	rawDataStr := string(rawData)

	// Parse auction date
	var auctionDate *time.Time
	if data.AuctionInfo.AuctionDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", data.AuctionInfo.AuctionDate)
		if err != nil {
			// Try alternative format without time
			t, err = time.Parse("2006-01-02", data.AuctionInfo.AuctionDate)
			if err == nil {
				auctionDate = &t
			}
		} else {
			auctionDate = &t
		}
	}

	// Parse recording date
	var recordingDate *time.Time
	if data.AuctionInfo.RecordingDate != "" {
		t, err := time.Parse("2006-01-02T15:04:05.000Z", data.AuctionInfo.RecordingDate)
		if err == nil {
			recordingDate = &t
		}
	}

	// Parse bank estimated value
	var bankEstimatedValue *int
	if data.AuctionInfo.EstimatedBankValue != "" {
		val, err := strconv.Atoi(data.AuctionInfo.EstimatedBankValue)
		if err == nil {
			bankEstimatedValue = &val
		}
	}

	// Build location string for PostGIS (POINT format)
	var location *string
	if data.PropertyInfo.Latitude != 0 && data.PropertyInfo.Longitude != 0 {
		locStr := fmt.Sprintf("POINT(%f %f)", data.PropertyInfo.Longitude, data.PropertyInfo.Latitude)
		location = &locStr
	}

	// Get loan type from first mortgage
	var loanType *string
	var lenderName *string
	if len(data.CurrentMortgages) > 0 {
		loanType = &data.CurrentMortgages[0].LoanType
		lenderName = &data.CurrentMortgages[0].LenderName
	}

	// Helper to get string pointer
	strPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	// Helper to get int pointer
	intPtr := func(i int) *int {
		if i == 0 {
			return nil
		}
		return &i
	}

	// Helper to get float64 pointer
	float64Ptr := func(f float64) *float64 {
		if f == 0 {
			return nil
		}
		return &f
	}

	// Helper to get bool pointer
	boolPtr := func(b bool) *bool {
		return &b
	}

	prop := &domain.Property{
		Source:   "reapi",
		SourceID: strconv.FormatInt(data.ID, 10),

		// Address
		Address:  data.PropertyInfo.Address.Address,
		Street:   strPtr(data.PropertyInfo.Address.Address),
		City:     strPtr(data.PropertyInfo.Address.City),
		State:    data.PropertyInfo.Address.State,
		Zip:      strPtr(data.PropertyInfo.Address.Zip),
		County:   strPtr(data.PropertyInfo.Address.County),
		Location: location,

		// Property details
		PropertyType: strPtr(data.PropertyType),
		Bedrooms:     intPtr(data.PropertyInfo.Bedrooms),
		Bathrooms:    float64Ptr(data.PropertyInfo.Bathrooms),
		Sqft:         intPtr(data.PropertyInfo.LivingSquareFeet),
		YearBuilt:    intPtr(data.PropertyInfo.YearBuilt),

		// Valuation
		AVM:           intPtr(data.EstimatedValue),
		AssessedValue: intPtr(data.TaxInfo.AssessedValue),

		// Mortgage & equity
		OpenMortgageBalance: intPtr(data.OpenMortgageBalance),
		EquityAmount:        intPtr(data.EstimatedEquity),
		EquityPercent:       intPtr(data.EquityPercent),
		LoanType:            loanType,
		LenderName:          lenderName,

		// Foreclosure
		PreForeclosure:  boolPtr(data.PreForeclosure),
		Auction:         boolPtr(data.Auction),
		AuctionDate:     auctionDate,
		AuctionTime:     strPtr(data.AuctionInfo.AuctionTime),
		AuctionLocation: strPtr(data.AuctionInfo.AuctionStreetAddress),
		NoticeType:      strPtr(data.NoticeType),
		RecordingDate:   recordingDate,

		// Auction enrichment
		BankEstimatedValue:      bankEstimatedValue,
		TrusteeName:             strPtr(data.AuctionInfo.TrusteeFullName),
		TrusteePhone:            strPtr(data.AuctionInfo.TrusteePhone),
		ForeclosureDocumentType: strPtr(data.AuctionInfo.DocumentType),

		// Owner
		Owner1FirstName:       strPtr(data.OwnerInfo.Owner1FirstName),
		Owner1LastName:        strPtr(data.OwnerInfo.Owner1LastName),
		Owner1FullName:        strPtr(data.OwnerInfo.Owner1FullName),
		Owner1Type:            strPtr(data.OwnerInfo.Owner1Type),
		Owner2FirstName:       strPtr(data.OwnerInfo.Owner2FirstName),
		Owner2LastName:        strPtr(data.OwnerInfo.Owner2LastName),
		OwnerOccupied:         boolPtr(data.OwnerOccupied),
		InvestorBuyer:         boolPtr(data.InvestorBuyer),
		CorporateOwned:        boolPtr(data.CorporateOwned),
		OwnershipLengthMonths: intPtr(data.OwnerInfo.OwnershipLength),
		MailStreet:            strPtr(data.OwnerInfo.MailAddress.Address),
		MailCity:              strPtr(data.OwnerInfo.MailAddress.City),
		MailState:             strPtr(data.OwnerInfo.MailAddress.State),
		MailZip:               strPtr(data.OwnerInfo.MailAddress.Zip),

		// Extras
		SuggestedRent:     strPtr(data.Demographics.SuggestedRent),
		MedianIncome:      strPtr(data.Demographics.MedianIncome),
		HOA:               boolPtr(data.PropertyInfo.HOA),
		TaxAmount:         strPtr(data.TaxInfo.TaxAmount),
		TaxDelinquentYear: strPtr(data.TaxInfo.TaxDelinquentYear),

		// Metadata
		RawData:      &rawDataStr,
		DiscoveredAt: time.Now(),
		EnrichedAt:   timePtr(time.Now()),
	}

	return prop, nil
}

// MapSkipTraceToLead converts REAPI SkipTrace response to domain.Lead
// Returns nil if no valid mobile phone found
// Strategy: Try owner first, then fallback to any valid mobile (relative/co-owner)
func MapSkipTraceToLead(propertyID uuid.UUID, ownerFirstName, ownerLastName string, skipTrace *reapi.SkipTraceResponse, daysToAuction int) (*domain.Lead, error) {
	if !skipTrace.Match {
		return nil, nil
	}

	var primaryPhone *reapi.Phone
	var matchedPersonID string
	var isOwnerMatch bool = false

	// First pass: Try to match owner by name
	for _, phone := range skipTrace.Output.Identity.Phones {
		// Apply phone filters first
		if phone.PhoneType != "mobile" || !phone.IsConnected || phone.DoNotCall {
			continue
		}

		// Check if this phone belongs to the owner
		for _, name := range skipTrace.Output.Identity.Names {
			if name.PersonID == phone.PersonID {
				if name.FirstName == ownerFirstName && name.LastName == ownerLastName {
					primaryPhone = &phone
					matchedPersonID = phone.PersonID
					isOwnerMatch = true
					break
				}
			}
		}
		if primaryPhone != nil {
			break
		}
	}

	// Second pass: If no owner phone found, take ANY valid mobile at the address (relative/co-owner)
	if primaryPhone == nil {
		for _, phone := range skipTrace.Output.Identity.Phones {
			if phone.PhoneType == "mobile" && phone.IsConnected && !phone.DoNotCall {
				primaryPhone = &phone
				matchedPersonID = phone.PersonID
				isOwnerMatch = false // This is a relative/co-owner, not the exact owner
				break
			}
		}
	}

	// If still no valid phone found, return nil (don't create lead)
	if primaryPhone == nil {
		return nil, nil
	}

	// Find primary email for matched person
	var primaryEmail *string
	for _, email := range skipTrace.Output.Identity.Emails {
		if email.PersonID == matchedPersonID {
			primaryEmail = &email.Email
			break
		}
	}

	// Marshal alt phones to JSON (all phones from skip trace for reference)
	var altPhones *string
	if len(skipTrace.Output.Identity.Phones) > 1 {
		altPhonesData, err := json.Marshal(skipTrace.Output.Identity.Phones)
		if err == nil {
			altPhonesStr := string(altPhonesData)
			altPhones = &altPhonesStr
		}
	}

	// Marshal alt emails to JSON
	var altEmails *string
	if len(skipTrace.Output.Identity.Emails) > 1 {
		altEmailsData, err := json.Marshal(skipTrace.Output.Identity.Emails)
		if err == nil {
			altEmailsStr := string(altEmailsData)
			altEmails = &altEmailsStr
		}
	}

	// Find demographics for matched person
	var ownerAge *int
	var ownerGender *string
	if skipTrace.Output.Demographics.Age > 0 {
		age := skipTrace.Output.Demographics.Age
		ownerAge = &age
	}
	if skipTrace.Output.Demographics.Gender != "" {
		ownerGender = &skipTrace.Output.Demographics.Gender
	}

	// Helper functions
	strPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	boolPtr := func(b bool) *bool {
		return &b
	}

	lead := &domain.Lead{
		PropertyID: propertyID,
		Status:     "new",

		// Contact (might be relative's phone if owner phone not found)
		Phone:           strPtr(primaryPhone.Phone),
		PhoneType:       strPtr(primaryPhone.PhoneType),
		PhoneConnected:  boolPtr(primaryPhone.IsConnected),
		PhoneDNC:        boolPtr(primaryPhone.DoNotCall),
		PhoneOwnerMatch: boolPtr(isOwnerMatch),
		AltPhones:       altPhones,
		Email:           primaryEmail,
		AltEmails:       altEmails,

		// Demographics
		OwnerAge:    ownerAge,
		OwnerGender: ownerGender,

		// Scoring
		DaysToAuction: &daysToAuction,
	}

	return lead, nil
}

// timePtr helper
func timePtr(t time.Time) *time.Time {
	return &t
}
