package entity

import (
	"time"
)

type User struct {
	ID            string     `db:"id"`
	Email         string     `db:"email"`
	Username      string     `db:"username"`
	Password      string     `db:"password"`
	FullName      string     `db:"full_name"`
	Phone         *string    `db:"phone"`
	IsActive      bool       `db:"is_active"`
	EmailVerified bool       `db:"email_verified"`
	PhoneVerified bool       `db:"phone_verified"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}
