package main

import (
	"fmt"
	"os"
)

type config struct {
	DatabaseURL   string
	HTTPAddr      string
	MigrationFile string
}

func loadConfig() (config, error) {
	cfg := config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		HTTPAddr:      getenvDefault("HTTP_ADDR", ":8080"),
		MigrationFile: getenvDefault("MIGRATION_FILE", "migrations/001_create_jobs.sql"),
	}

	if cfg.DatabaseURL == "" {
		return config{}, fmt.Errorf("DATABASE_URL is not set")
	}

	return cfg, nil
}

func getenvDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
