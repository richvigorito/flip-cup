package main

import (
	"testing"
	"time"

	"flip-cup/internal/game"
)

func TestLoadRuntimeConfigDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("GAME_CLEANUP_INTERVAL", "")
	t.Setenv("GAME_STALE_AFTER", "")

	config, err := loadRuntimeConfig()
	if err != nil {
		t.Fatalf("expected defaults to load, got error: %v", err)
	}

	if config.port != "8080" {
		t.Fatalf("expected default port 8080, got %s", config.port)
	}

	if config.cleanupInterval != game.DefaultCleanupInterval {
		t.Fatalf("expected cleanup interval %s, got %s", game.DefaultCleanupInterval, config.cleanupInterval)
	}

	if config.staleAfter != game.DefaultStaleAfter {
		t.Fatalf("expected stale threshold %s, got %s", game.DefaultStaleAfter, config.staleAfter)
	}
}

func TestLoadRuntimeConfigFromEnv(t *testing.T) {
	t.Setenv("PORT", "19090")
	t.Setenv("GAME_CLEANUP_INTERVAL", "15m")
	t.Setenv("GAME_STALE_AFTER", "45m")

	config, err := loadRuntimeConfig()
	if err != nil {
		t.Fatalf("expected config from env to load, got error: %v", err)
	}

	if config.port != "19090" {
		t.Fatalf("expected port 19090, got %s", config.port)
	}

	if config.cleanupInterval != 15*time.Minute {
		t.Fatalf("expected cleanup interval 15m, got %s", config.cleanupInterval)
	}

	if config.staleAfter != 45*time.Minute {
		t.Fatalf("expected stale threshold 45m, got %s", config.staleAfter)
	}
}

func TestLoadRuntimeConfigRejectsInvalidDuration(t *testing.T) {
	t.Setenv("GAME_CLEANUP_INTERVAL", "not-a-duration")

	_, err := loadRuntimeConfig()
	if err == nil {
		t.Fatal("expected invalid duration error")
	}
}

func TestLoadRuntimeConfigRejectsNonPositiveDuration(t *testing.T) {
	t.Setenv("GAME_STALE_AFTER", "0s")

	_, err := loadRuntimeConfig()
	if err == nil {
		t.Fatal("expected non-positive duration error")
	}
}
