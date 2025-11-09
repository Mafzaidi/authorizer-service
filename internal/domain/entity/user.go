package entity

import (
	"time"
)

type User struct {
	ID            string                 `db:"id"`
	Email         string                 `db:"email"`
	Username      string                 `db:"username"`
	FullName      string                 `db:"full_name"`
	Phone         string                 `db:"phone"`
	PasswordHash  string                 `db:"password_hash"`
	Role          string                 `db:"role"`
	IsActive      bool                   `db:"is_active"`
	EmailVerified bool                   `db:"email_verified"`
	PhoneVerified bool                   `db:"phone_verified"`
	LastLoginAt   *time.Time             `db:"last_login_at"`
	CreatedAt     time.Time              `db:"created_at"`
	UpdatedAt     time.Time              `db:"updated_at"`
	DeletedAt     *time.Time             `db:"deleted_at"`
	Metadata      map[string]interface{} `db:"metadata"`
}
