package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
)

type userRepositoryPGX struct {
	pool *pgxpool.Pool
}

func NewUserRepositoryPGX(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepositoryPGX{
		pool: pool,
	}
}

func (r *userRepositoryPGX) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO authorizer_service.users 
			(id, email, username, password, full_name, phone, is_active, email_verified)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.Username, user.Password, user.FullName,
		user.Phone, user.IsActive,
		user.EmailVerified,
	)

	return err
}

func (r *userRepositoryPGX) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `SELECT * FROM authorizer_service.users WHERE id = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, id)
	return scanUser(row)
}

func (r *userRepositoryPGX) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT * FROM authorizer_service.users WHERE email = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, email)
	return scanUser(row)
}

func (r *userRepositoryPGX) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE authorizer_service.users
		SET full_name = $1,
			phone = $2,
			is_active = $3,
			updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
	`
	_, err := r.pool.Exec(ctx,
		query, user.FullName, user.Phone, user.IsActive, user.ID,
	)
	return err
}

func (r *userRepositoryPGX) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE authorizer_service.users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (r *userRepositoryPGX) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	query := `
		SELECT * FROM authorizer_service.users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func scanUser(row pgx.Row) (*entity.User, error) {
	var u entity.User

	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.Password,
		&u.FullName,
		&u.Phone,
		&u.IsActive,
		&u.EmailVerified,
		&u.PhoneVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &u, nil
}
