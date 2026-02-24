package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
)

type Permission struct {
	ID            string
	ApplicationID pgtype.Text
	Code          string
	Description   pgtype.Text
	Version       int
	CreatedBy     pgtype.Text
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
	DeletedAt     pgtype.Timestamp
}

func (p *Permission) ToEntity() *entity.Permission {
	var appID *string
	if p.ApplicationID.Valid {
		appID = &p.ApplicationID.String
	}

	var desc *string
	if p.Description.Valid {
		desc = &p.Description.String
	}

	var createdBy *string
	if p.CreatedBy.Valid {
		desc = &p.CreatedBy.String
	}

	var deletedAt *time.Time
	if p.DeletedAt.Valid {
		deletedAt = &p.DeletedAt.Time
	}

	return &entity.Permission{
		ID:            p.ID,
		ApplicationID: appID,
		Code:          p.Code,
		Description:   desc,
		Version:       p.Version,
		CreatedBy:     createdBy,
		CreatedAt:     p.CreatedAt.Time,
		UpdatedAt:     p.UpdatedAt.Time,
		DeletedAt:     deletedAt,
	}
}
