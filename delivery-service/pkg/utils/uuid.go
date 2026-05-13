package utils

import "github.com/google/uuid"

func ParseUUID(value string) (uuid.UUID, bool) {
	id, err := uuid.Parse(value)
	return id, err == nil
}
