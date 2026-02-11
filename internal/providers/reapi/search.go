package reapi

import (
	"encoding/json"
	"fmt"
	"time"
)

// SearchPropertyIDs searches for properties and returns only IDs (FREE - no credits)
func (c *Client) SearchPropertyIDs(state string, auctionDateMin, auctionDateMax time.Time, equityMin int) ([]int64, error) {
	var allIDs []int64
	start := 0

	for {
		req := SearchRequest{
			Size:             250,
			Start:            start,
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

		allIDs = append(allIDs, resp.Data...)

		// Check if we need to paginate
		if start+len(resp.Data) >= resp.ResultCount {
			break
		}

		start += 250
		time.Sleep(rateLimitDelay) // Rate limit between pagination calls
	}

	return allIDs, nil
}
