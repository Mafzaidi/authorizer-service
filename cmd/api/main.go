package main

import (
	"log"

	"github.com/mafzaidi/authorizer/config"
	"github.com/mafzaidi/authorizer/internal/app"
	"github.com/mafzaidi/authorizer/internal/infrastructure/persistence/postgres"
	"github.com/mafzaidi/authorizer/internal/infrastructure/persistence/redis"
)

func main() {
	cfg := config.GetConfig()
	pool, err := postgres.NewPostgreSQL(cfg)
	redis := redis.NewRedisClient(cfg)
	if err != nil {
		panic(err)
	}

	s := app.NewServer(cfg, pool, redis)
	if err := s.Run(); err != nil {
		log.Fatalf("Server could not be ran: %v", err)
	}
}
