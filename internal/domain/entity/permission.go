package entity

import "time"

type Permission struct {
	ID            string     `db:"id"`
	ApplicationID string     `db:"application _id"`
	Resource      string     `db:"resource"`
	Action        string     `db:"action"`
	Slug          string     `db:"slug"`
	Description   string     `db:"description"`
	Version       int        `db:"version"`
	CreatedBy     string     `db:"created_by"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}
