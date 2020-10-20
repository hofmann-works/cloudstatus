package config

import (
	"os"
	"strconv"
)

type Config struct {
	PollInterval int
}

func New() *Config {
	return &Config{
		PollInterval: getEnvAsInt("CLOUDSTATUS_PollInterval", 30),
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
