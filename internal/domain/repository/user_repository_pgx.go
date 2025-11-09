package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localdev.me/authorizer/internal/domain/entity"
)

type userRepositoryPGX struct {
	pool *pgxpool.Pool
}

func NewUserRepositoryPGX(pool *pgxpool.Pool) UserRepository {
	return &userRepositoryPGX{
		pool: pool,
	}
}

func (r *userRepositoryPGX) Create(user *entity.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metadataJSON, _ := json.Marshal(user.Metadata)

	query := `
		INSERT INTO users 
			(id, email, username, full_name, phone, password_hash, role, is_active, email_verified, metadata)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.Username, user.FullName,
		user.Phone, user.PasswordHash, user.Role, user.IsActive,
		user.EmailVerified, metadataJSON,
	)

	return err
}

func (r *userRepositoryPGX) GetByID(id string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, id)
	return scanUser(row)
}

func (r *userRepositoryPGX) GetByEmail(email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, email)
	return scanUser(row)
}

func (r *userRepositoryPGX) Update(user *entity.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET full_name = $1,
			phone = $2,
			role = $3,
			is_active = $4,
			updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
	`
	_, err := r.pool.Exec(ctx,
		query, user.FullName, user.Phone, user.Role, user.IsActive, user.ID,
	)
	return err
}

func (r *userRepositoryPGX) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx, `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (r *userRepositoryPGX) List(limit, offset int) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT * FROM users
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
	var metadataJSON []byte

	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.FullName,
		&u.Phone,
		&u.PasswordHash,
		&u.Role,
		&u.IsActive,
		&u.EmailVerified,
		&u.LastLoginAt,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
		&metadataJSON,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	_ = json.Unmarshal(metadataJSON, &u.Metadata)
	return &u, nil
}
