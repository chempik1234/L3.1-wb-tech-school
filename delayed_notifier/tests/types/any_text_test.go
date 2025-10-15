package tests

import (
	"delayed_notifier/pkg/types"
	"testing"
)

func TestAnyText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "regular text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "text with special characters",
			input:    "Hello, 世界! 🎉",
			expected: "Hello, 世界! 🎉",
		},
		{
			name:     "text with numbers",
			input:    "123 ABC",
			expected: "123 ABC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test types.NewAnyText
			anyText := types.NewAnyText(tt.input)

			// Test String() method
			result := anyText.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}

			// Test type conversion
			if string(anyText) != tt.expected {
				t.Errorf("Type conversion failed: expected '%s', got '%s'", tt.expected, string(anyText))
			}
		})
	}
}

func TestAnyText_Equality(t *testing.T) {
	text1 := types.NewAnyText("same text")
	text2 := types.NewAnyText("same text")
	text3 := types.NewAnyText("different text")

	if text1.String() != text2.String() {
		t.Error("Same text values should be equal")
	}

	if text1.String() == text3.String() {
		t.Error("Different text values should not be equal")
	}
}
