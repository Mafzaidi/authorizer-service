package entity

import "time"

type Role struct {
	ID            string     `db:"id"`
	ApplicationID *string    `db:"application _id"`
	Code          string     `db:"code"`
	Name          string     `db:"role"`
	Description   *string    `db:"description"`
	Scope         *string    `db:"scope"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}
