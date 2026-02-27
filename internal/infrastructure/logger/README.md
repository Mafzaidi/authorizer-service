# Logger Package

This package provides a structured logging implementation for the Authorizer application.

## Overview

The logger implements the `service.Logger` interface defined in the domain layer, providing structured JSON logging with multiple log levels.

## Features

- **Structured Logging**: All logs are output in JSON format with timestamp, level, message, and fields
- **Multiple Log Levels**: Supports INFO, WARN, ERROR, and DEBUG levels
- **Type-Safe Fields**: Uses the `service.Fields` type for structured log fields
- **Token Truncation**: Utility function for safely logging sensitive tokens

## Usage

### Creating a Logger

```go
import "github.com/mafzaidi/authorizer/internal/infrastructure/logger"

// Create a new logger instance
log := logger.New()
```

### Logging Messages

```go
// Info level
log.Info("User logged in", service.Fields{
    "user_id": "123",
    "email": "user@example.com",
})

// Warning level
log.Warn("Rate limit approaching", service.Fields{
    "user_id": "123",
    "requests": 95,
    "limit": 100,
})

// Error level
log.Error("Database connection failed", service.Fields{
    "error": err.Error(),
    "host": "localhost",
})

// Debug level (for development)
log.Debug("Processing request", service.Fields{
    "request_id": "abc-123",
    "path": "/api/users",
})
```

### Log Output Format

All logs are output as JSON:

```json
{
  "timestamp": "2024-01-15T10:30:45Z",
  "level": "INFO",
  "message": "User logged in",
  "fields": {
    "user_id": "123",
    "email": "user@example.com"
  }
}
```

### Token Truncation

For safely logging tokens or other sensitive data:

```go
import "github.com/mafzaidi/authorizer/internal/infrastructure/logger"

token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0"
truncated := logger.TruncateToken(token)
// Returns: "...wIn0"

log.Info("Token validated", service.Fields{
    "token": truncated,
})
```

## Implementation Details

- Logs are written to `os.Stdout`
- Timestamps are in UTC using RFC3339 format
- If JSON marshaling fails, falls back to simple text logging
- The logger is safe for concurrent use (backed by Go's standard `log.Logger`)

## Interface Compliance

This logger implements the `service.Logger` interface:

```go
type Logger interface {
    Info(message string, fields Fields)
    Warn(message string, fields Fields)
    Error(message string, fields Fields)
}
```

The `Debug` method is provided as an additional utility but is not part of the interface.

## Testing

Run tests with:

```bash
go test ./internal/infrastructure/logger/... -v
```

## Dependencies

- Standard library only (`encoding/json`, `log`, `os`, `time`)
- Domain service package for the `Logger` interface and `Fields` type
