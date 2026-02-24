package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mafzaidi/authorizer/internal/domain/entity"
)

type User struct {
	ID            string
	Email         string
	Username      string
	Password      string
	FullName      string
	Phone         pgtype.Text
	IsActive      bool
	EmailVerified bool
	PhoneVerified bool
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
	DeletedAt     pgtype.Timestamp
}

func (u *User) ToEntity() *entity.User {

	var phone *string
	if u.Phone.Valid {
		phone = &u.Phone.String
	}

	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}

	return &entity.User{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		Password:      u.Password,
		FullName:      u.FullName,
		Phone:         phone,
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		PhoneVerified: u.PhoneVerified,
		CreatedAt:     u.CreatedAt.Time,
		UpdatedAt:     u.UpdatedAt.Time,
		DeletedAt:     deletedAt,
	}
}
