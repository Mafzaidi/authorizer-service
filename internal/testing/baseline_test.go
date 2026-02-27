package testing

import (
	"testing"

	"github.com/mafzaidi/authorizer/internal/testing/comparison"
	"github.com/mafzaidi/authorizer/internal/testing/property"
	"github.com/mafzaidi/authorizer/internal/testing/structural"
)

// TestBaselineInfrastructure verifies the testing infrastructure is set up correctly
func TestBaselineInfrastructure(t *testing.T) {
	t.Run("ComparisonUtilities", func(t *testing.T) {
		// Test that comparison utilities work
		resp := &comparison.HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: map[string]interface{}{
				"status": "success",
			},
		}

		// Test matchers
		if !comparison.StatusCodeMatcher(200)(resp) {
			t.Error("StatusCodeMatcher failed")
		}

		if !comparison.BodyFieldMatcher("status", "success")(resp) {
			t.Error("BodyFieldMatcher failed")
		}

		if !comparison.AllMatchers(
			comparison.StatusCodeMatcher(200),
			comparison.BodyFieldMatcher("status", "success"),
		)(resp) {
			t.Error("AllMatchers failed")
		}
	})

	t.Run("PropertyGenerators", func(t *testing.T) {
		// Test that property generators work
		config := property.DefaultPropertyTestConfig()
		if config.MinSuccessfulTests != 100 {
			t.Errorf("Expected MinSuccessfulTests to be 100, got %d", config.MinSuccessfulTests)
		}

		// Test email generator
		emailGen := property.ValidEmail()
		email, ok := emailGen.Sample()
		if !ok {
			t.Error("Failed to generate email")
		}
		if email == nil {
			t.Error("Generated email is nil")
		}

		// Test password generator
		passwordGen := property.ValidPassword()
		password, ok := passwordGen.Sample()
		if !ok {
			t.Error("Failed to generate password")
		}
		if password == nil {
			t.Error("Generated password is nil")
		}
	})

	t.Run("StructuralAnalyzer", func(t *testing.T) {
		// Test that structural analyzer can be created
		analyzer := structural.NewArchitectureAnalyzer("../../..")
		if analyzer == nil {
			t.Error("Failed to create architecture analyzer")
		}

		// Test folder structure verification
		missingDirs, err := analyzer.VerifyFolderStructure()
		if err != nil {
			t.Errorf("Error verifying folder structure: %v", err)
		}

		// We expect some directories might be missing initially
		t.Logf("Missing directories: %v", missingDirs)
	})
}

// TestStructuralAnalyzerFunctions tests individual analyzer functions
func TestStructuralAnalyzerFunctions(t *testing.T) {
	analyzer := structural.NewArchitectureAnalyzer("../../..")

	t.Run("VerifyFolderStructure", func(t *testing.T) {
		missingDirs, err := analyzer.VerifyFolderStructure()
		if err != nil {
			t.Fatalf("Error verifying folder structure: %v", err)
		}

		// Log the result for visibility
		if len(missingDirs) > 0 {
			t.Logf("Missing directories (expected during refactoring): %v", missingDirs)
		} else {
			t.Log("All required directories exist")
		}
	})

	t.Run("CheckResourcePackageImports", func(t *testing.T) {
		filesWithResourceImport, err := analyzer.CheckResourcePackageImports()
		if err != nil {
			t.Fatalf("Error checking resource package imports: %v", err)
		}

		// Log the result for visibility
		if len(filesWithResourceImport) > 0 {
			t.Logf("Files importing resource package (should be removed): %v", filesWithResourceImport)
		} else {
			t.Log("No files import resource package")
		}
	})
}

// TestComparisonUtilities tests the comparison utilities in detail
func TestComparisonUtilities(t *testing.T) {
	t.Run("CompareResponses_Identical", func(t *testing.T) {
		resp1 := &comparison.HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: map[string]interface{}{
				"status": "success",
				"data":   "test",
			},
		}

		resp2 := &comparison.HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: map[string]interface{}{
				"status": "success",
				"data":   "test",
			},
		}

		if !comparison.CompareResponses(t, resp1, resp2) {
			t.Error("Expected responses to be identical")
		}
	})

	t.Run("CompareResponses_DifferentStatusCode", func(t *testing.T) {
		resp1 := &comparison.HTTPResponse{
			StatusCode: 200,
			Headers:    map[string]string{},
			Body:       map[string]interface{}{},
		}

		resp2 := &comparison.HTTPResponse{
			StatusCode: 404,
			Headers:    map[string]string{},
			Body:       map[string]interface{}{},
		}

		// Create a sub-test to capture the error
		subT := &testing.T{}
		if comparison.CompareResponses(subT, resp1, resp2) {
			t.Error("Expected responses to be different")
		}
	})
}
