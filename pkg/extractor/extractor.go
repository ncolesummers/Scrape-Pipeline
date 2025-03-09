package extractor

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	readability "github.com/go-shiori/go-readability"

	"github.com/ncolesummers/scrape-pipeline/internal/models"
)

// ReadabilityExtractor implements the Extractor interface using go-readability
type ReadabilityExtractor struct {
	siteRules       map[string]SiteRule
	extractImages   bool
	extractMetadata bool
}

// SiteRule defines custom extraction rules for a specific site
type SiteRule struct {
	ArticleSelector string
	TitleSelector   string
	AuthorSelector  string
	DateSelector    string
}

// Config holds configuration for the extractor
type Config struct {
	SiteSpecificRules map[string]SiteRule
	ExtractImages     bool
	ExtractMetadata   bool
}

// NewReadabilityExtractor creates a new ReadabilityExtractor with the given configuration
func NewReadabilityExtractor(config Config) *ReadabilityExtractor {
	return &ReadabilityExtractor{
		extractImages:   config.ExtractImages,
		extractMetadata: config.ExtractMetadata,
		siteRules:       config.SiteSpecificRules,
	}
}

// Extract extracts the main content from raw HTML
func (e *ReadabilityExtractor) Extract(ctx context.Context, rawContent *models.RawContent) (*models.ExtractedContent, error) {
	// Check if context is canceled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawContent.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Check if we have site-specific rules
	hostname := parsedURL.Hostname()

	// Apply site-specific rules if they exist
	// Note: In a more complete implementation, we would use these rules
	// to customize the extraction process
	_ = e.siteRules[hostname] // Placeholder for future implementation

	// Extract the article using go-readability
	article, err := readability.FromReader(strings.NewReader(rawContent.HTML), parsedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract content: %w", err)
	}

	// Create the extracted content
	extracted := &models.ExtractedContent{
		URL:     rawContent.URL,
		Title:   article.Title,
		Content: article.Content,
		// Extract more metadata if available
		Author:    article.Byline,
		Published: article.SiteName,
		Language:  detectLanguage(article.TextContent),
		Tags:      extractTags(article.TextContent),
	}

	// Extract images if configured
	if e.extractImages && article.Image != "" {
		imgURL := article.Image
		if !strings.HasPrefix(imgURL, "http") && !strings.HasPrefix(imgURL, "//") {
			// Convert relative URL to absolute
			base, err := url.Parse(rawContent.URL)
			if err == nil {
				imgURLObj, err := url.Parse(imgURL)
				if err == nil {
					imgURL = base.ResolveReference(imgURLObj).String()
				}
			}
		}

		extracted.Images = append(extracted.Images, models.ImageInfo{
			URL:         imgURL,
			Alt:         article.Title,
			Description: article.Excerpt,
		})
	}

	// Apply site-specific extraction for additional metadata
	if e.extractMetadata {
		e.extractAdditionalMetadata(ctx, extracted, rawContent.HTML, hostname)
	}

	return extracted, nil
}

// extractAdditionalMetadata extracts additional metadata from the HTML using site-specific rules
func (e *ReadabilityExtractor) extractAdditionalMetadata(ctx context.Context, content *models.ExtractedContent, html, hostname string) {
	// This is a placeholder for more sophisticated metadata extraction
	// In a real implementation, this would use the site-specific rules to extract
	// additional metadata like author, publication date, tags, etc.

	// For example, using a DOM parser to extract specific elements based on selectors
	// from the site-specific rules.
}

// detectLanguage is a simple placeholder for language detection
// In a real implementation, this would use a language detection library
func detectLanguage(text string) string {
	// This is a placeholder for more sophisticated language detection
	// For now, we'll just return "en" (English)
	return "en"
}

// extractTags is a simple placeholder for tag extraction
// In a real implementation, this would use a more sophisticated approach
func extractTags(text string) []string {
	// This is a placeholder for more sophisticated tag extraction
	return []string{}
}
