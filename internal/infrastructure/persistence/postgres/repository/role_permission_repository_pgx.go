package repository

import (
	"context"

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

func (r *rolePermRepositoryPGX) Grant(ctx context.Context, roleID, permID string) error {
	query := `
		INSERT INTO authorizer_service.role_permissions 
			(role_id, permission_id)
		VALUES 
			($1, $2);
	`
	_, err := r.pool.Exec(ctx, query, roleID, permID)

	return err
}

func (r *rolePermRepositoryPGX) Revoke(ctx context.Context, roleID, permID string) error {
	query := `
		DELETE FROM authorizer_service.role_permissions
		WHERE role_id = $1 AND permission_id = $2;
	`

	_, err := r.pool.Exec(ctx, query, roleID, permID)
	return err
}

func (r *rolePermRepositoryPGX) Replace(ctx context.Context, roleID string, permIDs []string) error {
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

func (r *rolePermRepositoryPGX) GetPermsByRole(ctx context.Context, roleID string) ([]*entity.Permission, error) {
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

func (r *rolePermRepositoryPGX) GetPermsByRoles(ctx context.Context, roleIDs []string) ([]*entity.Permission, error) {
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
