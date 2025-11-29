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

type roleRepositoryPGX struct {
	pool *pgxpool.Pool
}

func NewRoleRepositoryPGX(pool *pgxpool.Pool) repository.RoleRepository {
	return &roleRepositoryPGX{
		pool: pool,
	}
}

func (r *roleRepositoryPGX) Create(role *entity.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO authorizer_service.roles 
			(id, application_id, code, name, description)
		VALUES 
			($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		role.ID, role.ApplicationID, role.Code, role.Name, role.Description,
	)

	return err
}

func (r *roleRepositoryPGX) Update(role *entity.Role) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE authorizer_service.roles
		SET application_id = $1,
			code = $2,
			name = $3,
			description = $4,
			updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
	`
	_, err := r.pool.Exec(ctx,
		query, role.ApplicationID, role.Code, role.Name, role.Description, role.ID,
	)
	return err
}

func (r *roleRepositoryPGX) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx, `UPDATE authorizer_service.roles SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (r *roleRepositoryPGX) GetByID(id string) (*entity.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.roles WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	return scanRole(row)
}

func (r *roleRepositoryPGX) GetByCode(code string) (*entity.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.roles WHERE code = $1`

	row := r.pool.QueryRow(ctx, query, code)
	return scanRole(row)
}

func (r *roleRepositoryPGX) GetByApp(appID string) ([]*entity.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.roles WHERE application_id = $1`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		role, err := scanRole(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *roleRepositoryPGX) GetByAppAndCode(appID, code string) (*entity.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.roles WHERE application_id = $1 AND code = $1`

	row := r.pool.QueryRow(ctx, query, appID, code)
	return scanRole(row)
}

func (r *roleRepositoryPGX) List(limit, offset int) ([]*entity.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT * FROM authorizer_service.roles
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		role, err := scanRole(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func scanRole(row pgx.Row) (*entity.Role, error) {
	var r entity.Role

	err := row.Scan(
		&r.ID,
		&r.ApplicationID,
		&r.Code,
		&r.Name,
		&r.Description,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}
