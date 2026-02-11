package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	REAPIKey       string
	DatabaseURL    string
	PipelineStates []string
	AuctionDaysMin int
	AuctionDaysMax int
	EquityMin      int
	DailyHour      int
}

// Load reads environment variables and returns a Config
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{}

	// Required fields
	cfg.REAPIKey = os.Getenv("REAPI_API_KEY")
	if cfg.REAPIKey == "" {
		return nil, fmt.Errorf("REAPI_API_KEY is required")
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	// Parse pipeline states (comma-separated)
	statesStr := os.Getenv("PIPELINE_STATES")
	if statesStr == "" {
		statesStr = "TX" // default to Texas
	}
	cfg.PipelineStates = strings.Split(statesStr, ",")
	for i := range cfg.PipelineStates {
		cfg.PipelineStates[i] = strings.TrimSpace(cfg.PipelineStates[i])
	}

	// Parse integer fields with defaults
	var err error
	cfg.AuctionDaysMin, err = parseIntWithDefault("PIPELINE_AUCTION_DAYS_MIN", 3)
	if err != nil {
		return nil, err
	}

	cfg.AuctionDaysMax, err = parseIntWithDefault("PIPELINE_AUCTION_DAYS_MAX", 21)
	if err != nil {
		return nil, err
	}

	cfg.EquityMin, err = parseIntWithDefault("PIPELINE_EQUITY_MIN", 25)
	if err != nil {
		return nil, err
	}

	cfg.DailyHour, err = parseIntWithDefault("PIPELINE_DAILY_HOUR", 6)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// parseIntWithDefault parses an environment variable as int, with a default value
func parseIntWithDefault(key string, defaultValue int) (int, error) {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %s", key, valueStr)
	}

	return value, nil
}
