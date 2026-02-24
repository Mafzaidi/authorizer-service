package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
)

type userRoleRepositoryPGX struct {
	pool *pgxpool.Pool
}

func NewUserRoleRepositoryPGX(pool *pgxpool.Pool) repository.UserRoleRepository {
	return &userRoleRepositoryPGX{
		pool: pool,
	}
}

func (r *userRoleRepositoryPGX) Assign(ctx context.Context, userID string, roleIDs []string) error {
	query := `
		INSERT INTO authorizer_service.user_roles 
			(user_id, role_id)
		SELECT $1, unnest($2::uuid[]);
	`
	_, err := r.pool.Exec(ctx, query, userID, roleIDs)

	return err
}

func (r *userRoleRepositoryPGX) Unassign(ctx context.Context, userID, roleID string) error {
	query := `
		DELETE FROM authorizer_service.user_roles
		WHERE user_id = $1 AND role_id = $2;
	`

	_, err := r.pool.Exec(ctx, query, userID, roleID)
	return err
}

func (r *userRoleRepositoryPGX) Replace(ctx context.Context, userID string, roleIDs []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	delQuery := `
		DELETE FROM authorizer_service.user_roles
		WHERE user_id = $1;
	`
	if _, err := tx.Exec(ctx, delQuery, userID); err != nil {
		return err
	}

	if len(roleIDs) == 0 {
		return tx.Commit(ctx)
	}

	insQuery := `
		INSERT INTO authorizer_service.user_roles (user_id, role_id)
		SELECT $1, unnest($2::uuid[]);
	`
	if _, err := tx.Exec(ctx, insQuery, userID, roleIDs); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *userRoleRepositoryPGX) GetRolesByUser(ctx context.Context, userID string) ([]*entity.Role, error) {
	query := `
		SELECT r.id, r.application_id, r.code, r.name, r.description, r.deleted_at
		FROM authorizer_service.roles r
		INNER JOIN authorizer_service.user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND r.deleted_at IS NULL;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRoles(rows)
}

func (r *userRoleRepositoryPGX) GetRolesByUserAndApp(ctx context.Context, userID, appID string) ([]*entity.Role, error) {
	query := `
		SELECT r.*
		FROM authorizer_service.roles r
		INNER JOIN authorizer_service.user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND r.application_id = $2 AND r.deleted_at IS NULL;
	`

	rows, err := r.pool.Query(ctx, query, userID, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRoles(rows)
}

func (r *userRoleRepositoryPGX) GetGlobalRolesByUser(ctx context.Context, userID string) ([]*entity.Role, error) {
	query := `
		SELECT r.*
		FROM authorizer_service.roles r
		INNER JOIN authorizer_service.user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND r.scope ='GLOBAL' AND r.deleted_at IS NULL;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRoles(rows)
}

func (r *userRoleRepositoryPGX) GetUsersByRole(ctx context.Context, roleID string) ([]*entity.User, error) {
	query := `
		SELECT u.id, u.username, u.full_name, r.deleted_at
		FROM authorizer_service.roles r
		INNER JOIN authorizer_service.user_roles ur ON ur.role_id = r.id
		INNER JOIN authorizer_service.users u ON u.id = ur.user_id
		WHERE r.role_id = $1 AND r.deleted_at IS NULL;
	`

	rows, err := r.pool.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
