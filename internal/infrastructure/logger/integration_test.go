package logger

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/mafzaidi/authorizer/internal/domain/service"
)

// TestLoggerImplementsInterface verifies that Logger implements service.Logger interface
func TestLoggerImplementsInterface(t *testing.T) {
	var _ service.Logger = (*Logger)(nil)
}

// TestLoggerAsInterface tests using the logger through the service.Logger interface
func TestLoggerAsInterface(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{Logger: log.New(&buf, "", 0)}

	// Use logger as service.Logger interface
	var serviceLogger service.Logger = logger

	// Test Info
	serviceLogger.Info("interface test", service.Fields{
		"test": "value",
	})

	output := buf.String()
	var entry LogEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got: %s", entry.Level)
	}
	if entry.Message != "interface test" {
		t.Errorf("Expected message 'interface test', got: %s", entry.Message)
	}
	if entry.Fields["test"] != "value" {
		t.Errorf("Expected field test=value, got: %v", entry.Fields["test"])
	}
}

// TestLoggerWithNilFields tests that logger handles nil fields gracefully
func TestLoggerWithNilFields(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{Logger: log.New(&buf, "", 0)}

	// Should not panic with nil fields
	logger.Info("test", nil)
	logger.Warn("test", nil)
	logger.Error("test", nil)

	// Verify output is valid JSON
	lines := bytes.Split(buf.Bytes(), []byte("\n"))
	for i, line := range lines {
		if len(line) == 0 {
			continue
		}
		var entry LogEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			t.Errorf("Line %d: Failed to parse log output: %v", i, err)
		}
	}
}
