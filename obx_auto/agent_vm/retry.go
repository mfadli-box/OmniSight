package main

import (
	"log"
	"time"
)

// RetryConfig defines retry behavior.
type RetryConfig struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

// DefaultRetry returns standard retry config: 3 attempts, 1s base, 10s max.
func DefaultRetry() RetryConfig {
	return RetryConfig{MaxRetries: 3, BaseDelay: 1 * time.Second, MaxDelay: 10 * time.Duration(1000)}
}

// WithRetry executes fn up to cfg.MaxRetries times with exponential backoff.
// Returns nil on first success. Returns last error after all retries exhausted.
func WithRetry(cfg RetryConfig, label string, fn func() error) error {
	var lastErr error
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt < cfg.MaxRetries {
			delay := cfg.BaseDelay * time.Duration(1<<uint(attempt))
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
			log.Printf("[agent] %s failed (attempt %d/%d): %v — retrying in %v", label, attempt+1, cfg.MaxRetries+1, lastErr, delay)
			time.Sleep(delay)
		}
	}
	log.Printf("[agent] %s failed after %d attempts: %v", label, cfg.MaxRetries+1, lastErr)
	return lastErr
}
