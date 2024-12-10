package parser

import (
	"bufio"
	"strings"
	"testing"
)

func TestJSONObjectParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []any
	}{
		{
			name:     "parse empty object",
			input:    `{}`,
			expected: []any{nil},
		},
		{
			name:     "parse simple object",
			input:    `{"key": "value"}`,
			expected: []any{"value", nil},
		},
		{
			name:     "parse object with multiple properties",
			input:    `{"key1": "value1", "key2": "value2"}`,
			expected: []any{"value1", "value2", nil},
		},
		{
			name:     "parse nested object",
			input:    `{"outer": {"inner": "value"}}`,
			expected: []any{"value", nil, nil},
		},
		{
			name:     "parse object with different value types",
			input:    `{"string": "value", "number": 42, "boolean": true, "null": null}`,
			expected: []any{"value", 42, true, nil, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []any
			parser, _ := NewJSONParser(bufio.NewReader(strings.NewReader(tt.input)), func(v any) error {
				result = append(result, v)
				return nil
			})
			JSONObject.ParseValue(parser)

			if !CompareArray(result, tt.expected) {
				t.Errorf("JSONObject ParseValue result = %v, expected %v", result, tt.expected)
			}
		})
	}
}
