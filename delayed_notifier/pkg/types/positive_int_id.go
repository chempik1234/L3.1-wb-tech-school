package types

import (
	"errors"
)

// ErrLessThanZero is an error when given value is less than 0
var ErrLessThanZero = errors.New("value must be greater than 0")

// PositiveIntID is a value type for int id, wrapper around int
//
// check > 0
type PositiveIntID struct {
	value int
}

// NewPositiveIntID creates a new PositiveIntID from int, validating it (>0)
func NewPositiveIntID(value int) (PositiveIntID, error) {
	if value <= 0 {
		return PositiveIntID{}, ErrLessThanZero
	}
	return PositiveIntID{value}, nil
}

// Value returns value of types.PositiveIntID of type int
func (v PositiveIntID) Value() int {
	return v.value
}
