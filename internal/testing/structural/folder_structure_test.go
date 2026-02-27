package structural

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFolderStructureExists verifies all required directories exist
// **Validates: Requirements 4.3**
func TestFolderStructureExists(t *testing.T) {
	// Get the project root (3 levels up from internal/testing/structural)
	projectRoot, err := filepath.Abs(filepath.Join("..", "..", ".."))
	require.NoError(t, err, "Failed to get project root")

	// Define required directories according to clean architecture
	requiredDirs := []string{
		// Domain layer
		"internal/domain/entity",
		"internal/domain/repository",
		"internal/domain/service",
		
		// Usecase layer
		"internal/usecase",
		
		// Infrastructure layer
		"internal/infrastructure/config",
		"internal/infrastructure/logger",
		"internal/infrastructure/auth",
		"internal/infrastructure/persistence",
		
		// Delivery layer
		"internal/delivery/http/handler",
		"internal/delivery/http/middleware",
		"internal/delivery/http/router",
	}

	// Verify each required directory exists
	for _, dir := range requiredDirs {
		fullPath := filepath.Join(projectRoot, dir)
		info, err := os.Stat(fullPath)
		
		assert.NoError(t, err, "Directory should exist: %s", dir)
		if err == nil {
			assert.True(t, info.IsDir(), "Path should be a directory: %s", dir)
		}
	}
}

// TestCleanArchitectureLayering verifies the structure follows clean architecture pattern
// **Validates: Requirements 4.3**
func TestCleanArchitectureLayering(t *testing.T) {
	// Get the project root (3 levels up from internal/testing/structural)
	projectRoot, err := filepath.Abs(filepath.Join("..", "..", ".."))
	require.NoError(t, err, "Failed to get project root")

	// Verify layer hierarchy exists
	layers := []struct {
		name string
		path string
	}{
		{"Domain", "internal/domain"},
		{"UseCase", "internal/usecase"},
		{"Infrastructure", "internal/infrastructure"},
		{"Delivery", "internal/delivery"},
	}

	for _, layer := range layers {
		fullPath := filepath.Join(projectRoot, layer.path)
		info, err := os.Stat(fullPath)
		
		assert.NoError(t, err, "%s layer should exist at: %s", layer.name, layer.path)
		if err == nil {
			assert.True(t, info.IsDir(), "%s layer path should be a directory", layer.name)
		}
	}

	// Verify domain sublayers
	domainSubLayers := []string{
		"internal/domain/entity",
		"internal/domain/repository",
		"internal/domain/service",
	}

	for _, subLayer := range domainSubLayers {
		fullPath := filepath.Join(projectRoot, subLayer)
		info, err := os.Stat(fullPath)
		
		assert.NoError(t, err, "Domain sublayer should exist: %s", subLayer)
		if err == nil {
			assert.True(t, info.IsDir(), "Domain sublayer should be a directory: %s", subLayer)
		}
	}

	// Verify infrastructure sublayers
	infraSubLayers := []string{
		"internal/infrastructure/config",
		"internal/infrastructure/logger",
		"internal/infrastructure/auth",
		"internal/infrastructure/persistence",
	}

	for _, subLayer := range infraSubLayers {
		fullPath := filepath.Join(projectRoot, subLayer)
		info, err := os.Stat(fullPath)
		
		assert.NoError(t, err, "Infrastructure sublayer should exist: %s", subLayer)
		if err == nil {
			assert.True(t, info.IsDir(), "Infrastructure sublayer should be a directory: %s", subLayer)
		}
	}

	// Verify delivery sublayers
	deliverySubLayers := []string{
		"internal/delivery/http/handler",
		"internal/delivery/http/middleware",
		"internal/delivery/http/router",
	}

	for _, subLayer := range deliverySubLayers {
		fullPath := filepath.Join(projectRoot, subLayer)
		info, err := os.Stat(fullPath)
		
		assert.NoError(t, err, "Delivery sublayer should exist: %s", subLayer)
		if err == nil {
			assert.True(t, info.IsDir(), "Delivery sublayer should be a directory: %s", subLayer)
		}
	}
}
