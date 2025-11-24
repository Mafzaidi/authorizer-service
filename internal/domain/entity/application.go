package entity

import "time"

type Application struct {
	ID        string                 `db:"id"`
	Slug      string                 `db:"slug"`
	Name      string                 `db:"name"`
	Metadata  map[string]interface{} `db:"metadata"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedAt time.Time              `db:"updated_at"`
	DeletedAt *time.Time             `db:"deleted_at"`
}
