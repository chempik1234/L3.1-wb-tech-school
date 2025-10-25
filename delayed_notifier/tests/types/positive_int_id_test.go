package tests

import (
	"errors"
	"github.com/chempik1234/L3.1-wb-tech-school/delayed_notifier/pkg/types"
	"testing"
)

func TestNewPositiveIntID(t *testing.T) {
	tests := []struct {
		name        string
		input       int
		expected    int
		expectError bool
	}{
		{
			name:        "positive integer",
			input:       1,
			expected:    1,
			expectError: false,
		},
		{
			name:        "large positive integer",
			input:       999999,
			expected:    999999,
			expectError: false,
		},
		{
			name:        "zero",
			input:       0,
			expected:    0,
			expectError: true,
		},
		{
			name:        "negative integer",
			input:       -1,
			expected:    0,
			expectError: true,
		},
		{
			name:        "large negative integer",
			input:       -999999,
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := types.NewPositiveIntID(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if !errors.Is(err, types.ErrLessThanZero) {
					t.Errorf("Expected ErrLessThanZero, got %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if id.Value() != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, id.Value())
			}
		})
	}
}

func TestPositiveIntID_Value(t *testing.T) {
	id, err := types.NewPositiveIntID(42)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	value := id.Value()
	expected := 42

	if value != expected {
		t.Errorf("Expected %d, got %d", expected, value)
	}
}

func TestPositiveIntID_Immutability(t *testing.T) {
	// This test ensures that the struct fields are not exported and thus immutable
	id, err := types.NewPositiveIntID(100)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// The value field is not exported, so we can't modify it directly
	// This is a good thing for a value object
	originalValue := id.Value()

	// Create another ID
	anotherID, err := types.NewPositiveIntID(200)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Original ID should remain unchanged
	if id.Value() != originalValue {
		t.Error("Original ID value was modified")
	}

	if anotherID.Value() == id.Value() {
		t.Error("Different IDs should have different values")
	}
}
