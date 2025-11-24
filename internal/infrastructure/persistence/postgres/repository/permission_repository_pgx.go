package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
)

type permissionRepositoryPGX struct {
	pool *pgxpool.Pool
}

func NewPermissionRepositoryPGX(pool *pgxpool.Pool) repository.PermissionRepository {
	return &permissionRepositoryPGX{
		pool: pool,
	}
}

func (r *permissionRepositoryPGX) Create(perm *entity.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO authorizer_service.permissions 
			(id, application _id, resource, action, slug, description, version, created_by)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		perm.ID, perm.ApplicationID, perm.Resource, perm.Action,
		perm.Slug, perm.Description, perm.Version, perm.CreatedBy,
	)

	return err
}

func (r *permissionRepositoryPGX) GetByID(id string) (*entity.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.permissions WHERE id = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, id)
	return scanPermission(row)
}

func (r *permissionRepositoryPGX) GetBySlug(slug string) (*entity.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.permissions WHERE slug = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, slug)
	return scanPermission(row)
}

func (r *permissionRepositoryPGX) GetByApplication(appID string) (*entity.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.permissions WHERE application_id = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, appID)
	return scanPermission(row)
}

func (r *permissionRepositoryPGX) Update(perm *entity.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE authorizer_service.permissions
		SET application _id = $1,
			resource = $2,
			action = $3,
			slug = $4,
			description = $5,
			version = $6,
			created_by = $7,
			updated_at = NOW()
		WHERE id = $8 AND deleted_at IS NULL
	`
	_, err := r.pool.Exec(ctx,
		query, perm.ApplicationID, perm.Resource, perm.Action,
		perm.Slug, perm.Description, perm.Version, perm.CreatedBy,
	)
	return err
}

func (r *permissionRepositoryPGX) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx, `UPDATE authorizer_service.permissions SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (r *permissionRepositoryPGX) List(limit, offset int) ([]*entity.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT * FROM authorizer_service.permissions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []*entity.Permission
	for rows.Next() {
		perm, err := scanPermission(rows)
		if err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	return perms, rows.Err()
}

func scanPermission(row pgx.Row) (*entity.Permission, error) {
	var p entity.Permission

	err := row.Scan(
		&p.ID,
		&p.ApplicationID,
		&p.Resource,
		&p.Action,
		&p.Slug,
		&p.Description,
		&p.Version,
		&p.CreatedBy,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}
