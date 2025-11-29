package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
)

type appRepositoryPGX struct {
	pool *pgxpool.Pool
}

func NewAppRepositoryPGX(pool *pgxpool.Pool) repository.AppRepository {
	return &appRepositoryPGX{
		pool: pool,
	}
}

func (r *appRepositoryPGX) Create(app *entity.Application) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metadataJSON, _ := json.Marshal(app.Metadata)

	query := `
		INSERT INTO authorizer_service.applications 
			(id, code, name, description, metadata)
		VALUES 
			($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		app.ID, app.Code, app.Name, app.Description, metadataJSON,
	)

	return err
}

func (r *appRepositoryPGX) GetByID(id string) (*entity.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.applications WHERE id = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, id)
	return scanApp(row)
}

func (r *appRepositoryPGX) GetByCode(code string) (*entity.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.applications WHERE code = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, code)
	return scanApp(row)
}

func (r *appRepositoryPGX) Update(app *entity.Application) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE authorizer_service.applications
		SET name = $1,
			updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := r.pool.Exec(ctx,
		query, app.Name, app.ID,
	)
	return err
}

func (r *appRepositoryPGX) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx, `UPDATE authorizer_service.applications SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (r *appRepositoryPGX) List(limit, offset int) ([]*entity.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT * FROM authorizer_service.applications
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []*entity.Application
	for rows.Next() {
		app, err := scanApp(rows)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, rows.Err()
}

func scanApp(row pgx.Row) (*entity.Application, error) {
	var a entity.Application
	var metadataJSON []byte

	err := row.Scan(
		&a.ID,
		&a.Code,
		&a.Name,
		&a.Description,
		&metadataJSON,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	_ = json.Unmarshal(metadataJSON, &a.Metadata)
	return &a, nil
}
