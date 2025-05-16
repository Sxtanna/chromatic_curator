package common

import (
	"testing"
)

func TestParseTextToColorInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		// Integer input
		{
			name:     "Integer input",
			input:    "16711680", // Decimal for red (0xFF0000)
			expected: 16711680,
			wantErr:  false,
		},
		// Color name input
		{
			name:     "Color name - Red",
			input:    "Red",
			expected: 16711680, // 0xFF0000
			wantErr:  false,
		},
		{
			name:     "Color name - case insensitive",
			input:    "blue",
			expected: 255, // 0x0000FF
			wantErr:  false,
		},
		// Hex code with # prefix
		{
			name:     "Hex code with # prefix - Red",
			input:    "#FF0000",
			expected: 16711680, // 0xFF0000
			wantErr:  false,
		},
		// Hex code without # prefix
		{
			name:     "Hex code without # prefix - Green",
			input:    "00FF00",
			expected: 65280, // 0x00FF00
			wantErr:  false,
		},
		// Invalid input
		{
			name:     "Invalid hex code",
			input:    "GGGGGG",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
			wantErr:  true,
		},
		// Short hex format
		{
			name:     "Short hex format with # prefix",
			input:    "#F00",
			expected: 16711680, // 0xFF0000
			wantErr:  false,
		},
		{
			name:     "Short hex format without # prefix",
			input:    "F00",
			expected: 16711680, // 0xFF0000
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTextToColorInt(tt.input)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTextToColorInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check result
			if got != tt.expected {
				t.Errorf("ParseTextToColorInt() = %v, want %v", got, tt.expected)
			}
		})
	}
}
