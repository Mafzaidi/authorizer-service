package config

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

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
		Secret        string
		TokenExpiry   time.Duration
		RefreshExpiry time.Duration
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
		viper.AddConfigPath(".")
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

		cfg.Redis.Host = getEnvOrDefault("REDIS_DB_HOST", cfg.Redis.Host)
		cfg.Redis.Port = getEnvOrDefault("REDIS_DB_PORT", cfg.Redis.Port)
		cfg.Redis.User = getEnvOrDefault("REDIS_USER", cfg.Redis.User)
		cfg.Redis.Password = getEnvOrDefault("REDIS_PASSWORD", cfg.Redis.Password)
		cfg.Redis.DBName = getEnvOrDefault("REDIS_DB_NAME", cfg.PostgresDB.DBName)

		cfg.JWT.Secret = getEnvOrDefault("JWT_SECRET", cfg.JWT.Secret)

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
