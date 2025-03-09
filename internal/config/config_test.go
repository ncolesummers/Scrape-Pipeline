package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file for testing
	tempConfig := `
scrapers:
  - name: example-blog
    url: https://example.com/blog
    rate_limit: 1
    concurrency: 2
    user_agent: "Mozilla/5.0 Bot"
extraction:
  preserve_headings: true
  extract_images: false
`
	tempFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write([]byte(tempConfig)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test loading the config file
	config, err := LoadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify the loaded configuration
	if len(config.Scrapers) != 1 {
		t.Errorf("Expected 1 scraper, got %d", len(config.Scrapers))
	}

	if config.Scrapers[0].Name != "example-blog" {
		t.Errorf("Expected scraper name 'example-blog', got '%s'", config.Scrapers[0].Name)
	}

	if config.Scrapers[0].RateLimit != 1 {
		t.Errorf("Expected rate limit 1, got %d", config.Scrapers[0].RateLimit)
	}

	if config.Extraction.PreserveHeadings != true {
		t.Errorf("Expected preserve_headings to be true")
	}
}

func TestValidateConfig(t *testing.T) {
	// Test validation with valid configuration
	validConfig := Config{
		Scrapers: []ScraperConfig{
			{
				Name:        "test-blog",
				URL:         "https://test.com/blog",
				RateLimit:   2,
				Concurrency: 1,
				UserAgent:   "Test Bot",
			},
		},
	}

	if err := validConfig.Validate(); err != nil {
		t.Errorf("Valid config failed validation: %v", err)
	}

	// Test validation with invalid configuration
	invalidConfig := Config{
		Scrapers: []ScraperConfig{
			{
				Name:      "missing-url",
				RateLimit: 1,
				// URL is missing
			},
		},
	}

	if err := invalidConfig.Validate(); err == nil {
		t.Errorf("Invalid config passed validation when it should have failed")
	}
}

func TestExampleConfigFile(t *testing.T) {
	// Path to the example config file
	exampleConfigPath := "../../config.yaml.example"

	// Skip the test if the example config file doesn't exist
	if _, err := os.Stat(exampleConfigPath); os.IsNotExist(err) {
		t.Skipf("Example config file not found at %s, skipping test", exampleConfigPath)
		return
	}

	// Try to load the example config file
	_, err := LoadConfig(exampleConfigPath)
	if err != nil {
		t.Errorf("Failed to load example config file: %v", err)
	}
}
