package extractor

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

func TestJSONExtractor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		base     string
		fields   []string
		expected [][]string
	}{
		{
			name:     "Simple array extraction",
			input:    `{"data":[{"id":1,"name":"John"},{"id":2,"name":"Jane"}]}`,
			base:     ".data",
			fields:   []string{"id", "name"},
			expected: [][]string{{"id", "name"}, {"1", "John"}, {"2", "Jane"}},
		},
		{
			name:     "Empty array extraction",
			input:    `{"data":[]}`,
			base:     ".data",
			fields:   []string{"id", "name"},
			expected: [][]string{{"id", "name"}},
		},
		{
			name:     "Missing fields extraction",
			input:    `{"data":[{"id":1},{"name":"Jane"}]}`,
			base:     ".data",
			fields:   []string{"id", "name"},
			expected: [][]string{{"id", "name"}, {"1", ""}, {"", "Jane"}},
		},
		{
			name:     "Nested fields extraction",
			input:    `{"data":[{"user":{"id":1,"details":{"name":"John"}}},{"user":{"id":2,"details":{"name":"Jane"}}}]}`,
			base:     ".data",
			fields:   []string{"user.id", "user.details.name"},
			expected: [][]string{{"user.id", "user.details.name"}, {"1", "John"}, {"2", "Jane"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare input reader
			reader := bufio.NewReader(strings.NewReader(tt.input))

			// Prepare output buffer
			var output bytes.Buffer
			writer := csv.NewWriter(&output)

			// Create and run extractor
			extractor, _ := NewJSONExtractor(reader, writer, tt.base, tt.fields)
			extractor.Extract()
			writer.Flush()

			// Parse output
			csvReader := csv.NewReader(strings.NewReader(output.String()))
			result, err := csvReader.ReadAll()
			if err != nil {
				t.Errorf("Failed to read CSV output: %v", err)
			}

			// Compare results
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d rows, got %d", len(tt.expected), len(result))
				return
			}

			for i := range result {
				if len(result[i]) != len(tt.expected[i]) {
					t.Errorf("Row %d: expected %d columns, got %d", i, len(tt.expected[i]), len(result[i]))
					continue
				}
				for j := range result[i] {
					if result[i][j] != tt.expected[i][j] {
						t.Errorf("Row %d, Column %d: expected %s, got %s", i, j, tt.expected[i][j], result[i][j])
					}
				}
			}
		})
	}
}

func TestGetAbsolutePath(t *testing.T) {
	tests := []struct {
		base     string
		field    string
		expected string
	}{
		{"data", "id", "data.id"},
		{"users", "name", "users.name"},
		{"items", "details.price", "items.details.price"},
	}

	for _, tt := range tests {
		result := getAbsolutePath(tt.base, tt.field)
		if result != tt.expected {
			t.Errorf("getAbsolutePath(%s, %s) = %s; want %s", tt.base, tt.field, result, tt.expected)
		}
	}
}

func TestShouldUpdate(t *testing.T) {
	extractor := &JSONExtractor{
		base:    "data",
		targets: []string{"id", "name"},
	}

	tests := []struct {
		field    string
		expected bool
	}{
		{"data.id", true},
		{"data.name", true},
		{"data.age", false},
		{"users.id", false},
	}

	for _, tt := range tests {
		result := extractor.shouldUpdate(tt.field)
		if result != tt.expected {
			t.Errorf("shouldUpdate(%s) = %v; want %v", tt.field, result, tt.expected)
		}
	}
}
