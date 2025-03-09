package config

import (
	"errors"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Storage    StorageConfig    `yaml:"storage"`
	Scrapers   []ScraperConfig  `yaml:"scrapers"`
	Embedding  EmbeddingConfig  `yaml:"embedding"`
	Chunking   ChunkingConfig   `yaml:"chunking"`
	Quality    QualityConfig    `yaml:"quality"`
	Extraction ExtractionConfig `yaml:"extraction"`
}

// ScraperConfig contains configuration for a web scraper
type ScraperConfig struct {
	Name             string `yaml:"name"`
	URL              string `yaml:"url"`
	UserAgent        string `yaml:"user_agent"`
	RateLimit        int    `yaml:"rate_limit"`
	Concurrency      int    `yaml:"concurrency"`
	RespectRobotsTxt bool   `yaml:"respect_robots_txt"`
}

// ExtractionConfig contains configuration for content extraction
type ExtractionConfig struct {
	PreserveHeadings bool `yaml:"preserve_headings"`
	ExtractImages    bool `yaml:"extract_images"`
}

// ChunkingConfig contains configuration for document chunking
type ChunkingConfig struct {
	MaxTokens int `yaml:"max_tokens"`
	Overlap   int `yaml:"overlap"`
}

// QualityConfig contains configuration for quality control
type QualityConfig struct {
	MinContentLength   int     `yaml:"min_content_length"`
	DuplicateThreshold float64 `yaml:"duplicate_threshold"`
}

// EmbeddingConfig contains configuration for embedding generation
type EmbeddingConfig struct {
	Model     string `yaml:"model"`
	BatchSize int    `yaml:"batch_size"`
}

// StorageConfig contains configuration for vector storage
type StorageConfig struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

// LoadConfig loads configuration from a file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if len(c.Scrapers) == 0 {
		return errors.New("at least one scraper configuration is required")
	}

	for i, scraper := range c.Scrapers {
		if scraper.Name == "" {
			return fmt.Errorf("scraper #%d is missing a name", i+1)
		}
		if scraper.URL == "" {
			return fmt.Errorf("scraper '%s' is missing a URL", scraper.Name)
		}
		if scraper.RateLimit <= 0 {
			// Set default rate limit if invalid
			c.Scrapers[i].RateLimit = 1
		}
		if scraper.Concurrency <= 0 {
			// Set default concurrency if invalid
			c.Scrapers[i].Concurrency = 1
		}
		if scraper.UserAgent == "" {
			// Set default user agent if missing
			c.Scrapers[i].UserAgent = "Scrape-Pipeline/1.0"
		}
	}

	return nil
}

// WriteDefaultConfig writes a default configuration to a file
func WriteDefaultConfig(filepath string) error {
	// Create a default configuration
	config := &Config{
		Scrapers: []ScraperConfig{
			{
				Name:             "Default Scraper",
				URL:              "https://example.com",
				RateLimit:        1,
				Concurrency:      1,
				RespectRobotsTxt: true,
			},
		},
		Extraction: ExtractionConfig{
			PreserveHeadings: true,
			ExtractImages:    true,
		},
		Chunking: ChunkingConfig{
			MaxTokens: 1000,
			Overlap:   200,
		},
		Quality: QualityConfig{
			MinContentLength:   100,
			DuplicateThreshold: 0.5,
		},
		Embedding: EmbeddingConfig{
			Model:     "default_model",
			BatchSize: 32,
		},
		Storage: StorageConfig{
			Type: "local",
			Path: "./data",
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling default config: %w", err)
	}

	// Write to file
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing default config file: %w", err)
	}

	return nil
}
