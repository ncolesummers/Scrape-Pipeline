package scraper

import (
	"context"
	"regexp"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"

	"github.com/ncolesummers/scrape-pipeline/internal/models"
)

// CollyScraper implements the Scraper interface using Colly
type CollyScraper struct {
	collector     *colly.Collector
	rateLimiter   *colly.LimitRule
	userAgent     string
	domainFilters []string
	urlFilters    []string
	mutex         sync.Mutex
}

// Config holds configuration for the scraper
type Config struct {
	RateLimitRules     map[string]float64
	UserAgent          string
	AllowedDomains     []string
	ProxyURLs          []string
	DenyURLPatterns    []string
	AllowURLPatterns   []string
	MaxConcurrency     int
	RetryCount         int
	RetryDelaySeconds  int
	TimeoutSeconds     int
	RateLimitPerDomain float64
	MaxDepth           int
	RespectRobotsTxt   bool
}

// NewCollyScraper creates a new CollyScraper with the given configuration
func NewCollyScraper(config Config) (*CollyScraper, error) {
	c := colly.NewCollector(
		colly.UserAgent(config.UserAgent),
		colly.MaxDepth(config.MaxDepth),
	)

	// Configure the collector
	c.AllowURLRevisit = false
	c.ParseHTTPErrorResponse = true

	// Set allowed domains if specified
	if len(config.AllowedDomains) > 0 {
		c.AllowedDomains = config.AllowedDomains
	}

	// Add URL filters
	for _, pattern := range config.AllowURLPatterns {
		c.URLFilters = append(c.URLFilters, regexp.MustCompile(pattern))
	}

	// Add URL deny filters using callbacks
	if len(config.DenyURLPatterns) > 0 {
		c.OnRequest(func(r *colly.Request) {
			for _, pattern := range config.DenyURLPatterns {
				if regexp.MustCompile(pattern).MatchString(r.URL.String()) {
					r.Abort()
					return
				}
			}
		})
	}

	// Set up rate limiting
	if err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: config.MaxConcurrency,
		Delay:       time.Duration(1000/config.RateLimitPerDomain) * time.Millisecond,
		RandomDelay: time.Duration(500/config.RateLimitPerDomain) * time.Millisecond,
	}); err != nil {
		return nil, err
	}

	// Set up domain-specific rate limiting
	for domain, rateLimit := range config.RateLimitRules {
		if err := c.Limit(&colly.LimitRule{
			DomainGlob:  domain,
			Parallelism: config.MaxConcurrency,
			Delay:       time.Duration(1000/rateLimit) * time.Millisecond,
			RandomDelay: time.Duration(500/rateLimit) * time.Millisecond,
		}); err != nil {
			return nil, err
		}
	}

	// Respect robots.txt if configured
	if config.RespectRobotsTxt {
		extensions.RandomUserAgent(c)
		extensions.Referer(c)
	}

	// Set timeout
	c.SetRequestTimeout(time.Duration(config.TimeoutSeconds) * time.Second)

	// Set up retry with callbacks
	if config.RetryCount > 0 {
		// Use Colly's retry extension
		c.OnError(func(r *colly.Response, err error) {
			if r.StatusCode >= 500 || r.StatusCode == 0 || r.StatusCode == 429 {
				retries := 0
				if r.Request.Ctx.GetAny("retries") != nil {
					retries = r.Request.Ctx.GetAny("retries").(int)
				}
				if retries < config.RetryCount {
					retries++
					r.Request.Ctx.Put("retries", retries)
					time.Sleep(time.Duration(config.RetryDelaySeconds) * time.Second)
					r.Request.Retry()
				}
			}
		})
	}

	// Set up proxy if configured
	if len(config.ProxyURLs) > 0 {
		// Set proxy for each request
		c.OnRequest(func(r *colly.Request) {
			// Simple round-robin proxy selection
			proxyIndex := 0
			if r.Ctx.GetAny("proxy_index") != nil {
				proxyIndex = r.Ctx.GetAny("proxy_index").(int)
			}
			proxyURL := config.ProxyURLs[proxyIndex%len(config.ProxyURLs)]
			r.ProxyURL = proxyURL
			// Increment for next request
			r.Ctx.Put("proxy_index", proxyIndex+1)
		})
	}

	return &CollyScraper{
		collector:     c,
		userAgent:     config.UserAgent,
		domainFilters: config.AllowedDomains,
		urlFilters:    config.AllowURLPatterns,
	}, nil
}

// Scrape extracts content from URLs and returns the raw content
func (s *CollyScraper) Scrape(ctx context.Context, urls []string) (<-chan *models.RawContent, <-chan error) {
	contentChan := make(chan *models.RawContent)
	errorChan := make(chan error)

	go func() {
		defer close(contentChan)
		defer close(errorChan)

		// Set up collectors
		s.collector.OnResponse(func(r *colly.Response) {
			// Check if context is canceled
			select {
			case <-ctx.Done():
				return
			default:
				// Process the response
				content := &models.RawContent{
					URL:         r.Request.URL.String(),
					HTML:        string(r.Body),
					Timestamp:   time.Now().Unix(),
					StatusCode:  r.StatusCode,
					ContentType: r.Headers.Get("Content-Type"),
					Headers:     make(map[string]string),
				}

				// Copy headers (correctly handling http.Header)
				for name, values := range *r.Headers {
					if len(values) > 0 {
						content.Headers[name] = values[0]
					}
				}

				// Send content to channel
				select {
				case contentChan <- content:
				case <-ctx.Done():
					return
				}
			}
		})

		s.collector.OnError(func(r *colly.Response, err error) {
			// Check if context is canceled
			select {
			case <-ctx.Done():
				return
			default:
				// Send error to channel
				select {
				case errorChan <- err:
				case <-ctx.Done():
					return
				}
			}
		})

		// Start scraping
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.collector.Visit(url); err != nil {
					select {
					case errorChan <- err:
					case <-ctx.Done():
						return
					}
				}
			}
		}

		// Wait for scraping to complete
		s.collector.Wait()
	}()

	return contentChan, errorChan
}

// AddURLs adds URLs to the scraping queue
func (s *CollyScraper) AddURLs(urls []string) error {
	for _, url := range urls {
		if err := s.collector.Visit(url); err != nil {
			return err
		}
	}
	return nil
}

// SetRateLimit sets the rate limit for scraping per domain
func (s *CollyScraper) SetRateLimit(requestsPerSecond float64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       time.Duration(1000/requestsPerSecond) * time.Millisecond,
		RandomDelay: time.Duration(500/requestsPerSecond) * time.Millisecond,
	})
}
