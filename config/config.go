package config

import "os"

// Config holds application configuration loaded from environment variables.
type Config struct {
	DatabaseURL string
	Port        string
}

// Load reads configuration from environment variables, applying defaults where necessary.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        port,
	}
}
