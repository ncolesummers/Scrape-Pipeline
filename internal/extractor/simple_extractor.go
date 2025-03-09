// This is a simple HTML content extractor implementation for testing purposes.
// For production use, consider using the more advanced implementation in pkg/extractor/extractor.go
// which uses the go-readability library for better content extraction.

package extractor

import (
	"strings"

	"golang.org/x/net/html"

	"github.com/ncolesummers/scrape-pipeline/internal/config"
	"github.com/ncolesummers/scrape-pipeline/internal/scraper"
)

// Image represents an extracted image
type Image struct {
	URL         string
	Alt         string
	Description string
}

// Content represents the extracted content from HTML
type Content struct {
	Metadata  map[string]string
	Title     string
	Content   string
	URL       string
	Images    []Image
	WordCount int
}

// Extractor is the interface for content extractors
type Extractor interface {
	// Extract extracts content from scraped HTML
	Extract(result *scraper.ScrapeResult) (*Content, error)
}

// SimpleExtractor implements the Extractor interface using basic HTML parsing
// This is a simplified version for testing purposes only
type SimpleExtractor struct {
	config config.ExtractionConfig
}

// NewExtractor creates a new extractor instance based on the provided configuration
func NewExtractor(cfg config.ExtractionConfig) (Extractor, error) {
	return &SimpleExtractor{
		config: cfg,
	}, nil
}

// Extract extracts content from scraped HTML
func (e *SimpleExtractor) Extract(result *scraper.ScrapeResult) (*Content, error) {
	// In a real implementation, we would use the go-readability library
	// For this test implementation, we'll use a simple HTML parser

	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(result.HTML))
	if err != nil {
		return nil, err
	}

	// Initialize the content
	content := &Content{
		URL:      result.URL,
		Metadata: make(map[string]string),
		Images:   []Image{},
	}

	// Extract metadata and content
	e.extractTitle(doc, content)
	e.extractMetadata(doc, content)
	e.extractMainContent(doc, content)

	if e.config.ExtractImages {
		e.extractImages(doc, content)
	}

	// Count words in content
	content.WordCount = len(strings.Fields(content.Content))

	return content, nil
}

func (e *SimpleExtractor) extractTitle(doc *html.Node, content *Content) {
	var extractTitle func(*html.Node)
	extractTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				content.Title = n.FirstChild.Data
				content.Metadata["title"] = content.Title
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractTitle(c)
		}
	}
	extractTitle(doc)
}

func (e *SimpleExtractor) extractMetadata(doc *html.Node, content *Content) {
	var extractMeta func(*html.Node)
	extractMeta = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			var name, property, contentValue string
			for _, attr := range n.Attr {
				if attr.Key == "name" {
					name = attr.Val
				} else if attr.Key == "property" {
					property = attr.Val
				} else if attr.Key == "content" {
					contentValue = attr.Val
				}
			}

			if name != "" && contentValue != "" {
				content.Metadata[name] = contentValue
			} else if property != "" && contentValue != "" {
				if strings.HasPrefix(property, "og:") {
					content.Metadata[property] = contentValue
				} else if property == "article:published_time" {
					content.Metadata["published_time"] = contentValue
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractMeta(c)
		}
	}
	extractMeta(doc)
}

func (e *SimpleExtractor) extractMainContent(doc *html.Node, content *Content) {
	// In a real implementation, this would be much more sophisticated
	// For now, we'll just extract the content from the main article tag
	var sb strings.Builder
	var extractArticle func(*html.Node)
	extractArticle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "article" {
			e.extractNodeText(n, &sb, true)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractArticle(c)
		}
	}
	extractArticle(doc)
	content.Content = sb.String()
}

func (e *SimpleExtractor) extractNodeText(n *html.Node, sb *strings.Builder, isRoot bool) {
	// Skip non-content elements
	if n.Type == html.ElementNode {
		if n.Data == "script" || n.Data == "style" || n.Data == "nav" || n.Data == "footer" {
			return
		}

		// Handle headings
		if e.config.PreserveHeadings && strings.HasPrefix(n.Data, "h") && len(n.Data) == 2 {
			if !isRoot {
				sb.WriteString("\n\n")
			}
		}
	}

	// Add text content
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			sb.WriteString(text)
			sb.WriteString(" ")
		}
	}

	// Process child nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		e.extractNodeText(c, sb, false)
	}

	// Add line breaks after certain elements
	if n.Type == html.ElementNode {
		if n.Data == "p" || n.Data == "div" || n.Data == "br" ||
			n.Data == "li" || n.Data == "h1" || n.Data == "h2" ||
			n.Data == "h3" || n.Data == "h4" || n.Data == "h5" || n.Data == "h6" {
			sb.WriteString("\n")
		}
	}
}

func (e *SimpleExtractor) extractImages(doc *html.Node, content *Content) {
	var extractImgs func(*html.Node)
	extractImgs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			var src, alt string
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					src = attr.Val
				} else if attr.Key == "alt" {
					alt = attr.Val
				}
			}

			if src != "" {
				content.Images = append(content.Images, Image{
					URL: src,
					Alt: alt,
				})
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractImgs(c)
		}
	}
	extractImgs(doc)
}
