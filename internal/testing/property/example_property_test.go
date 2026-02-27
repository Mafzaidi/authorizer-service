package property

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

// This file contains example property-based tests that demonstrate how to use
// the testing infrastructure for validating the clean architecture refactoring.

// TestExample_EmailGeneration demonstrates email generator
func TestExample_EmailGeneration(t *testing.T) {
	properties := gopter.NewProperties(DefaultPropertyTestConfig().ToGopterParameters())

	properties.Property("Generated emails are valid",
		prop.ForAll(
			func(email string) bool {
				// All generated emails should contain @
				return len(email) > 0 && email[len(email)-12:] == "@example.com"
			},
			ValidEmail(),
		))

	properties.TestingRun(t)
}

// TestExample_PasswordGeneration demonstrates password generator
func TestExample_PasswordGeneration(t *testing.T) {
	properties := gopter.NewProperties(DefaultPropertyTestConfig().ToGopterParameters())

	properties.Property("Generated passwords meet minimum length",
		prop.ForAll(
			func(password string) bool {
				// All generated passwords should be at least 8 characters
				return len(password) >= 8
			},
			ValidPassword(),
		))

	properties.TestingRun(t)
}

// TestExample_LoginRequestGeneration demonstrates login request generator
func TestExample_LoginRequestGeneration(t *testing.T) {
	properties := gopter.NewProperties(DefaultPropertyTestConfig().ToGopterParameters())

	properties.Property("Generated login requests have all required fields",
		prop.ForAll(
			func(req LoginRequest) bool {
				// All fields should be non-empty
				return req.Application != "" &&
					req.Email != "" &&
					req.Password != "" &&
					len(req.Password) >= 8
			},
			ValidLoginRequest(),
		))

	properties.TestingRun(t)
}

// TestExample_EndpointGeneration demonstrates endpoint generator
func TestExample_EndpointGeneration(t *testing.T) {
	properties := gopter.NewProperties(DefaultPropertyTestConfig().ToGopterParameters())

	properties.Property("Generated endpoints have valid HTTP methods",
		prop.ForAll(
			func(endpoint Endpoint) bool {
				validMethods := map[string]bool{
					"GET":    true,
					"POST":   true,
					"PUT":    true,
					"DELETE": true,
					"PATCH":  true,
				}
				return validMethods[endpoint.Method] && endpoint.Path != ""
			},
			CommonEndpoints(),
		))

	properties.TestingRun(t)
}

// TestExample_UserRoleGeneration demonstrates user role generator
func TestExample_UserRoleGeneration(t *testing.T) {
	properties := gopter.NewProperties(DefaultPropertyTestConfig().ToGopterParameters())

	properties.Property("Generated user roles have valid IDs",
		prop.ForAll(
			func(userRole UserRole) bool {
				return userRole.UserID != "" && userRole.RoleID != ""
			},
			ValidUserRole(),
		))

	properties.TestingRun(t)
}

// TestExample_PermissionGeneration demonstrates permission generator
func TestExample_PermissionGeneration(t *testing.T) {
	properties := gopter.NewProperties(DefaultPropertyTestConfig().ToGopterParameters())

	properties.Property("Generated permissions have all required fields",
		prop.ForAll(
			func(perm Permission) bool {
				return perm.Code != "" && perm.Name != "" && perm.Description != ""
			},
			ValidPermission(),
		))

	properties.TestingRun(t)
}

// TestExample_CacheableOperationGeneration demonstrates cacheable operation generator
func TestExample_CacheableOperationGeneration(t *testing.T) {
	properties := gopter.NewProperties(DefaultPropertyTestConfig().ToGopterParameters())

	properties.Property("Generated cacheable operations have valid types",
		prop.ForAll(
			func(op CacheableOperation) bool {
				validTypes := map[string]bool{
					"GetUser":        true,
					"GetRole":        true,
					"GetPermission":  true,
					"GetApplication": true,
				}
				return validTypes[op.Type] && op.Key != "" && op.UserID != ""
			},
			ValidCacheableOperation(),
		))

	properties.TestingRun(t)
}

// TestExample_TestUserGeneration demonstrates test user generator
func TestExample_TestUserGeneration(t *testing.T) {
	properties := gopter.NewProperties(DefaultPropertyTestConfig().ToGopterParameters())

	properties.Property("Generated test users have valid credentials",
		prop.ForAll(
			func(user TestUser) bool {
				return user.ID != "" &&
					user.Email != "" &&
					len(user.Password) >= 8 &&
					user.Roles != nil &&
					user.Permissions != nil
			},
			ValidTestUser(),
		))

	properties.TestingRun(t)
}
