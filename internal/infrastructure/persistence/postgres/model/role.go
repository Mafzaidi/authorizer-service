package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
)

type Role struct {
	ID            string
	ApplicationID pgtype.Text
	Code          string
	Name          string
	Description   pgtype.Text
	Scope         pgtype.Text
	CreatedAt     pgtype.Timestamp
	DeletedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
}

func (r *Role) ToEntity() *entity.Role {
	var appID *string
	if r.ApplicationID.Valid {
		appID = &r.ApplicationID.String
	}

	var desc *string
	if r.Description.Valid {
		desc = &r.Description.String
	}

	var scope *string
	if r.Scope.Valid {
		scope = &r.Scope.String
	}

	var deletedAt *time.Time
	if r.DeletedAt.Valid {
		deletedAt = &r.DeletedAt.Time
	}

	return &entity.Role{
		ID:            r.ID,
		ApplicationID: appID,
		Code:          r.Code,
		Name:          r.Name,
		Description:   desc,
		Scope:         scope,
		CreatedAt:     r.CreatedAt.Time,
		UpdatedAt:     r.UpdatedAt.Time,
		DeletedAt:     deletedAt,
	}
}
