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
