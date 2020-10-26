package config

import (
	"os"
	"strconv"
)

type Config struct {
	PollInterval    int
	AzureStatusURL  string
	GitHubStatusURL string
}

func New() *Config {
	return &Config{
		PollInterval:    getEnvAsInt("CLOUDSTATUS_PollInterval", 3),
		AzureStatusURL:  getEnv("CLOUDSTATUS_AzureStatusURL", "https://status.dev.azure.com/_apis/status/health?geographies=EU,US&api-version=6.0-preview.1"),
		GitHubStatusURL: getEnv("CloudStatus_GitHubStatusURL", "https://kctbh9vrtdwd.statuspage.io/api/v2/summary.json"),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
