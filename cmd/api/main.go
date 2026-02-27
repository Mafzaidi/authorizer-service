package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/internal/delivery/http/handler"
	"github.com/mafzaidi/authorizer/internal/infrastructure/config"
	"github.com/mafzaidi/authorizer/internal/delivery/http/middleware"
	"github.com/mafzaidi/authorizer/internal/delivery/http/router"
	"github.com/mafzaidi/authorizer/internal/domain/service"
	"github.com/mafzaidi/authorizer/internal/infrastructure/auth"
	infraConfig "github.com/mafzaidi/authorizer/internal/infrastructure/config"
	"github.com/mafzaidi/authorizer/internal/infrastructure/logger"
	"github.com/mafzaidi/authorizer/internal/infrastructure/persistence/postgres"
	postgresRepo "github.com/mafzaidi/authorizer/internal/infrastructure/persistence/postgres/repository"
	"github.com/mafzaidi/authorizer/internal/infrastructure/persistence/redis"
	redisRepo "github.com/mafzaidi/authorizer/internal/infrastructure/persistence/redis/repository"
	appUsecase "github.com/mafzaidi/authorizer/internal/usecase/application"
	authUsecase "github.com/mafzaidi/authorizer/internal/usecase/auth"
	permUsecase "github.com/mafzaidi/authorizer/internal/usecase/permission"
	roleUsecase "github.com/mafzaidi/authorizer/internal/usecase/role"
	userUsecase "github.com/mafzaidi/authorizer/internal/usecase/user"
)

func main() {
	// 1. Initialize logger first (needed for config loading)
	log := logger.New()
	log.Info("Starting Authorizer application", logger.Fields{})

	// 2. Load configuration
	cfg := config.GetConfig()
	log.Info("Configuration loaded successfully", logger.Fields{})

	// 3. Initialize database connection
	pgConn, err := postgres.NewPostgreSQL(cfg)
	if err != nil {
		log.Error("Failed to connect to PostgreSQL", logger.Fields{
			"error": err.Error(),
		})
		panic(fmt.Sprintf("Failed to connect to PostgreSQL: %v", err))
	}
	pool := pgConn.Pool
	log.Info("PostgreSQL connection established", logger.Fields{})

	// 4. Initialize Redis connection
	redisConn := redis.NewRedisClient(cfg)
	redisClient := redisConn.Client
	log.Info("Redis connection established", logger.Fields{})

	// 5. Initialize repositories
	// PostgreSQL repositories
	userRepo := postgresRepo.NewUserRepositoryPGX(pool)
	roleRepo := postgresRepo.NewRoleRepositoryPGX(pool)
	permRepo := postgresRepo.NewPermRepositoryPGX(pool)
	appRepo := postgresRepo.NewAppRepositoryPGX(pool)
	userRoleRepo := postgresRepo.NewUserRoleRepositoryPGX(pool)
	rolePermRepo := postgresRepo.NewRolePermRepositoryPGX(pool)

	// Redis repositories
	authRepo := redisRepo.NewAuthRepository(redisClient)

	log.Info("All repositories initialized", logger.Fields{})

	// 6. Initialize domain services
	authService := service.NewAuthService(
		userRoleRepo,
		roleRepo,
		rolePermRepo,
		appRepo,
	)
	log.Info("Domain services initialized", logger.Fields{})

	// 7. Initialize infrastructure services
	jwtService := auth.NewJWTService(log)
	jwksService := auth.NewJWKSService()
	log.Info("Infrastructure services initialized", logger.Fields{})

	// 8. Initialize use cases
	authUC := authUsecase.NewAuthUseCase(
		authRepo,
		userRepo,
		authService,
		jwtService,
		log,
	)

	userUC := userUsecase.NewUserUsecase(
		userRepo,
		roleRepo,
		userRoleRepo,
		log,
	)

	roleUC := roleUsecase.NewRoleUsecase(
		roleRepo,
		appRepo,
		permRepo,
		rolePermRepo,
		log,
	)

	permUC := permUsecase.NewPermUsecase(
		permRepo,
		appRepo,
		log,
	)

	appUC := appUsecase.NewAppUsecase(
		appRepo,
		log,
	)

	log.Info("All use cases initialized", logger.Fields{})

	// 9. Initialize handlers
	authHandler := handler.NewAuthHandler(
		authUC,
		jwksService,
		cfg,
		log,
	)

	userHandler := handler.NewUserHandler(
		userUC,
		log,
	)

	roleHandler := handler.NewRoleHandler(
		roleUC,
		log,
	)

	permHandler := handler.NewPermHandler(
		permUC,
		log,
	)

	appHandler := handler.NewAppHandler(
		appUC,
		log,
	)

	healthHandler := handler.NewHealthHandler(log)

	log.Info("All handlers initialized", logger.Fields{})

	// 10. Initialize middleware
	// Note: Middleware uses new infrastructure config, but we need to convert from old config
	// This will be cleaned up when handlers are fully migrated to new config
	jwtMiddleware := middleware.JWTAuthMiddleware(jwtService, convertToInfraConfig(cfg), log)
	log.Info("Middleware initialized", logger.Fields{})

	// 11. Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// 12. Setup router with all dependencies
	err = router.Setup(e, &router.RouterConfig{
		AuthHandler:   authHandler,
		UserHandler:   userHandler,
		RoleHandler:   roleHandler,
		PermHandler:   permHandler,
		AppHandler:    appHandler,
		HealthHandler: healthHandler,
		JWTMiddleware: jwtMiddleware,
		Logger:        log,
	})
	if err != nil {
		log.Error("Failed to setup router", logger.Fields{
			"error": err.Error(),
		})
		panic(fmt.Sprintf("Failed to setup router: %v", err))
	}
	log.Info("Router configured successfully", logger.Fields{})

	// 13. Start server with graceful shutdown
	startServer(e, cfg, log)
}

