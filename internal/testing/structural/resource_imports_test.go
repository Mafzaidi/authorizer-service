package structural

import (
	"testing"
)

// TestNoResourcePackageImports verifies that no files import internal/app/resource
// This validates Requirement 1.4: The System SHALL ensure no import statement references internal/app/resource
func TestNoResourcePackageImports(t *testing.T) {
	analyzer := NewArchitectureAnalyzer("../../..")

	filesWithResourceImport, err := analyzer.CheckResourcePackageImports()
	if err != nil {
		t.Fatalf("Failed to check resource package imports: %v", err)
	}

	if len(filesWithResourceImport) > 0 {
		t.Errorf("Found %d file(s) importing internal/app/resource package:", len(filesWithResourceImport))
		for _, file := range filesWithResourceImport {
			t.Errorf("  - %s", file)
		}
		t.Fatal("Resource package imports must be removed before completing this task")
	}

	t.Log("✓ No files import internal/app/resource package")
}
