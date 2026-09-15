package main

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("MIGRATION_FILE", "")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() returned error: %v", err)
	}

	if cfg.DatabaseURL != "postgres://example" {
		t.Fatalf("expected database url %q, got %q", "postgres://example", cfg.DatabaseURL)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("expected http addr %q, got %q", ":8080", cfg.HTTPAddr)
	}
	if cfg.MigrationFile != "migrations/001_create_jobs.sql" {
		t.Fatalf("expected migration file %q, got %q", "migrations/001_create_jobs.sql", cfg.MigrationFile)
	}
}

func TestLoadConfigRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	_, err := loadConfig()
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is missing")
	}
}

func TestLoadConfigUsesOverrides(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("MIGRATION_FILE", "custom.sql")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("loadConfig() returned error: %v", err)
	}

	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("expected http addr override, got %q", cfg.HTTPAddr)
	}
	if cfg.MigrationFile != "custom.sql" {
		t.Fatalf("expected migration file override, got %q",
			cfg.MigrationFile)
	}
}
