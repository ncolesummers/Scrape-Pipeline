package extractor

import (
	"testing"

	"github.com/ncolesummers/scrape-pipeline/internal/config"
	"github.com/ncolesummers/scrape-pipeline/internal/scraper"
)

// TestNewExtractor tests the creation of a new extractor
// This test uses the simple extractor implementation
func TestNewExtractor(t *testing.T) {
	cfg := config.ExtractionConfig{
		PreserveHeadings: true,
		ExtractImages:    false,
	}

	extractor, err := NewExtractor(cfg)
	if err != nil {
		t.Fatalf("Failed to create extractor: %v", err)
	}

	if extractor == nil {
		t.Fatal("Extractor is nil")
	}

	// Check that we get a SimpleExtractor
	_, ok := extractor.(*SimpleExtractor)
	if !ok {
		t.Errorf("Expected a SimpleExtractor implementation")
	}
}

// TestExtract tests the content extraction functionality
// This test uses the simple extractor implementation
func TestExtract(t *testing.T) {
	// Create test HTML
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Test Article</title>
    <meta name="author" content="Test Author">
    <meta name="description" content="Test Description">
</head>
<body>
    <header>
        <nav>
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="/about">About</a></li>
            </ul>
        </nav>
    </header>
    <main>
        <article>
            <h1>Main Article Heading</h1>
            <p>This is the first paragraph of the article. It contains some text that should be extracted.</p>
            <h2>Section Heading</h2>
            <p>This is the second paragraph with more <b>important content</b> that should be preserved.</p>
            <ul>
                <li>List item 1</li>
                <li>List item 2</li>
            </ul>
            <img src="/images/test.jpg" alt="Test Image">
        </article>
    </main>
    <footer>
        <p>Copyright 2023</p>
    </footer>
</body>
</html>`

	// Create a scrape result
	result := &scraper.ScrapeResult{
		URL:  "https://example.com/article",
		HTML: html,
	}

	// Create an extractor
	cfg := config.ExtractionConfig{
		PreserveHeadings: true,
		ExtractImages:    true,
	}

	extractor, err := NewExtractor(cfg)
	if err != nil {
		t.Fatalf("Failed to create extractor: %v", err)
	}

	// Extract content
	content, err := extractor.Extract(result)
	if err != nil {
		t.Fatalf("Failed to extract content: %v", err)
	}

	// Verify the extracted content
	if content == nil {
		t.Fatal("Extracted content is nil")
	}

	if content.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got '%s'", content.Title)
	}

	// Main content should contain the article text
	if !contains(content.Content, "Main Article Heading") {
		t.Errorf("Content should contain 'Main Article Heading'")
	}

	if !contains(content.Content, "This is the first paragraph") {
		t.Errorf("Content should contain the first paragraph")
	}

	// Headers should be preserved
	if cfg.PreserveHeadings && !contains(content.Content, "Section Heading") {
		t.Errorf("Content should preserve section headings")
	}

	// Ensure navigation and footer are removed
	if contains(content.Content, "Home") || contains(content.Content, "About") {
		t.Errorf("Navigation menu should be removed from content")
	}

	if contains(content.Content, "Copyright 2023") {
		t.Errorf("Footer should be removed from content")
	}

	// Check if images are preserved when requested
	if cfg.ExtractImages && len(content.Images) == 0 {
		t.Errorf("Images should be extracted when extract_images is true")
	}
}

// TestExtractMetadata tests metadata extraction capabilities
// This test uses the simple extractor implementation
func TestExtractMetadata(t *testing.T) {
	// Create test HTML with metadata
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Metadata Test</title>
    <meta name="author" content="John Doe">
    <meta name="description" content="This is a test article">
    <meta property="og:title" content="OpenGraph Title">
    <meta property="article:published_time" content="2023-05-15T12:00:00Z">
</head>
<body>
    <article>
        <h1>Article with Metadata</h1>
        <p>This article has metadata in the head section.</p>
    </article>
</body>
</html>`

	// Create a scrape result
	result := &scraper.ScrapeResult{
		URL:  "https://example.com/metadata-article",
		HTML: html,
	}

	// Create an extractor
	cfg := config.ExtractionConfig{
		PreserveHeadings: true,
		ExtractImages:    false,
	}

	extractor, err := NewExtractor(cfg)
	if err != nil {
		t.Fatalf("Failed to create extractor: %v", err)
	}

	// Extract content
	content, err := extractor.Extract(result)
	if err != nil {
		t.Fatalf("Failed to extract content: %v", err)
	}

	// Verify the extracted metadata
	if content.Metadata == nil {
		t.Fatal("Extracted metadata is nil")
	}

	if content.Metadata["title"] != "Metadata Test" {
		t.Errorf("Expected title 'Metadata Test', got '%s'", content.Metadata["title"])
	}

	if content.Metadata["author"] != "John Doe" {
		t.Errorf("Expected author 'John Doe', got '%s'", content.Metadata["author"])
	}

	if content.Metadata["description"] != "This is a test article" {
		t.Errorf("Expected description 'This is a test article', got '%s'", content.Metadata["description"])
	}

	if content.Metadata["og:title"] != "OpenGraph Title" {
		t.Errorf("Expected og:title 'OpenGraph Title', got '%s'", content.Metadata["og:title"])
	}

	if content.Metadata["published_time"] != "2023-05-15T12:00:00Z" {
		t.Errorf("Expected published_time '2023-05-15T12:00:00Z', got '%s'", content.Metadata["published_time"])
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && (s == substr || s != "" && s != substr && s[0:len(substr)] == substr)
}
