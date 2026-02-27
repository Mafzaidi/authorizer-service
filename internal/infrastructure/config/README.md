# Infrastructure Config Package

This package provides configuration loading functionality for the Authorizer application, following clean architecture principles.

## Features

- Load configuration from YAML files and environment variables
- Support for structured logging via logger interface
- RSA key pair loading for JWT operations
- Backward compatible with the old `config` package

## Usage

### New Way (Recommended)

```go
import (
    "github.com/mafzaidi/authorizer/internal/infrastructure/config"
    "github.com/mafzaidi/authorizer/internal/infrastructure/logger"
)

func main() {
    // Initialize logger first
    log := logger.New()
    
    // Load config with logger
    cfg, err := config.Load(log)
    if err != nil {
        log.Error("Failed to load config", logger.Fields{"error": err.Error()})
        return
    }
    
    // Use config
    fmt.Printf("Server port: %d\n", cfg.Server.Port)
}
```

### Old Way (Still Supported)

```go
import "github.com/mafzaidi/authorizer/config"

func main() {
    cfg := config.GetConfig()
    fmt.Printf("Server port: %d\n", cfg.Server.Port)
}
```

## Configuration Structure

The configuration is loaded from:
1. YAML file (default: `/app/config/config.yaml`)
2. Environment variables (override YAML values)
3. `.env` file (if present)

### Environment Variables

- `CONFIG_PATH`: Path to config directory (default: `/app/config`)
- `POSTGRES_DB_HOST`: PostgreSQL host
- `POSTGRES_DB_PORT`: PostgreSQL port
- `POSTGRES_USER`: PostgreSQL user
- `POSTGRES_PASSWORD`: PostgreSQL password
- `POSTGRES_DB_NAME`: PostgreSQL database name
- `REDIS_HOST`: Redis host
- `REDIS_PORT`: Redis port
- `JWT_PRIVATE_KEY`: RSA private key in PEM format (as string)
- `JWT_PRIVATE_KEY_PATH`: Path to RSA private key file

## Migration Guide

The new config package is designed to be used in the refactored clean architecture. During the migration:

1. **Phase 1-7**: Old config package (`github.com/mafzaidi/authorizer/config`) remains in use
2. **Phase 8**: New code uses the infrastructure config package
3. **Phase 9**: Old config package can be removed after all references are updated

## Testing

Run tests with:

```bash
go test ./internal/infrastructure/config/
```

Note: Some tests require a valid private key file to be present.
