package parser

import (
	"bufio"
	"strings"
	"testing"
)

func TestJSONStringParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "parse empty string",
			input:    `""`,
			expected: "",
		},
		{
			name:     "parse simple string",
			input:    `"hello"`,
			expected: "hello",
		},
		{
			name:     "parse string with spaces",
			input:    `"hello world"`,
			expected: "hello world",
		},
		{
			name:     "parse string with special characters",
			input:    `"hello\nworld"`,
			expected: `hello\nworld`,
		},
		{
			name:     "parse string with unicode",
			input:    `"hello\u0020world"`,
			expected: `hello\u0020world`,
		},
		{
			name:     "parse string with escaped quotes",
			input:    `"hello\"world\""`,
			expected: `hello\"world\"`,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result any
			parser, _ := NewJSONParser(bufio.NewReader(strings.NewReader(tt.input)), func(v any) error {
				result = v
				return nil
			})
			JSONString.ParseValue(parser)

			if result != tt.expected {
				t.Errorf("JSONString ParseValue result = %v, expected %v", result, tt.expected)
			}
		})
	}
}
