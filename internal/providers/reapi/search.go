package reapi

import (
	"encoding/json"
	"fmt"
	"time"
)

// SearchPropertyIDs searches for properties and returns only IDs (FREE - no credits)
// ids_only mode returns all results up to the size limit in a single call
func (c *Client) SearchPropertyIDs(state string, auctionDateMin, auctionDateMax time.Time, equityMin int) ([]int64, error) {
	req := SearchRequest{
		Size:             5000, // High limit to get all results in one call
		State:            state,
		Auction:          true,
		AuctionDateMin:   auctionDateMin.Format("2006-01-02"),
		AuctionDateMax:   auctionDateMax.Format("2006-01-02"),
		EquityPercentMin: equityMin,
		CorporateOwned:   false,
		PropertyType:     "SFR",
		IDsOnly:          true,
	}

	respBody, err := c.doRequest("POST", "/v2/PropertySearch", req)
	if err != nil {
		return nil, fmt.Errorf("property search failed: %w", err)
	}

	var resp SearchResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search API returned status %d", resp.StatusCode)
	}

	return resp.Data, nil
}
