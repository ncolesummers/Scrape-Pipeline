package main

import (
	"flag"
	"os"
	"testing"
)

func TestConfigFlagParsing(t *testing.T) {
	// Save original args and flags
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Reset the flag parsing state
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test custom configuration path
	os.Args = []string{"cmd", "-config", "custom-config.yaml"}
	configPath := parseFlags()

	if configPath != "custom-config.yaml" {
		t.Errorf("Expected config path to be 'custom-config.yaml', got '%s'", configPath)
	}

	// Reset the flag parsing state again
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Test default configuration path
	os.Args = []string{"cmd"}
	configPath = parseFlags()

	if configPath != "config.yaml" {
		t.Errorf("Expected default config path to be 'config.yaml', got '%s'", configPath)
	}
}

func TestConfigFileNotFound(t *testing.T) {
	// Test loading a non-existent config file
	nonExistentPath := "non-existent-config.yaml"
	_, err := loadConfig(nonExistentPath)

	if err == nil {
		t.Error("Expected error when loading non-existent config file, got nil")
	}
}

// Skip the actual loadConfig testing here since it requires mocking the file system
// or creating a real temporary file. In a real implementation, you would use
// the same approach as in the previous test, but use a proper config file.
