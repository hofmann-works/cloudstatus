package config

import (
	"os"
	"strconv"
)

type Config struct {
	PollInterval    int
	AzureStatusURL  string
	GitHubStatusURL string
	PGHost          string
	PGPort          int
	PGDatabase      string
	PGUser          string
	PGPassword      string
}

func New() *Config {
	return &Config{
		PollInterval:    getEnvAsInt("CLOUDSTATUS_PollInterval", 10),
		AzureStatusURL:  getEnv("CLOUDSTATUS_AzureStatusURL", "https://status.dev.azure.com/_apis/status/health?geographies=EU,US&api-version=6.0-preview.1"),
		GitHubStatusURL: getEnv("CLOUDSTATUS_GitHubStatusURL", "https://kctbh9vrtdwd.statuspage.io/api/v2/summary.json"),
		PGHost:          getEnv("CLOUDSTATUS_PGHost", "localhost"),
		PGPort:          getEnvAsInt("CLOUDSTATUS_PGPort", 5432),
		PGDatabase:      getEnv("CLOUDSTATUS_PGDatabase", "cloudstatus"),
		PGUser:          getEnv("CLOUDSTATUS_PGUser", "cloudstatus"),
		PGPassword:      getEnv("CLOUDSTATUS_PGPassword", "mypw"),
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
