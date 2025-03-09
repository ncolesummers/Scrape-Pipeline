package models

import (
	"context"
)

// RawContent represents the raw content fetched from a URL
type RawContent struct {
	Headers     map[string]string
	URL         string
	HTML        string
	ContentType string
	Timestamp   int64
	StatusCode  int
}

// ExtractedContent represents content after extraction from HTML
type ExtractedContent struct {
	URL       string
	Title     string
	Content   string
	Author    string
	Published string
	Updated   string
	Language  string
	Tags      []string
	Images    []ImageInfo
}

// ImageInfo represents metadata about an image in the content
type ImageInfo struct {
	URL         string
	Alt         string
	Description string
	Width       int
	Height      int
}

// NormalizedContent represents content after normalization
type NormalizedContent struct {
	ID       string
	Original *ExtractedContent
	Text     string
}

// ContentChunk represents a chunk of content ready for embedding
type ContentChunk struct {
	ID       string
	Text     string
	Metadata map[string]interface{}
	Source   string
	Start    int
	End      int
}

// VectorEmbedding represents a vector embedding of a content chunk
type VectorEmbedding struct {
	Chunk   *ContentChunk
	ID      string
	Model   string
	Version string
	Vector  []float32
}

// Scraper defines the interface for the web scraping module
type Scraper interface {
	// Scrape extracts content from URLs and returns the raw content
	Scrape(ctx context.Context, urls []string) (<-chan *RawContent, <-chan error)

	// AddURLs adds URLs to the scraping queue
	AddURLs(urls []string) error

	// SetRateLimit sets the rate limit for scraping per domain
	SetRateLimit(requestsPerSecond float64) error
}

// Extractor defines the interface for the content extraction module
type Extractor interface {
	// Extract extracts the main content from raw HTML
	Extract(ctx context.Context, rawContent *RawContent) (*ExtractedContent, error)
}

// Normalizer defines the interface for the text normalization module
type Normalizer interface {
	// Normalize normalizes the extracted content
	Normalize(ctx context.Context, content *ExtractedContent) (*NormalizedContent, error)
}

// Chunker defines the interface for the document chunking module
type Chunker interface {
	// Chunk splits the normalized content into chunks
	Chunk(ctx context.Context, content *NormalizedContent) ([]*ContentChunk, error)
}

// QualityControl defines the interface for the quality control module
type QualityControl interface {
	// Check evaluates the quality of content chunks
	Check(ctx context.Context, chunks []*ContentChunk) ([]*ContentChunk, []error)

	// IsDuplicate checks if a chunk is a duplicate of existing content
	IsDuplicate(ctx context.Context, chunk *ContentChunk) (bool, float64, error)
}

// Embedder defines the interface for the embedding service module
type Embedder interface {
	// Embed converts text chunks to vector embeddings
	Embed(ctx context.Context, chunks []*ContentChunk) ([]*VectorEmbedding, error)
}

// VectorStorage defines the interface for the vector storage module
type VectorStorage interface {
	// Store saves vector embeddings to storage
	Store(ctx context.Context, embeddings []*VectorEmbedding) error

	// Query searches for similar vectors
	Query(ctx context.Context, queryVector []float32, limit int) ([]*VectorEmbedding, error)

	// Delete removes embeddings from storage
	Delete(ctx context.Context, ids []string) error
}

// Observer defines the interface for the observability module
type Observer interface {
	// RecordMetric records a metric
	RecordMetric(name string, value float64, labels map[string]string) error

	// StartSpan starts a tracing span
	StartSpan(ctx context.Context, name string) (context.Context, func())

	// Log logs a message with the given level and fields
	Log(level string, msg string, fields map[string]interface{})
}
