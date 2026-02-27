package logger

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"

	"github.com/mafzaidi/authorizer/internal/domain/service"
)

func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{Logger: log.New(&buf, "", 0)}

	logger.Info("test message", service.Fields{
		"key1": "value1",
		"key2": 123,
	})

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected log level INFO, got: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected message 'test message', got: %s", output)
	}

	// Verify JSON structure
	var entry LogEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse log output as JSON: %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got: %s", entry.Level)
	}
	if entry.Message != "test message" {
		t.Errorf("Expected message 'test message', got: %s", entry.Message)
	}
	if entry.Fields["key1"] != "value1" {
		t.Errorf("Expected field key1=value1, got: %v", entry.Fields["key1"])
	}
}

func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{Logger: log.New(&buf, "", 0)}

	logger.Error("error occurred", service.Fields{
		"error": "something went wrong",
	})

	output := buf.String()
	if !strings.Contains(output, "ERROR") {
		t.Errorf("Expected log level ERROR, got: %s", output)
	}
	if !strings.Contains(output, "error occurred") {
		t.Errorf("Expected message 'error occurred', got: %s", output)
	}
}

func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{Logger: log.New(&buf, "", 0)}

	logger.Warn("warning message", service.Fields{
		"warning": "be careful",
	})

	output := buf.String()
	if !strings.Contains(output, "WARN") {
		t.Errorf("Expected log level WARN, got: %s", output)
	}
	if !strings.Contains(output, "warning message") {
		t.Errorf("Expected message 'warning message', got: %s", output)
	}
}

func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{Logger: log.New(&buf, "", 0)}

	logger.Debug("debug info", service.Fields{
		"detail": "debugging",
	})

	output := buf.String()
	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Expected log level DEBUG, got: %s", output)
	}
	if !strings.Contains(output, "debug info") {
		t.Errorf("Expected message 'debug info', got: %s", output)
	}
}

func TestLogger_EmptyFields(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{Logger: log.New(&buf, "", 0)}

	logger.Info("message without fields", nil)

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected log level INFO, got: %s", output)
	}
	if !strings.Contains(output, "message without fields") {
		t.Errorf("Expected message 'message without fields', got: %s", output)
	}

	// Verify JSON structure
	var entry LogEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse log output as JSON: %v", err)
	}
}

func TestTruncateToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "normal token",
			token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0",
			expected: "...wIn0",
		},
		{
			name:     "short token",
			token:    "abc",
			expected: "",
		},
		{
			name:     "empty token",
			token:    "",
			expected: "",
		},
		{
			name:     "exactly 4 chars",
			token:    "abcd",
			expected: "...abcd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateToken(tt.token)
			if result != tt.expected {
				t.Errorf("TruncateToken(%q) = %q, want %q", tt.token, result, tt.expected)
			}
		})
	}
}

func TestNew(t *testing.T) {
	logger := New()
	if logger == nil {
		t.Fatal("New() returned nil")
	}
	if logger.Logger == nil {
		t.Fatal("New() returned logger with nil internal logger")
	}
}
