package handlers

import (
	"github.com/google/uuid"
)

// ParseUUID converts a string to UUID
func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
