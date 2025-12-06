package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
)

type rolePermRepositoryPGX struct {
	pool *pgxpool.Pool
}

func NewRolePermRepositoryPGX(pool *pgxpool.Pool) repository.RolePermRepository {
	return &rolePermRepositoryPGX{
		pool: pool,
	}
}

func (r *rolePermRepositoryPGX) Grant(roleID, permID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO authorizer_service.role_permissions 
			(role_id, permission_id)
		VALUES 
			($1, $2);
	`
	_, err := r.pool.Exec(ctx, query, roleID, permID)

	return err
}

func (r *rolePermRepositoryPGX) Revoke(roleID, permID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM authorizer_service.role_permissions
		WHERE role_id = $1 AND permission_id = $2;
	`

	_, err := r.pool.Exec(ctx, query, roleID, permID)
	return err
}

func (r *rolePermRepositoryPGX) Replace(roleID string, permIDs []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	delQuery := `
		DELETE FROM authorizer_service.role_permissions
		WHERE role_id = $1;
	`
	if _, err := tx.Exec(ctx, delQuery, roleID); err != nil {
		return err
	}

	if len(permIDs) == 0 {
		return tx.Commit(ctx)
	}

	insQuery := `
		INSERT INTO authorizer_service.role_permissions (role_id, permission_id)
		SELECT $1, unnest($2::uuid[]);
	`
	if _, err := tx.Exec(ctx, insQuery, roleID, permIDs); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *rolePermRepositoryPGX) GetPermsByRole(roleID string) ([]*entity.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT p.id, p.code, p.description, p.deleted_at
		FROM authorizer_service.permissions p
		INNER JOIN authorizer_service.role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1 AND r.deleted_at IS NULL;
	`

	rows, err := r.pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []*entity.Permission
	for rows.Next() {
		var perm entity.Permission
		if err := rows.Scan(&perm.ID, &perm.Code, &perm.Description, &perm.DeletedAt); err != nil {
			return nil, err
		}
		perms = append(perms, &perm)
	}

	return perms, nil
}

func (r *rolePermRepositoryPGX) GetPermsByRoles(roleIDs []string) ([]*entity.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT p.id, p.code, p.description, p.deleted_at
		FROM authorizer_service.permissions p
		INNER JOIN authorizer_service.role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = ANY($1) AND p.deleted_at IS NULL;
	`

	rows, err := r.pool.Query(ctx, query, roleIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []*entity.Permission
	for rows.Next() {
		var perm entity.Permission
		if err := rows.Scan(&perm.ID, &perm.Code, &perm.Description, &perm.DeletedAt); err != nil {
			return nil, err
		}
		perms = append(perms, &perm)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return perms, nil
}
