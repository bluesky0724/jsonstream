package parser

import (
	"bufio"
	"strings"
	"testing"
)

func TestJSONNumberParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "parse integer value",
			input:    "123",
			expected: 123,
		},
		{
			name:     "parse float value",
			input:    "123.456",
			expected: 123.456,
		},
		{
			name:     "parse negative value",
			input:    "-123",
			expected: -123,
		},
		{
			name:     "parse scientific notation",
			input:    "1.23e2",
			expected: 123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result any
			parser, _ := NewJSONParser(bufio.NewReader(strings.NewReader(tt.input)), func(v any) error {
				result = v
				return nil
			})
			JSONNumber.ParseValue(parser)

			if result != tt.expected {
				t.Errorf("JSONNumber ParseValue result = %v, expected %v", result, tt.expected)
			}
		})
	}
}
