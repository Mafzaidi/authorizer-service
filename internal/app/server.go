package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mafzaidi/authorizer/config"
	"github.com/mafzaidi/authorizer/internal/infrastructure/persistence/postgres"
	"github.com/mafzaidi/authorizer/internal/infrastructure/persistence/redis"
)

type Server struct {
	echo      *echo.Echo
	cfg       *config.Config
	postgreDB *postgres.PostgreSQL
	redis     *redis.Redis
}

func NewServer(cfg *config.Config, postgreDB *postgres.PostgreSQL, redis *redis.Redis) *Server {
	return &Server{
		echo:      echo.New(),
		cfg:       cfg,
		postgreDB: postgreDB,
		redis:     redis,
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
	server := &http.Server{
		Addr: addr,
	}

	go func() {
		if err := s.echo.StartServer(server); err != nil {
			log.Fatalf("Could not start the server: %v", err)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		log.Fatalf("An error has occured mapping the handlers: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	<-quit
	log.Println("Server is stopping...")

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.echo.Shutdown(ctx)
}
