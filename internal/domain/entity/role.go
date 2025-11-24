package entity

import "time"

type Role struct {
	ID            string     `db:"id"`
	ApplicationID string     `db:"application _id"`
	Slug          string     `db:"slug"`
	Name          string     `db:"role"`
	Description   string     `db:"description"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}
