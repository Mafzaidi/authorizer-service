package config

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		envKey   string
		envValue string
		fallback string
		want     string
	}{
		{
			name:     "returns env value when set",
			envKey:   "TEST_ENV_VAR",
			envValue: "test_value",
			fallback: "default_value",
			want:     "test_value",
		},
		{
			name:     "returns fallback when env not set",
			envKey:   "NONEXISTENT_ENV_VAR",
			envValue: "",
			fallback: "default_value",
			want:     "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			got := getEnvOrDefault(tt.envKey, tt.fallback)
			if got != tt.want {
				t.Errorf("getEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateKID(t *testing.T) {
	// Create a test private key
	privateKey, err := loadPrivateKeyFromEnvOrFile()
	if err != nil {
		t.Skip("Skipping test: private key not available")
	}

	kid := generateKID(&privateKey.PublicKey)
	
	// KID should not be empty
	if kid == "" {
		t.Error("generateKID() returned empty string")
	}

	// KID should be consistent for the same key
	kid2 := generateKID(&privateKey.PublicKey)
	if kid != kid2 {
		t.Errorf("generateKID() not consistent: %v != %v", kid, kid2)
	}
}
