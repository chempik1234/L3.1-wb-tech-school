package types

import (
	"fmt"
	"time"
)

// DateOnly is a value type that stores year, day, month
type DateOnly struct {
	year  int
	month time.Month
	day   int
}

// NewDateOnlyFromTime creates a new DateOnly from time.Time, truncating time of day, leaving only date
func NewDateOnlyFromTime(t time.Time) DateOnly {
	return DateOnly{year: t.Year(), month: t.Month(), day: t.Day()}
}

// NewDateOnlyFromString creates a new DateOnly from 'YYYY-MM-DD' string
func NewDateOnlyFromString(s string) (DateOnly, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return DateOnly{}, fmt.Errorf("invalid date format: %w", err)
	}
	return NewDateOnlyFromTime(t), nil
}

// GreaterOrEqualThan returns if given date is later or equal than another one
func (d DateOnly) GreaterOrEqualThan(other DateOnly) bool {
	if d.year != other.year {
		return d.year >= other.year
	}
	if d.month != other.month {
		return d.month >= other.month
	}
	return d.day >= other.day
}

// Value returns value of types.DateOnly converted to time.Time
func (d DateOnly) Value() time.Time {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, time.Local)
}

// String returns value of types.DateOnly converted to string
func (d DateOnly) String() string {
	return d.Value().Format("2006-01-02")
}
