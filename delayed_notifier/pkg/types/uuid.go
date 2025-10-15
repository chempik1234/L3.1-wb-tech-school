package types

import (
	"fmt"
	"github.com/google/uuid"
)

// UUID is a value type that stores a valid UUID
type UUID struct {
	value uuid.UUID
}

// GenerateUUID creates a new valid UUID
func GenerateUUID() UUID {
	return UUID{
		value: uuid.New(),
	}
}

// NewUUID creates UUID from given string, returns err if invalid
func NewUUID(id string) (UUID, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return UUID{}, fmt.Errorf("invalid uuid '%s': %w", id, err)
	}
	return UUID{
		value: idUUID,
	}, nil
}

// Value returns value of types.UUID of type uuid.UUID
func (v UUID) Value() uuid.UUID {
	return v.value
}

// String returns value of types.UUID converted to string
func (v UUID) String() string {
	return v.value.String()
}
