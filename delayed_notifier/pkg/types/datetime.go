package types

import (
	"fmt"
	"time"
)

const datetimeFormat = "2006-01-02 15:04:05"

// DateTime is a value type that contains a datetime
type DateTime struct {
	val time.Time
}

// NewDateTime creates a new DateTime from given time.Time value
func NewDateTime(val time.Time) DateTime {
	return DateTime{val: val}
}

// NewDateTimeFromString creates a new DateTime from string “2006-01-02 15:04:05“
func NewDateTimeFromString(val string) (DateTime, error) {
	timeParsed, err := time.Parse(datetimeFormat, val)
	if err != nil {
		return DateTime{}, fmt.Errorf("invalid datetime format: %w", err)
	}

	return NewDateTime(timeParsed), nil
}

// GreaterOrEqualThan returns if given date is later or equal than another one
func (d DateTime) GreaterOrEqualThan(other DateTime) bool {
	return d.val.After(other.val) || d.val.Equal(other.val)
}

// Value returns value of types.DateTime converted to time.Time
func (d DateTime) Value() time.Time {
	return d.val
}

// String returns value of types.DateTime converted to string
func (d DateTime) String() string {
	return d.Value().Format(datetimeFormat)
}
