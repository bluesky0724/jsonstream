package parser

import (
	"bufio"
	"strings"
	"testing"
)

func TestJSONTrueParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "parse true value",
			input:    "true",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result any
			parser, _ := NewJSONParser(bufio.NewReader(strings.NewReader(tt.input)), func(v any) error {
				result = v
				return nil
			})
			JSONTrue.ParseValue(parser)

			if result != tt.expected {
				t.Errorf("JSONTrue ParseValue result = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestJSONFalseParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "parse false value",
			input:    "false",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result any
			parser, _ := NewJSONParser(bufio.NewReader(strings.NewReader(tt.input)), func(v any) error {
				result = v
				return nil
			})
			JSONFalse.ParseValue(parser)

			if result != tt.expected {
				t.Errorf("JSONFalse ParseValue result = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestJSONNullParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			name:     "parse null value",
			input:    "null",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result any
			parser, _ := NewJSONParser(bufio.NewReader(strings.NewReader(tt.input)), func(v any) error {
				result = v
				return nil
			})
			JSONNull.ParseValue(parser)

			if result != tt.expected {
				t.Errorf("JSONNull ParseValue result = %v, expected %v", result, tt.expected)
			}
		})
	}
}
