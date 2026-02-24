package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/internal/infrastructure/persistence/postgres/model"
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

	return scanUsers(rows)
}

func scanUser(row pgx.Row) (*entity.User, error) {
	var (
		user  entity.User
		model model.User
	)

	err := row.Scan(
		&model.ID,
		&model.Email,
		&model.Username,
		&model.Password,
		&model.FullName,
		&model.Phone,
		&model.IsActive,
		&model.EmailVerified,
		&model.PhoneVerified,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	user = *model.ToEntity()
	return &user, nil
}

func scanUsers(rows pgx.Rows) ([]*entity.User, error) {
	var users []*entity.User

	for rows.Next() {
		var model model.User
		if err := rows.Scan(
			&model.ID,
			&model.Email,
			&model.Username,
			&model.Password,
			&model.FullName,
			&model.Phone,
			&model.IsActive,
			&model.EmailVerified,
			&model.PhoneVerified,
			&model.CreatedAt,
			&model.UpdatedAt,
			&model.DeletedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, model.ToEntity())
	}

	return users, rows.Err()
}
