package property

import (
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
)

// HTTPMethod generates random HTTP methods
func HTTPMethod() gopter.Gen {
	return gen.OneConstOf("GET", "POST", "PUT", "DELETE", "PATCH")
}

// Endpoint represents an API endpoint
type Endpoint struct {
	Method string
	Path   string
}

// CommonEndpoints returns a generator for common API endpoints
func CommonEndpoints() gopter.Gen {
	endpoints := []Endpoint{
		{Method: "POST", Path: "/api/v1/auth/login"},
		{Method: "POST", Path: "/api/v1/auth/logout"},
		{Method: "GET", Path: "/api/v1/auth/jwks"},
		{Method: "POST", Path: "/api/v1/users"},
		{Method: "GET", Path: "/api/v1/users"},
		{Method: "GET", Path: "/api/v1/users/:id"},
		{Method: "PUT", Path: "/api/v1/users/:id"},
		{Method: "DELETE", Path: "/api/v1/users/:id"},
		{Method: "POST", Path: "/api/v1/roles"},
		{Method: "GET", Path: "/api/v1/roles"},
		{Method: "GET", Path: "/api/v1/roles/:id"},
		{Method: "PUT", Path: "/api/v1/roles/:id"},
		{Method: "DELETE", Path: "/api/v1/roles/:id"},
		{Method: "POST", Path: "/api/v1/permissions"},
		{Method: "GET", Path: "/api/v1/permissions"},
		{Method: "GET", Path: "/api/v1/permissions/:id"},
		{Method: "PUT", Path: "/api/v1/permissions/:id"},
		{Method: "DELETE", Path: "/api/v1/permissions/:id"},
		{Method: "POST", Path: "/api/v1/applications"},
		{Method: "GET", Path: "/api/v1/applications"},
		{Method: "GET", Path: "/api/v1/applications/:id"},
		{Method: "PUT", Path: "/api/v1/applications/:id"},
		{Method: "DELETE", Path: "/api/v1/applications/:id"},
		{Method: "GET", Path: "/health"},
	}

	return gen.OneGenOf(
		gen.Const(endpoints[0]),
		gen.Const(endpoints[1]),
		gen.Const(endpoints[2]),
		gen.Const(endpoints[3]),
		gen.Const(endpoints[4]),
		gen.Const(endpoints[5]),
		gen.Const(endpoints[6]),
		gen.Const(endpoints[7]),
		gen.Const(endpoints[8]),
		gen.Const(endpoints[9]),
		gen.Const(endpoints[10]),
		gen.Const(endpoints[11]),
		gen.Const(endpoints[12]),
		gen.Const(endpoints[13]),
		gen.Const(endpoints[14]),
		gen.Const(endpoints[15]),
		gen.Const(endpoints[16]),
		gen.Const(endpoints[17]),
		gen.Const(endpoints[18]),
		gen.Const(endpoints[19]),
		gen.Const(endpoints[20]),
		gen.Const(endpoints[21]),
		gen.Const(endpoints[22]),
		gen.Const(endpoints[23]),
	)
}

// ValidEmail generates valid email addresses
func ValidEmail() gopter.Gen {
	return gen.Identifier().
		Map(func(name string) string {
			return name + "@example.com"
		})
}

// ValidPassword generates valid passwords
func ValidPassword() gopter.Gen {
	return gen.AlphaString().
		SuchThat(func(s string) bool {
			return len(s) >= 8
		})
}

// LoginRequest represents a login request
type LoginRequest struct {
	Application string
	Email       string
	Password    string
}

// ValidLoginRequest generates valid login requests
func ValidLoginRequest() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),
		ValidEmail(),
		ValidPassword(),
	).Map(func(values []interface{}) LoginRequest {
		return LoginRequest{
			Application: values[0].(string),
			Email:       values[1].(string),
			Password:    values[2].(string),
		}
	})
}

// UserRole represents a user role
type UserRole struct {
	UserID string
	RoleID string
}

// ValidUserRole generates valid user roles
func ValidUserRole() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),
		gen.Identifier(),
	).Map(func(values []interface{}) UserRole {
		return UserRole{
			UserID: values[0].(string),
			RoleID: values[1].(string),
		}
	})
}

// Permission represents a permission
type Permission struct {
	Code        string
	Name        string
	Description string
}

// ValidPermission generates valid permissions
func ValidPermission() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),
		gen.Identifier(),
		gen.Identifier(),
	).Map(func(values []interface{}) Permission {
		return Permission{
			Code:        values[0].(string),
			Name:        values[1].(string),
			Description: values[2].(string),
		}
	})
}

// CacheableOperation represents an operation that can be cached
type CacheableOperation struct {
	Type   string
	Key    string
	UserID string
}

// ValidCacheableOperation generates valid cacheable operations
func ValidCacheableOperation() gopter.Gen {
	return gopter.CombineGens(
		gen.OneConstOf("GetUser", "GetRole", "GetPermission", "GetApplication"),
		gen.Identifier(),
		gen.Identifier(),
	).Map(func(values []interface{}) CacheableOperation {
		return CacheableOperation{
			Type:   values[0].(string),
			Key:    values[1].(string),
			UserID: values[2].(string),
		}
	})
}

// TestUser represents a test user with roles and permissions
type TestUser struct {
	ID          string
	Email       string
	Password    string
	Roles       []string
	Permissions []string
}

// ValidTestUser generates valid test users
func ValidTestUser() gopter.Gen {
	return gopter.CombineGens(
		gen.Identifier(),
		ValidEmail(),
		ValidPassword(),
		gen.SliceOf(gen.Identifier()),
		gen.SliceOf(gen.Identifier()),
	).Map(func(values []interface{}) TestUser {
		return TestUser{
			ID:          values[0].(string),
			Email:       values[1].(string),
			Password:    values[2].(string),
			Roles:       values[3].([]string),
			Permissions: values[4].([]string),
		}
	})
}

// PropertyTestConfig holds configuration for property-based tests
type PropertyTestConfig struct {
	MinSuccessfulTests int
	MaxSize            int
	Workers            int
}

// DefaultPropertyTestConfig returns default configuration
func DefaultPropertyTestConfig() *PropertyTestConfig {
	return &PropertyTestConfig{
		MinSuccessfulTests: 100,
		MaxSize:            100,
		Workers:            1,
	}
}

// ToGopterParameters converts config to gopter parameters
func (c *PropertyTestConfig) ToGopterParameters() *gopter.TestParameters {
	return gopter.DefaultTestParameters()
}
