package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectID(hex string) (primitive.ObjectID, error) {
	if hex == "" {
		return primitive.NilObjectID, nil
	}
	return primitive.ObjectIDFromHex(hex)
}
