package idgen

import (
	"github.com/google/uuid"
)

func NewUUIDv7() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic("failed to generate UUID v7: " + err.Error())
	}
	return id.String()
}

func IsUUIDv7(id string) bool {
	u, err := uuid.Parse(id)
	if err != nil {
		return false
	}

	return u.Version() == 7
}
