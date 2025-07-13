package main

import (
	"testing"

	"ping-badge-be/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	// Test that config loads without error
	cfg := config.Load()
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.DatabaseURL)
	assert.NotEmpty(t, cfg.JWTSecret)
}
