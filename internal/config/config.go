package config

import (
	"encoding/json"
	"os"
)

// Config holds the application configuration
type Config struct {
	YouTubeAPIKey     string `json:"youtube_api_key"`
	LiveChatId        string `json:"live_chat_id"`
	ProfanityDatabase string `json:"profanity_database"`
	MinSeverity       string `json:"min_severity"`
	// In how many seconds the bot will check for new messages from the API
	PollingInterval int `json:"polling_interval"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set defaults if not specified
	if config.PollingInterval == 0 {
		config.PollingInterval = 5 // Default to 5 seconds
	}

	if config.MinSeverity == "" {
		config.MinSeverity = "mild" // Default to mild severity
	}

	if config.ProfanityDatabase == "" {
		config.ProfanityDatabase = "data/profanity.json" // Default location
	}

	return &config, nil
}
