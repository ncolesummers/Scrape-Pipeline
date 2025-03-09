package scraper

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ncolesummers/scrape-pipeline/internal/config"
)

// TestNewScraper tests the creation of a new scraper
// This test uses the simple HTTP scraper implementation in http_scraper.go
func TestNewScraper(t *testing.T) {
	cfg := config.ScraperConfig{
		Name:        "test-scraper",
		URL:         "https://example.com",
		RateLimit:   1,
		Concurrency: 2,
		UserAgent:   "Test Bot",
	}

	scraper, err := NewScraper(cfg)
	if err != nil {
		t.Fatalf("Failed to create scraper: %v", err)
	}

	if scraper == nil {
		t.Fatal("Scraper is nil")
	}

	if scraper.Name() != "test-scraper" {
		t.Errorf("Expected scraper name 'test-scraper', got '%s'", scraper.Name())
	}
}

// TestScrape tests the basic scraping functionality
// This test uses the simple HTTP scraper implementation in http_scraper.go
func TestScrape(t *testing.T) {
	// Create a test server to mock responses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
</head>
<body>
    <h1>Test Heading</h1>
    <p>This is a test paragraph.</p>
    <a href="/page2">Link to page 2</a>
</body>
</html>`))
	}))
	defer server.Close()

	// Create a scraper config using the test server URL
	cfg := config.ScraperConfig{
		Name:        "test-scraper",
		URL:         server.URL,
		RateLimit:   1,
		Concurrency: 1,
		UserAgent:   "Test Bot",
	}

	// Create a new scraper with this config
	scraper, err := NewScraper(cfg)
	if err != nil {
		t.Fatalf("Failed to create scraper: %v", err)
	}

	// Scrape the test URL
	result, err := scraper.Scrape(server.URL)
	if err != nil {
		t.Fatalf("Failed to scrape URL: %v", err)
	}

	// Verify the scraped content
	if result == nil {
		t.Fatal("Scrape result is nil")
	}

	if result.URL != server.URL {
		t.Errorf("Expected URL '%s', got '%s'", server.URL, result.URL)
	}

	if !strings.Contains(result.HTML, "Test Heading") {
		t.Errorf("Expected HTML to contain 'Test Heading'")
	}

	if !strings.Contains(result.HTML, "This is a test paragraph") {
		t.Errorf("Expected HTML to contain 'This is a test paragraph'")
	}
}

// TestRespectRobotsTxt tests robots.txt functionality
// This test uses the simple HTTP scraper implementation in http_scraper.go
func TestRespectRobotsTxt(t *testing.T) {
	// Create a mock robots.txt server
	robotsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			w.Write([]byte(`User-agent: *
Disallow: /private/
`))
		} else if r.URL.Path == "/allowed" {
			w.Write([]byte("Allowed content"))
		} else if r.URL.Path == "/private/disallowed" {
			w.Write([]byte("Disallowed content"))
		}
	}))
	defer robotsServer.Close()

	// Create a scraper config with robots.txt respect enabled
	cfg := config.ScraperConfig{
		Name:             "robots-test",
		URL:              robotsServer.URL,
		RateLimit:        1,
		Concurrency:      1,
		UserAgent:        "Test Bot",
		RespectRobotsTxt: true,
	}

	// Create a new scraper with this config
	scraper, err := NewScraper(cfg)
	if err != nil {
		t.Fatalf("Failed to create scraper: %v", err)
	}

	// Test scraping an allowed URL
	allowed, err := scraper.Scrape(robotsServer.URL + "/allowed")
	if err != nil {
		t.Fatalf("Failed to scrape allowed URL: %v", err)
	}
	if allowed == nil || allowed.HTML != "Allowed content" {
		t.Errorf("Expected to be able to scrape allowed URL")
	}

	// Test scraping a disallowed URL
	disallowed, err := scraper.Scrape(robotsServer.URL + "/private/disallowed")
	if err == nil {
		t.Errorf("Expected error when scraping disallowed URL, got none")
	}
	if disallowed != nil {
		t.Errorf("Expected nil result for disallowed URL, got content: %s", disallowed.HTML)
	}
}
