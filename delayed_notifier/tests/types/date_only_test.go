package tests

import (
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
	"testing"
	"time"
)

func TestNewDateOnlyFromTime(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "regular date",
			input:    time.Date(2023, time.December, 25, 15, 30, 45, 0, time.UTC),
			expected: "2023-12-25",
		},
		{
			name:     "leap year date",
			input:    time.Date(2020, time.February, 29, 0, 0, 0, 0, time.UTC),
			expected: "2020-02-29",
		},
		{
			name:     "first day of year",
			input:    time.Date(2023, time.January, 1, 23, 59, 59, 999, time.UTC),
			expected: "2023-01-01",
		},
		{
			name:     "last day of year",
			input:    time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
			expected: "2023-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dateOnly := types.NewDateOnlyFromTime(tt.input)

			result := dateOnly.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}

			// Verify time components are truncated
			value := dateOnly.Value()
			if value.Hour() != 0 || value.Minute() != 0 || value.Second() != 0 || value.Nanosecond() != 0 {
				t.Error("Time components should be truncated to zero")
			}
		})
	}
}

func TestNewDateOnlyFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "valid date",
			input:       "2023-12-25",
			expected:    "2023-12-25",
			expectError: false,
		},
		{
			name:        "valid date with single digit month and day",
			input:       "2023-01-01",
			expected:    "2023-01-01",
			expectError: false,
		},
		{
			name:        "leap year date",
			input:       "2020-02-29",
			expected:    "2020-02-29",
			expectError: false,
		},
		{
			name:        "invalid format",
			input:       "2023/12/25",
			expected:    "",
			expectError: true,
		},
		{
			name:        "invalid date",
			input:       "2023-02-30",
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "wrong separator",
			input:       "2023.12.25",
			expected:    "",
			expectError: true,
		},
		{
			name:        "non-date string",
			input:       "not-a-date",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dateOnly, err := types.NewDateOnlyFromString(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			result := dateOnly.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestDateOnly_GreaterOrEqualThan(t *testing.T) {
	date1, _ := types.NewDateOnlyFromString("2023-01-01")
	date2, _ := types.NewDateOnlyFromString("2023-06-15")
	date3, _ := types.NewDateOnlyFromString("2023-06-15") // same as date2
	date4, _ := types.NewDateOnlyFromString("2024-01-01")

	tests := []struct {
		name     string
		date1    types.DateOnly
		date2    types.DateOnly
		expected bool
	}{
		{
			name:     "earlier date compared to later date",
			date1:    date1,
			date2:    date2,
			expected: false,
		},
		{
			name:     "later date compared to earlier date",
			date1:    date2,
			date2:    date1,
			expected: true,
		},
		{
			name:     "equal dates",
			date1:    date2,
			date2:    date3,
			expected: true,
		},
		{
			name:     "same date instance",
			date1:    date2,
			date2:    date2,
			expected: true,
		},
		{
			name:     "different year",
			date1:    date4,
			date2:    date1,
			expected: true,
		},
		{
			name:     "same year different month",
			date1:    date2,
			date2:    date1,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.date1.GreaterOrEqualThan(tt.date2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for %s >= %s",
					tt.expected, result, tt.date1.String(), tt.date2.String())
			}
		})
	}
}

func TestDateOnly_Value(t *testing.T) {
	dateOnly, _ := types.NewDateOnlyFromString("2023-12-25")
	value := dateOnly.Value()

	expectedYear := 2023
	expectedMonth := time.December
	expectedDay := 25

	if value.Year() != expectedYear {
		t.Errorf("Expected year %d, got %d", expectedYear, value.Year())
	}
	if value.Month() != expectedMonth {
		t.Errorf("Expected month %v, got %v", expectedMonth, value.Month())
	}
	if value.Day() != expectedDay {
		t.Errorf("Expected day %d, got %d", expectedDay, value.Day())
	}

	// Check that time is set to midnight
	if value.Hour() != 0 || value.Minute() != 0 || value.Second() != 0 || value.Nanosecond() != 0 {
		t.Error("Time should be set to midnight")
	}
}

func TestDateOnly_String(t *testing.T) {
	tests := []struct {
		name     string
		date     types.DateOnly
		expected string
	}{
		{
			name:     "regular date",
			date:     types.NewDateOnlyFromTime(time.Date(2023, 12, 25, 0, 0, 0, 0, time.Local)),
			expected: "2023-12-25",
		},
		{
			name:     "single digit month and day",
			date:     types.NewDateOnlyFromTime(time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)),
			expected: "2023-01-01",
		},
		{
			name:     "february in leap year",
			date:     types.NewDateOnlyFromTime(time.Date(2020, 2, 29, 0, 0, 0, 0, time.Local)),
			expected: "2020-02-29",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.date.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
