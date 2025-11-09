package main

import (
	"log"

	"localdev.me/authorizer/config"
	"localdev.me/authorizer/internal/infrastructure/database/postgres"
	"localdev.me/authorizer/internal/server"
)

func main() {
	cfg := config.GetConfig()
	pool, err := postgres.NewPostgreSQL(cfg)
	if err != nil {
		panic(err)
	}

	s := server.NewServer(cfg, pool)
	if err := s.Run(); err != nil {
		log.Fatalf("Server could not be ran: %v", err)
	}
}
