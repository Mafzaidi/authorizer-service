package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mafzaidi/authorizer/config"
)

type PostgreSQL struct {
	Pool *pgxpool.Pool
}

var (
	once       sync.Once
	dbInstance *PostgreSQL
	initErr    error
)

func NewPostgreSQL(conf *config.Config) (*PostgreSQL, error) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			conf.PostgresDB.User,
			conf.PostgresDB.Password,
			conf.PostgresDB.Host,
			conf.PostgresDB.Port,
			conf.PostgresDB.DBName,
		)

		cfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			initErr = fmt.Errorf("failed to parse DSN: %w", err)
			return
		}

		cfg.MaxConns = 10
		cfg.MinConns = 2
		cfg.MaxConnLifetime = 1 * time.Hour
		cfg.HealthCheckPeriod = 30 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		pool, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			initErr = fmt.Errorf("failed to create pgx pool: %w", err)
			return
		}

		if err = pool.Ping(ctx); err != nil {
			initErr = fmt.Errorf("failed to ping PostgreSQL: %w", err)
			return
		}

		dbInstance = &PostgreSQL{
			Pool: pool,
		}
	})

	if initErr != nil {
		return nil, initErr
	}
	if dbInstance == nil {
		return nil, errors.New("failed to initialize PostgreSQL")
	}
	return dbInstance, nil
}
