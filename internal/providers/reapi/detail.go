package reapi

import (
	"encoding/json"
	"fmt"
)

// GetPropertyDetail fetches full property details by ID (1 credit per call)
func (c *Client) GetPropertyDetail(id int64) (*PropertyDetailResponse, error) {
	req := PropertyDetailRequest{
		ID: id,
	}

	respBody, err := c.doRequest("POST", "/v2/PropertyDetail", req)
	if err != nil {
		return nil, fmt.Errorf("property detail request failed: %w", err)
	}

	var resp PropertyDetailResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse property detail response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("property detail API returned status %d", resp.StatusCode)
	}

	return &resp, nil
}