// startServer starts the HTTP server with graceful shutdown support
func startServer(e *echo.Echo, cfg *config.Config, log *logger.Logger) {
	// Server address
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// Start server in a goroutine
	go func() {
		log.Info("Starting HTTP server", logger.Fields{
			"host": cfg.Server.Host,
			"port": cfg.Server.Port,
			"addr": addr,
		})

		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed to start", logger.Fields{
				"error": err.Error(),
			})
			panic(fmt.Sprintf("Server failed to start: %v", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...", logger.Fields{})

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", logger.Fields{
			"error": err.Error(),
		})
	}

	log.Info("Server exited gracefully", logger.Fields{})
}

// convertToInfraConfig converts old config to new infrastructure config
// This is a temporary bridge until all components use the new config
func convertToInfraConfig(oldCfg *config.Config) *infraConfig.Config {
	return &infraConfig.Config{
		Server: &infraConfig.Server{
			Host: oldCfg.Server.Host,
			Port: oldCfg.Server.Port,
		},
		App: &infraConfig.App{
			Name:    oldCfg.App.Name,
			Version: oldCfg.App.Version,
		},
		PostgresDB: &infraConfig.PostgresDB{
			Host:     oldCfg.PostgresDB.Host,
			Port:     oldCfg.PostgresDB.Port,
			User:     oldCfg.PostgresDB.User,
			Password: oldCfg.PostgresDB.Password,
			DBName:   oldCfg.PostgresDB.DBName,
		},
		Redis: &infraConfig.Redis{
			Host:     oldCfg.Redis.Host,
			Port:     oldCfg.Redis.Port,
			User:     oldCfg.Redis.User,
			Password: oldCfg.Redis.Password,
			DBName:   oldCfg.Redis.DBName,
		},
		JWT: &infraConfig.JWT{
			PrivateKeyPath: oldCfg.JWT.PrivateKeyPath,
			PublicKeyPath:  oldCfg.JWT.PublicKeyPath,
			PrivateKey:     oldCfg.JWT.PrivateKey,
			PublicKey:      oldCfg.JWT.PublicKey,
			KeyID:          oldCfg.JWT.KeyID,
			TokenExpiry:    oldCfg.JWT.TokenExpiry,
			RefreshExpiry:  oldCfg.JWT.RefreshExpiry,
		},
	}
}
