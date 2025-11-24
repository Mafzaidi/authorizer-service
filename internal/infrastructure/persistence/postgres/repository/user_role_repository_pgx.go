package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
)

type userRoleRepositoryPGX struct {
	pool *pgxpool.Pool
}

func NewUserRoleRepositoryPGX(pool *pgxpool.Pool) repository.UserRoleRepository {
	return &userRoleRepositoryPGX{
		pool: pool,
	}
}

func (r *userRoleRepositoryPGX) Assign(userID string, roleID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO authorizer_service.user_roles 
			(user_id, role_id)
		VALUES 
			($1, $2)
	`
	_, err := r.pool.Exec(ctx, query, userID, roleID)

	return err
}

func (r *userRoleRepositoryPGX) Unassign(userID, roleID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM authorizer_service.user_roles
		WHERE user_id = $1 AND role_id = $2;
	`

	_, err := r.pool.Exec(ctx, query, userID, roleID)
	return err
}

func (r *userRoleRepositoryPGX) GetRolesByUser(userID string) ([]*entity.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT r.id, r.name, r.description, r.deleted_at
		FROM authorizer_service.roles r
		INNER JOIN authorizer_service.user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND r.deleted_at IS NULL;
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.DeletedAt); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

func (r *userRoleRepositoryPGX) GetUsersByRole(roleID string) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
