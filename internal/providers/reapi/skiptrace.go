package reapi

import (
	"encoding/json"
	"fmt"
)

// SkipTrace performs skip tracing to find contact information (1 credit per call)
func (c *Client) SkipTrace(req SkipTraceRequest) (*SkipTraceResponse, error) {
	respBody, err := c.doRequest("POST", "/v1/SkipTrace", req)
	if err != nil {
		return nil, fmt.Errorf("skip trace request failed: %w", err)
	}

	var resp SkipTraceResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse skip trace response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("skip trace API returned status %d", resp.StatusCode)
	}

	return &resp, nil
}
