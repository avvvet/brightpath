package pipeline

import (
	"time"

	"github.com/avvvet/brightpath/internal/domain"
)

// PassesOwnerFilter checks if property has a valid individual owner
// Filters out trusts, HOAs, corporations that slip through corporate_owned flag
func PassesOwnerFilter(prop *domain.Property) bool {
	// Must have owner first name
	if prop.Owner1FirstName == nil || *prop.Owner1FirstName == "" {
		return false
	}

	// Must be Individual type (not Trust, Corporation, etc)
	if prop.Owner1Type == nil || *prop.Owner1Type != "Individual" {
		return false
	}

	return true
}

// PassesInvestorFilter checks if property owner is not an investor
// Retail homeowners are more motivated than investors
func PassesInvestorFilter(prop *domain.Property) bool {
	if prop.InvestorBuyer == nil {
		return true // If unknown, allow through
	}

	return !*prop.InvestorBuyer
}

// PassesAllFilters runs all Go post-filters
func PassesAllFilters(prop *domain.Property) bool {
	return PassesOwnerFilter(prop) && PassesInvestorFilter(prop)
}

// CalculateDaysToAuction returns days until auction from now
func CalculateDaysToAuction(auctionDate time.Time) int {
	now := time.Now()
	duration := auctionDate.Sub(now)
	days := int(duration.Hours() / 24)
	return days
}

// GetPriorityTier returns priority tier based on days to auction and notice type
// Tier 1 (highest): 8-14 days + NTS
// Tier 2: 3-7 days
// Tier 3 (lowest): 15-21 days
func GetPriorityTier(daysToAuction int, noticeType string) int {
	if daysToAuction >= 8 && daysToAuction <= 14 && noticeType == "NTS" {
		return 1
	}
	if daysToAuction >= 3 && daysToAuction <= 7 {
		return 2
	}
	if daysToAuction >= 15 && daysToAuction <= 21 {
		return 3
	}
	return 3 // default to lowest priority
}
