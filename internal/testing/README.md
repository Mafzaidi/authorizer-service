# Testing Infrastructure for Clean Architecture Refactoring

This directory contains the baseline testing infrastructure for verifying the clean architecture refactoring of the Authorizer project.

## Overview

The testing infrastructure consists of three main components:

1. **Comparison Utilities** - For comparing old vs new behavior
2. **Property-Based Testing** - Using gopter for universal property verification
3. **Structural Analysis** - For verifying architecture compliance

## Components

### 1. Comparison Utilities (`comparison/`)

Located in `comparison/behavior_test_utils.go`, these utilities help compare HTTP responses before and after refactoring.

**Key Features:**
- Execute HTTP requests against Echo instances
- Compare responses (status codes, headers, body structure)
- Flexible matchers for response validation

**Example Usage:**
```go
import "github.com/mafzaidi/authorizer/internal/testing/comparison"

// Execute a request
req := comparison.HTTPRequest{
    Method: "POST",
    Path: "/api/v1/auth/login",
    Body: map[string]interface{}{
        "email": "test@example.com",
        "password": "password123",
    },
}

oldResp, _ := comparison.ExecuteRequest(oldEchoInstance, req)
newResp, _ := comparison.ExecuteRequest(newEchoInstance, req)

// Compare responses
if !comparison.CompareResponses(t, oldResp, newResp) {
    t.Error("Responses differ after refactoring")
}
```

### 2. Property-Based Testing (`property/`)

Located in `property/generators.go`, these utilities provide generators for property-based testing using gopter.

**Key Features:**
- Generators for common data types (emails, passwords, endpoints)
- Generators for domain entities (users, roles, permissions)
- Configurable test parameters

**Example Usage:**
```go
import (
    "github.com/leanovate/gopter"
    "github.com/leanovate/gopter/prop"
    "github.com/mafzaidi/authorizer/internal/testing/property"
)

func TestProperty_APIBackwardCompatibility(t *testing.T) {
    properties := gopter.NewProperties(property.DefaultPropertyTestConfig().ToGopterParameters())
    
    properties.Property("API responses identical before and after refactoring",
        prop.ForAll(
            func(endpoint property.Endpoint, req property.LoginRequest) bool {
                // Test logic here
                return true
            },
            property.CommonEndpoints(),
            property.ValidLoginRequest(),
        ))
    
    properties.TestingRun(t)
}
```

### 3. Structural Analysis (`structural/`)

Located in `structural/architecture_analyzer.go`, these tools verify architecture compliance.

**Key Features:**
- Check layer dependency rules
- Verify no resource package imports
- Check dependency injection patterns
- Verify folder structure

**Example Usage:**
```go
import "github.com/mafzaidi/authorizer/internal/testing/structural"

func TestArchitectureCompliance(t *testing.T) {
    analyzer := structural.NewArchitectureAnalyzer("../../..")
    
    // Check layer dependencies
    violations, err := analyzer.CheckLayerDependencies()
    if err != nil {
        t.Fatal(err)
    }
    
    if len(violations) > 0 {
        for _, v := range violations {
            t.Errorf("Violation: %s in layer %s imports %s from layer %s",
                v.File, v.Layer, v.ImportPath, v.ImportedLayer)
        }
    }
    
    // Check resource package imports
    files, err := analyzer.CheckResourcePackageImports()
    if err != nil {
        t.Fatal(err)
    }
    
    if len(files) > 0 {
        t.Errorf("Files still importing resource package: %v", files)
    }
}
```

## Running Tests

### Run all baseline tests:
```bash
go test ./internal/testing/...
```

### Run with verbose output:
```bash
go test -v ./internal/testing/...
```

### Run specific test:
```bash
go test -v ./internal/testing/ -run TestBaselineInfrastructure
```

## Property-Based Testing Configuration

Default configuration:
- **MinSuccessfulTests**: 100 iterations per property
- **MaxSize**: 100 (maximum size for generated values)
- **Workers**: 1 (sequential execution)

Customize by creating a new config:
```go
config := &property.PropertyTestConfig{
    MinSuccessfulTests: 200,
    MaxSize: 150,
    Workers: 4,
}
```

## Architecture Rules

The structural analyzer enforces these clean architecture rules:

1. **Domain Layer**: No dependencies on other layers
2. **UseCase Layer**: Can only depend on Domain layer
3. **Infrastructure Layer**: Can only depend on Domain layer
4. **Delivery Layer**: Can depend on Domain and UseCase layers

## Validation Requirements

According to the design document, the following properties must be validated:

### Behavioral Properties
- **Property 1**: Complete API Backward Compatibility (Requirements 6.1, 6.2, 9.1-9.5)
- **Property 2**: Authentication and Authorization Preservation (Requirement 6.3)
- **Property 3**: Redis Caching Behavior Preservation (Requirement 6.5)

### Structural Properties
- **Property 4**: Layer Dependency Rule Compliance (Requirements 4.4, 4.5)
- **Property 5**: Dependency Injection Pattern Compliance (Requirements 5.1-5.3)
- **Property 6**: Handler Structure Compliance (Requirements 7.2-7.4)
- **Property 7**: Middleware Dependency Injection Compliance (Requirements 8.1-8.4)

### Example-Based Verification
- **Example 1**: No Resource Package Imports (Requirement 1.4)
- **Example 2**: No Circular Dependencies (Requirement 2.6)
- **Example 3**: Folder Structure Exists (Requirement 4.3)
- **Example 4**: Test Coverage Per Layer (Requirement 5.5)
- **Example 5**: Existing Integration Tests Pass (Requirement 6.6)

## Next Steps

After setting up this baseline infrastructure, the refactoring can proceed with confidence that:

1. Behavioral changes will be detected by comparison tests
2. Universal properties will be verified across all inputs
3. Architecture violations will be caught early
4. The refactoring maintains backward compatibility

## References

- **gopter**: https://github.com/leanovate/gopter
- **Clean Architecture**: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- **Design Document**: `authorizer/.kiro/specs/clean-architecture-refactoring/design.md`
