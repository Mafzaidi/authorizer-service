package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
	"github.com/mafzaidi/authorizer/internal/domain/repository"
	"github.com/mafzaidi/authorizer/internal/infrastructure/persistence/postgres/model"
	"github.com/mafzaidi/authorizer/pkg/utils"
)

type permRepositoryPGX struct {
	pool *pgxpool.Pool
}

func NewPermRepositoryPGX(pool *pgxpool.Pool) repository.PermRepository {
	return &permRepositoryPGX{
		pool: pool,
	}
}

func (r *permRepositoryPGX) Create(ctx context.Context, perm *entity.Permission) error {
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

func (r *permRepositoryPGX) Update(ctx context.Context, perm *entity.Permission) error {
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

func (r *permRepositoryPGX) Upsert(ctx context.Context, perm *entity.Permission) error {
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

func (r *permRepositoryPGX) BulkUpsert(ctx context.Context, perms []*entity.Permission) error {
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

func (r *permRepositoryPGX) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `UPDATE authorizer_service.permissions SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (r *permRepositoryPGX) GetByID(ctx context.Context, id string) (*entity.Permission, error) {
	query := `SELECT * FROM authorizer_service.permissions WHERE id = $1 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, id)
	return scanPerm(row)
}

func (r *permRepositoryPGX) GetByAppAndCode(ctx context.Context, appID, code string) (*entity.Permission, error) {
	query := `SELECT * FROM authorizer_service.permissions WHERE application_id = $1 AND code = $2 AND deleted_at IS NULL`

	row := r.pool.QueryRow(ctx, query, appID, code)
	return scanPerm(row)
}

func (r *permRepositoryPGX) List(ctx context.Context, limit, offset int) ([]*entity.Permission, error) {
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

	return scanPerms(rows)
}

func (r *permRepositoryPGX) ListByApp(ctx context.Context, appID string) ([]*entity.Permission, error) {
	query := `SELECT * FROM authorizer_service.permissions WHERE application_id = $1 AND deleted_at IS NULL`

	rows, err := r.pool.Query(ctx, query, appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanPerms(rows)
}

func scanPerm(row pgx.Row) (*entity.Permission, error) {
	var (
		perm  entity.Permission
		model model.Permission
	)

	err := row.Scan(
		&model.ID,
		&model.ApplicationID,
		&model.Code,
		&model.Description,
		&model.Version,
		&model.CreatedBy,
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

	perm = *model.ToEntity()
	return &perm, nil
}

func scanPerms(rows pgx.Rows) ([]*entity.Permission, error) {
	var perms []*entity.Permission

	for rows.Next() {
		var model model.Permission
		if err := rows.Scan(
			&model.ID,
			&model.ApplicationID,
			&model.Code,
			&model.Description,
			&model.Version,
			&model.CreatedBy,
			&model.CreatedAt,
			&model.UpdatedAt,
			&model.DeletedAt,
		); err != nil {
			return nil, err
		}
		perms = append(perms, model.ToEntity())
	}

	return perms, rows.Err()
}
