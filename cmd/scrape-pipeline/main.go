package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ncolesummers/scrape-pipeline/internal/config"
)

func main() {
	// Parse command line flags
	configPath := parseFlags()

	fmt.Println("Starting Web Scraping and RAG System Pipeline")
	fmt.Printf("Using configuration file: %s\n", configPath)

	// Load configuration from file
	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	setupSignalHandler(cancel)

	// Initialize and run the pipeline
	fmt.Println("Initialized pipeline with configuration")
	fmt.Printf("Configured scrapers: %d\n", len(cfg.Scrapers))

	// Use the context in a simple timeout simulation
	// In the real implementation, this would be used to control the pipeline
	runPipeline(ctx)
}

// runPipeline runs the scraping pipeline with the given context
func runPipeline(ctx context.Context) {
	fmt.Println("Pipeline started")

	// Simulate pipeline operation
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("Pipeline completed successfully")
	case <-ctx.Done():
		fmt.Println("Pipeline was cancelled")
	}
}

// parseFlags parses command line flags and returns the path to the configuration file
func parseFlags() string {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()
	return *configPath
}

// loadConfig loads the configuration from the specified file
func loadConfig(path string) (*config.Config, error) {
	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file does not exist: %s", path)
	}

	// Load the configuration
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return cfg, nil
}

// setupSignalHandler sets up a handler for OS signals
func setupSignalHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nReceived shutdown signal. Shutting down gracefully...")
		cancel()
	}()
}
