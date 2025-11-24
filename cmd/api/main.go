package main

import (
	"log"

	"localdev.me/authorizer/config"
	"localdev.me/authorizer/internal/app"
	"localdev.me/authorizer/internal/infrastructure/persistence/postgres"
)

func main() {
	cfg := config.GetConfig()
	pool, err := postgres.NewPostgreSQL(cfg)
	if err != nil {
		panic(err)
	}

	s := app.NewServer(cfg, pool)
	if err := s.Run(); err != nil {
		log.Fatalf("Server could not be ran: %v", err)
	}
}
