package config

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Default paths for JWT private key (public key is derived from it).
const defaultPrivateKeyPath = "./private.pem"

type (
	Config struct {
		Server     *Server
		App        *App
		PostgresDB *PostgresDB
		Redis      *Redis
		JWT        *JWT
	}

	App struct {
		Name    string
		Version string
	}

	Server struct {
		Host string
		Port int
	}

	PostgresDB struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}

	Redis struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}

	JWT struct {
		PrivateKeyPath string
		PublicKeyPath  string
		PrivateKey     *rsa.PrivateKey
		PublicKey      *rsa.PublicKey
		KeyID          string
		TokenExpiry    time.Duration
		RefreshExpiry  time.Duration
	}
)

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig() *Config {
	once.Do(func() {

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(getEnvOrDefault("CONFIG_PATH", "/app/config"))
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic("Failed to read config.yaml: " + err.Error())
		}

		cfg := &Config{
			Server:     &Server{},
			App:        &App{},
			PostgresDB: &PostgresDB{},
			Redis:      &Redis{},
			JWT:        &JWT{},
		}

		if err := viper.Unmarshal(cfg); err != nil {
			panic("Failed to unmarshal config into struct: " + err.Error())
		}

		_ = godotenv.Load()

		cfg.PostgresDB.Host = getEnvOrDefault("POSTGRES_DB_HOST", cfg.PostgresDB.Host)
		cfg.PostgresDB.Port = getEnvOrDefault("POSTGRES_DB_PORT", cfg.PostgresDB.Port)
		cfg.PostgresDB.User = getEnvOrDefault("POSTGRES_USER", cfg.PostgresDB.User)
		cfg.PostgresDB.Password = getEnvOrDefault("POSTGRES_PASSWORD", cfg.PostgresDB.Password)
		cfg.PostgresDB.DBName = getEnvOrDefault("POSTGRES_DB_NAME", cfg.PostgresDB.DBName)

		cfg.Redis.Host = getEnvOrDefault("REDIS_HOST", cfg.Redis.Host)
		cfg.Redis.Port = getEnvOrDefault("REDIS_PORT", cfg.Redis.Port)

		// Load private key only (from env PEM or file). Public key is derived from it
		// so only one secret (private.pem) is needed; JWKS endpoint serves the public key.
		privateKey, err := loadPrivateKeyFromEnvOrFile()
		if err != nil {
			panic(fmt.Sprintf("Failed to load private key: %v", err))
		}
		cfg.JWT.PrivateKey = privateKey
		cfg.JWT.PublicKey = &privateKey.PublicKey

		cfg.JWT.KeyID = generateKID(cfg.JWT.PublicKey)

		if s := viper.GetString("jwt.tokenExpiry"); s != "" {
			cfg.JWT.TokenExpiry, _ = time.ParseDuration(s)
		}
		if s := viper.GetString("jwt.refreshExpiry"); s != "" {
			cfg.JWT.RefreshExpiry, _ = time.ParseDuration(s)
		}

		configInstance = cfg
	})

	return configInstance
}

func getEnvOrDefault(envKey, fallback string) string {
	if val := os.Getenv(envKey); val != "" {
		return val
	}
	return fallback
}

// loadPrivateKeyFromEnvOrFile loads private key from JWT_PRIVATE_KEY (PEM string) or from file.
// File paths tried in order: JWT_PRIVATE_KEY_PATH env, default ./private.pem, then /app/secrets/private.pem (K8s mount).
func loadPrivateKeyFromEnvOrFile() (*rsa.PrivateKey, error) {
	if pemStr := os.Getenv("JWT_PRIVATE_KEY"); pemStr != "" {
		return parsePrivateKeyPEM([]byte(pemStr))
	}

	paths := []string{
		getEnvOrDefault("JWT_PRIVATE_KEY_PATH", defaultPrivateKeyPath),
		defaultPrivateKeyPath,
	}
	for _, path := range paths {
		key, err := loadPrivateKey(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || strings.Contains(err.Error(), "no such file") {
				continue
			}
			return nil, err
		}
		return key, nil
	}
	return nil, fmt.Errorf("private key not found: tried paths %v and JWT_PRIVATE_KEY env", paths)
}

func parsePrivateKeyPEM(keyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	}
	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA private key")
	}
	return rsaKey, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}
	return parsePrivateKeyPEM(keyBytes)
}

func generateKID(pub *rsa.PublicKey) string {
	hash := sha256.Sum256(pub.N.Bytes())
	return base64.RawURLEncoding.EncodeToString(hash[:8])
}
