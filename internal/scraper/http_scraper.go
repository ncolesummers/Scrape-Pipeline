// This is a simple HTTP-based scraper implementation for testing purposes.
// For production use, consider using the more advanced implementation in pkg/scraper/scraper.go

package scraper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ncolesummers/scrape-pipeline/internal/config"
)

// ScrapeResult represents the result of a web scrape
type ScrapeResult struct {
	FetchedAt time.Time
	Headers   map[string]string
	URL       string
	HTML      string
	Status    int
}

// Scraper is the interface for web scrapers
type Scraper interface {
	// Name returns the name of the scraper
	Name() string
	// Scrape fetches the content of a URL
	Scrape(url string) (*ScrapeResult, error)
}

// HTTPScraper implements the Scraper interface using standard HTTP
type HTTPScraper struct {
	client          *http.Client
	name            string
	baseURL         string
	userAgent       string
	disallowedPaths []string
	rateLimit       float64
	concurrency     int
	respectRobots   bool
}

// NewScraper creates a new scraper instance based on the provided configuration
func NewScraper(cfg config.ScraperConfig) (Scraper, error) {
	// Create an HTTP client with reasonable timeout defaults
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	scraper := &HTTPScraper{
		name:          cfg.Name,
		baseURL:       cfg.URL,
		client:        client,
		userAgent:     cfg.UserAgent,
		rateLimit:     float64(cfg.RateLimit),
		concurrency:   cfg.Concurrency,
		respectRobots: cfg.RespectRobotsTxt,
	}

	// If we respect robots.txt, fetch and parse it
	if cfg.RespectRobotsTxt {
		// In a real implementation, we would fetch and parse the robots.txt file
		// For now, just simulate it with a basic implementation
		scraper.disallowedPaths = []string{"/private/"}
	}

	return scraper, nil
}

// Name returns the name of the scraper
func (s *HTTPScraper) Name() string {
	return s.name
}

// Scrape fetches the content of a URL
func (s *HTTPScraper) Scrape(url string) (*ScrapeResult, error) {
	// Simple robots.txt check - in a real implementation this would use a proper parser
	if s.respectRobots {
		for _, path := range s.disallowedPaths {
			if len(url) >= len(s.baseURL)+len(path) &&
				url[len(s.baseURL):len(s.baseURL)+len(path)] == path {
				return nil, errors.New("URL is disallowed by robots.txt")
			}
		}
	}

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml")

	// Send the request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Collect headers
	headers := make(map[string]string)
	for name, values := range resp.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}

	// Create and return the result
	result := &ScrapeResult{
		URL:       url,
		HTML:      string(body),
		Headers:   headers,
		Status:    resp.StatusCode,
		FetchedAt: time.Now(),
	}

	return result, nil
}
