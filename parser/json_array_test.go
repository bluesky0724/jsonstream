package parser

import (
	"bufio"
	"strings"
	"testing"
)

func TestJSONArrayParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []any
	}{
		{
			name:     "parse empty array",
			input:    `[]`,
			expected: []any{nil},
		},
		{
			name:     "parse array with single number",
			input:    `[1]`,
			expected: []any{1, nil},
		},
		{
			name:     "parse array with multiple numbers",
			input:    `[1,2,3]`,
			expected: []any{1, 2, 3, nil},
		},
		{
			name:     "parse array with strings",
			input:    `["hello","world"]`,
			expected: []any{"hello", "world", nil},
		},
		{
			name:     "parse array with mixed types",
			input:    `[1,"hello",true,null]`,
			expected: []any{1, "hello", true, nil, nil},
		},
		{
			name:     "parse nested arrays",
			input:    `[[1,2],[3,4]]`,
			expected: []any{1, 2, nil, 3, 4, nil, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result = []any{}
			parser, _ := NewJSONParser(bufio.NewReader(strings.NewReader(tt.input)), func(v any) error {
				result = append(result, v)
				return nil
			})
			JSONArray.ParseValue(parser)

			if !CompareArray(result, tt.expected) {
				t.Errorf("Parse() = %v\n, want %v", result, tt.expected)
			}
		})
	}
}
