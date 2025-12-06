package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localdev.me/authorizer/internal/domain/entity"
	"localdev.me/authorizer/internal/domain/repository"
	"localdev.me/authorizer/pkg/utils"
)

type permRepositoryPGX struct {
	pool *pgxpool.Pool
}

func NewPermRepositoryPGX(pool *pgxpool.Pool) repository.PermRepository {
	return &permRepositoryPGX{
		pool: pool,
	}
}

func (r *permRepositoryPGX) Create(perm *entity.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO authorizer_service.permissions 
			(id, application_id, code, description, version, created_by)
		VALUES 
			($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		perm.ID, perm.ApplicationID, perm.Code, perm.Description, perm.Version, perm.CreatedBy,
	)

	return err
}

func (r *permRepositoryPGX) Update(perm *entity.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE authorizer_service.permissions
		SET application_id = $1,
			code = $2,
			description = $3,
			version = $4,
			created_by = $5,
			updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
	`
	_, err := r.pool.Exec(ctx,
		query, perm.ApplicationID, perm.Code, perm.Description, perm.Version, perm.CreatedBy, perm.ID,
	)
	return err
}

func (r *permRepositoryPGX) Upsert(perm *entity.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO authorizer_service.permissions (id, application_id, code, description, version)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (application_id, code)
        DO UPDATE SET description = EXCLUDED.description,
                    updated_at = NOW()
	`
	_, err := r.pool.Exec(ctx,
		query, perm.ID, perm.ApplicationID, perm.Code, perm.Description, perm.Version,
	)
	return err
}

func (r *permRepositoryPGX) BulkUpsert(perms []*entity.Permission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	const chunkSize = 100

	for _, c := range utils.Chunk(perms, chunkSize) {
		if err := r.bulkUpsertChunk(ctx, c); err != nil {
			return err
		}
	}
	return nil
}

func (r *permRepositoryPGX) bulkUpsertChunk(ctx context.Context, perms []*entity.Permission) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	for _, p := range perms {
		batch.Queue(`
            INSERT INTO authorizer_service.permissions (id, application_id, code, description, version)
            VALUES ($1, $2, $3, $4, $5)
            ON CONFLICT (application_id, code)
            DO UPDATE SET
                description = EXCLUDED.description,
				version = EXCLUDED.version,
                updated_at = EXCLUDED.updated_at
        `,
			p.ID,
			p.ApplicationID,
			p.Code,
			p.Description,
			p.Version,
		)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *permRepositoryPGX) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.pool.Exec(ctx, `UPDATE authorizer_service.permissions SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (r *permRepositoryPGX) GetByID(id string) (*entity.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.permissions WHERE id = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, id)
	return scanPerm(row)
}

func (r *permRepositoryPGX) GetByAppAndCode(appID, code string) (*entity.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.permissions WHERE application_id = $1 AND code = $2 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, appID, code)
	return scanPerm(row)
}

func (r *permRepositoryPGX) List(limit, offset int) ([]*entity.Permission, error) {
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
		perm, err := scanPerm(rows)
		if err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	return perms, rows.Err()
}

func (r *permRepositoryPGX) ListByApp(appID string) ([]*entity.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT * FROM authorizer_service.permissions WHERE application_id = $1 AND deleted_at IS NULL`

	rows, err := r.pool.Query(ctx, query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []*entity.Permission
	for rows.Next() {
		perm, err := scanPerm(rows)
		if err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	return perms, rows.Err()
}

func scanPerm(row pgx.Row) (*entity.Permission, error) {
	var p entity.Permission

	err := row.Scan(
		&p.ID,
		&p.ApplicationID,
		&p.Code,
		&p.Description,
		&p.Version,
		&p.CreatedBy,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &p, nil
}
