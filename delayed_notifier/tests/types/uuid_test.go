package tests

import (
	"delayed_notifier/pkg/types"
	"github.com/google/uuid"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	// Generate multiple UUIDs to ensure they are different
	uuid1 := types.GenerateUUID()
	uuid2 := types.GenerateUUID()

	if uuid1.String() == uuid2.String() {
		t.Error("Generated UUIDs should be unique")
	}

	// Test that generated UUID is valid
	_, err := uuid.Parse(uuid1.String())
	if err != nil {
		t.Errorf("Generated UUID is invalid: %v", err)
	}

	_, err = uuid.Parse(uuid2.String())
	if err != nil {
		t.Errorf("Generated UUID is invalid: %v", err)
	}
}

func TestNewUUID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "valid UUID v4",
			input:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			expectError: false,
		},
		{
			name:        "valid UUID v1",
			input:       "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			expectError: false,
		},
		{
			name:        "invalid UUID format",
			input:       "not-a-uuid",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "wrong separator",
			input:       "f47ac10b_58cc_4372_a567_0e02b2c3d479",
			expectError: true,
		},
		{
			name:        "too short",
			input:       "f47ac10b-58cc-4372-a567-0e02b2c3d47",
			expectError: true,
		},
		{
			name:        "too long",
			input:       "f47ac10b-58cc-4372-a567-0e02b2c3d4799",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuidObj, err := types.NewUUID(tt.input)

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

			// Test that the string representation matches input
			result := uuidObj.String()
			if result != tt.input {
				t.Errorf("Expected '%s', got '%s'", tt.input, result)
			}

			// Test Value method
			value := uuidObj.Value()
			if value.String() != tt.input {
				t.Errorf("Value method returned unexpected UUID: expected '%s', got '%s'", tt.input, value.String())
			}
		})
	}
}

func TestUUID_Value(t *testing.T) {
	testUUID := "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	uuidObj, err := types.NewUUID(testUUID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	value := uuidObj.Value()

	if value.String() != testUUID {
		t.Errorf("Expected %s, got %s", testUUID, value.String())
	}

	// Verify it's a proper uuid.UUID type
	if value.Version() != 4 {
		t.Errorf("Expected UUID version 4, got %d", value.Version())
	}
}

func TestUUID_String(t *testing.T) {
	testUUID := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	uuidObj, err := types.NewUUID(testUUID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	result := uuidObj.String()
	if result != testUUID {
		t.Errorf("Expected '%s', got '%s'", testUUID, result)
	}
}

func TestUUID_Equality(t *testing.T) {
	uuidStr := "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	uuid1, err := types.NewUUID(uuidStr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	uuid2, err := types.NewUUID(uuidStr)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// They should have the same string representation
	if uuid1.String() != uuid2.String() {
		t.Error("UUIDs with same input string should be equal")
	}

	// They should have the same underlying UUID value
	if uuid1.Value() != uuid2.Value() {
		t.Error("UUID values should be equal for same input")
	}
}
