package parser

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestChunkSize(t *testing.T) {
	if ChunkSize != 1024 {
		t.Errorf("ChunkSize = %v, want 1024", ChunkSize)
	}
}

func TestJSONValueType(t *testing.T) {
	tests := []struct {
		name string
		jvt  JSONValueType
	}{
		{
			name: "empty parse value",
			jvt:  JSONValueType{ParseValue: nil},
		},
		{
			name: "with parse value function",
			jvt: JSONValueType{
				ParseValue: func(p *JSONParser) error { return nil },
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "empty parse value" && tt.jvt.ParseValue != nil {
				t.Error("Expected nil ParseValue for empty case")
			}
			if tt.name == "with parse value function" && tt.jvt.ParseValue == nil {
				t.Error("Expected non-nil ParseValue for function case")
			}
		})
	}
}

func TestNewJSONParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		handler func(any) error
	}{
		{
			name:    "create parser with empty input",
			input:   "",
			handler: func(v any) error { return nil },
		},
		{
			name:    "create parser with non-empty input",
			input:   "test input",
			handler: func(v any) error { return nil },
		},
		{
			name:    "create parser with whitespace input",
			input:   "  \t\n\r",
			handler: func(v any) error { return nil },
		},
		{
			name:    "create parser with JSON object",
			input:   "{\"key\": \"value\"}",
			handler: func(v any) error { return nil },
		},
		{
			name:    "create parser with JSON array",
			input:   "[1, 2, 3]",
			handler: func(v any) error { return nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			parser, err := NewJSONParser(reader, tt.handler)
			if err != nil {
				t.Fatal(fmt.Sprintf("NewJSONParser() encountered an error: %v", err))
			}

			if parser == nil {
				t.Fatal("NewJSONParser() returned nil")
			}
			if parser.reader != reader {
				t.Error("NewJSONParser() reader not set correctly")
			}
			if parser.parseHandler == nil {
				t.Error("NewJSONParser() handler not set correctly")
			}
			if parser.pos != 0 {
				t.Errorf("NewJSONParser() pos = %v, want 0", parser.pos)
			}
			if parser.NowField != "" {
				t.Errorf("NewJSONParser() NowField = %v, want empty string", parser.NowField)
			}
			// Don't check buffer content since streamData() is called in NewJSONParser
			if len(parser.buffer) == 0 && tt.input != "" && strings.TrimSpace(tt.input) != "" {
				t.Error("NewJSONParser() buffer should not be empty for non-empty input")
			}
		})
	}
}
func TestJSONParserSetParseHandler(t *testing.T) {
	tests := []struct {
		name string
		want any
	}{
		{
			name: "set handler with integer",
			want: 42,
		},
		{
			name: "set handler with string",
			want: "test",
		},
		{
			name: "set handler with boolean",
			want: true,
		},
		{
			name: "set handler with nil",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JSONParser{
				reader: bufio.NewReader(strings.NewReader("")),
			}
			var received any
			p.SetParseHandler(func(v any) error {
				received = v
				return nil
			})
			if p.parseHandler == nil {
				t.Error("SetParseHandler() failed to set handler")
			}
			p.parseHandler(tt.want)
			if received != tt.want {
				t.Errorf("ParseHandler received = %v, want %v", received, tt.want)
			}
		})
	}
}

func TestJSONParserStreamData(t *testing.T) {
	originalChunkSize := ChunkSize
	defer func() { ChunkSize = originalChunkSize }()

	tests := []struct {
		name      string
		input     string
		expected  string
		chunkSize int
	}{
		{
			name:      "read single chunk",
			input:     "test data",
			expected:  "test data",
			chunkSize: 100,
		},
		{
			name:      "empty input",
			input:     "",
			expected:  "",
			chunkSize: 1024,
		},
		{
			name:      "large input",
			input:     "large test data that exceeds chunk size",
			expected:  "large te",
			chunkSize: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ChunkSize = tt.chunkSize
			p := &JSONParser{
				reader: bufio.NewReader(strings.NewReader(tt.input)),
				buffer: "",
			}
			p.streamData()
			if !strings.Contains(p.buffer, tt.expected) {
				t.Errorf("streamData() buffer = %v, want %v", p.buffer, tt.expected)
			}
		})
	}
}
func TestJSONParserGoForward(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		field    string
		expected string
	}{
		{
			name:     "empty initial path",
			initial:  "",
			field:    "test",
			expected: "test",
		},
		{
			name:     "append dot",
			initial:  "parent",
			field:    ".",
			expected: "parent.",
		},
		{
			name:     "append field",
			initial:  "parent.",
			field:    "child",
			expected: "parent.child",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JSONParser{NowField: tt.initial}
			p.goForward(tt.field)
			if p.NowField != tt.expected {
				t.Errorf("goForward() = %v, want %v", p.NowField, tt.expected)
			}
		})
	}
}

func TestJSONParserGoBackward(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		field    string
		expected string
	}{
		{
			name:     "remove field",
			initial:  "parent.child",
			field:    "child",
			expected: "parent.",
		},
		{
			name:     "remove dot",
			initial:  "parent.",
			field:    ".",
			expected: "parent",
		},
		{
			name:     "remove from empty",
			initial:  "test",
			field:    "test",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JSONParser{NowField: tt.initial}
			p.goBackward(tt.field)
			if p.NowField != tt.expected {
				t.Errorf("goBackward() = %v, want %v", p.NowField, tt.expected)
			}
		})
	}
}

func TestJSONParserIncrementPos(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		initialPos  int
		initialBuf  string
		expectedPos int
	}{
		{
			name:        "increment within buffer",
			input:       "",
			initialPos:  0,
			initialBuf:  "test",
			expectedPos: 1,
		},
		{
			name:        "increment at buffer end",
			input:       "more data",
			initialPos:  3,
			initialBuf:  "test",
			expectedPos: 4,
		},
		{
			name:        "increment at buffer end, no more data",
			input:       "",
			initialPos:  3,
			initialBuf:  "test",
			expectedPos: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JSONParser{
				reader: bufio.NewReader(strings.NewReader("")),
				buffer: tt.initialBuf,
				pos:    tt.initialPos,
			}
			p.incrementPos()
			if p.pos != tt.expectedPos {
				t.Errorf("incrementPos() = pos %v, want %v", p.pos, tt.expectedPos)
			}
		})
	}
}

func TestJSONParserSkipWhitespace(t *testing.T) {
	tests := []struct {
		name        string
		initialBuf  string
		initialPos  int
		expectedPos int
	}{
		{
			name:        "skip spaces",
			initialBuf:  "   test",
			initialPos:  0,
			expectedPos: 3,
		},
		{
			name:        "skip tabs and newlines",
			initialBuf:  "\t\n\rtest",
			initialPos:  0,
			expectedPos: 3,
		},
		{
			name:        "no whitespace",
			initialBuf:  "test",
			initialPos:  0,
			expectedPos: 0,
		},
		{
			name:        "all whitespace",
			initialBuf:  "   ",
			initialPos:  0,
			expectedPos: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JSONParser{
				reader: bufio.NewReader(strings.NewReader("")),
				buffer: tt.initialBuf,
				pos:    tt.initialPos,
			}
			p.skipWhitespace()
			if p.pos != tt.expectedPos {
				t.Errorf("skipWhitespace() = pos %v, want %v", p.pos, tt.expectedPos)
			}
		})
	}
}

func TestJSONParserSubtractBuffer(t *testing.T) {
	tests := []struct {
		name        string
		initialBuf  string
		initialPos  int
		expectedBuf string
		expectedPos int
	}{
		{
			name:        "subtract from middle",
			initialBuf:  "test data",
			initialPos:  4,
			expectedBuf: " data",
			expectedPos: 0,
		},
		{
			name:        "subtract from start",
			initialBuf:  "test",
			initialPos:  0,
			expectedBuf: "test",
			expectedPos: 0,
		},
		{
			name:        "subtract all",
			initialBuf:  "test",
			initialPos:  4,
			expectedBuf: "",
			expectedPos: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JSONParser{
				reader: bufio.NewReader(strings.NewReader("")),
				buffer: tt.initialBuf,
				pos:    tt.initialPos,
			}
			p.subtractBuffer()
			if p.buffer != tt.expectedBuf {
				t.Errorf("subtractBuffer() = buffer %v, want %v", p.buffer, tt.expectedBuf)
			}
			if p.pos != tt.expectedPos {
				t.Errorf("subtractBuffer() = pos %v, want %v", p.pos, tt.expectedPos)
			}
		})
	}
}

func TestJSONParserConsume(t *testing.T) {
	tests := []struct {
		name        string
		initialBuf  string
		initialPos  int
		expectedBuf string
		expectedPos int
	}{
		{
			name:        "consume whitespace and data",
			initialBuf:  "   test data",
			initialPos:  0,
			expectedBuf: "test data",
			expectedPos: 0,
		},
		{
			name:        "consume from middle with whitespace",
			initialBuf:  "test   data",
			initialPos:  4,
			expectedBuf: "data",
			expectedPos: 0,
		},
		{
			name:        "consume all whitespace",
			initialBuf:  "   \t\n\r",
			initialPos:  0,
			expectedBuf: "",
			expectedPos: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &JSONParser{
				reader: bufio.NewReader(strings.NewReader("")),
				buffer: tt.initialBuf,
				pos:    tt.initialPos,
			}
			p.consume()
			if p.buffer != tt.expectedBuf {
				t.Errorf("consume() = buffer %v, want %v", p.buffer, tt.expectedBuf)
			}
			if p.pos != tt.expectedPos {
				t.Errorf("consume() = pos %v, want %v", p.pos, tt.expectedPos)
			}
		})
	}
}

// TestJSONParserParse tests the JSON parser's ability to parse various JSON structures
// including empty objects, simple values, arrays, nested objects, and mixed types.
// It verifies that the parser correctly handles different JSON constructs and
// produces the expected sequence of values through the parse handler.
func TestJSONParserParse(t *testing.T) {
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
			name: "parse object with one key",
			input: `{
				"key": "value"  
			}`,
			expected: []any{"value", nil},
		},
		{
			name:     "parse array",
			input:    `[1, 2, 3]`,
			expected: []any{1, 2, 3, nil},
		},
		{
			name:     "parse string",
			input:    `"test string"`,
			expected: []any{"test string"},
		},
		{
			name:     "parse number",
			input:    `42.5`,
			expected: []any{42.5},
		},
		{
			name:     "parse boolean true",
			input:    `true`,
			expected: []any{true},
		},
		{
			name:     "parse boolean false",
			input:    `false`,
			expected: []any{false},
		},
		{
			name:     "parse null",
			input:    `null`,
			expected: []any{nil},
		},
		{
			name: "parse nested objects",
			input: `{
				"person": {
					"name": "Alice",
					"age": 30,
					"address": {
						"city": "Wonderland",
						"zip": "12345"
					}
				}
			}`,
			expected: []any{"Alice", 30, "Wonderland", "12345", nil, nil, nil},
		},
		{
			name: "parse array of objects",
			input: `[
				{"name": "Alice", "age": 30},
				{"name": "Bob", "age": 25}
			]`,
			expected: []any{"Alice", 30, nil, "Bob", 25, nil, nil},
		},
		{
			name: "parse mixed types in array",
			input: `[
				{"name": "Alice", "age": 30},
				42,
				"simple string",
				true,
				null
			]`,
			expected: []any{"Alice", 30, nil, 42, "simple string", true, nil, nil},
		},
		{
			name: "parse deeply nested structure",
			input: `{
				"company": {
					"name": "Tech Corp",
					"employees": [
						{"name": "Alice", "role": "Developer"},
						{"name": "Bob", "role": "Manager"}
					],
					"locations": {
						"headquarters": {
							"city": "New York",
							"zip": "10001"
						},
						"branch": {
							"city": "San Francisco",
							"zip": "94105"
						}
					}
				}
			}`,
			expected: []any{"Tech Corp", "Alice", "Developer", nil, "Bob", "Manager", nil, nil, "New York", "10001", nil, "San Francisco", "94105", nil, nil, nil, nil},
		},
		{
			name: "parse array with mixed types",
			input: `[
				{"name": "Alice", "age": 30},
				[1, 2, 3],
				"simple string",
				null,
				true
			]`,
			expected: []any{"Alice", 30, nil, 1, 2, 3, nil, "simple string", nil, true, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result = []any{}
			parser, err := NewJSONParser(bufio.NewReader(strings.NewReader(tt.input)), nil)

			if err != nil {
				t.Fatal("NewJSONParser() returned nil")
			}
			parser.SetParseHandler(func(v any) error {
				result = append(result, v)
				return nil
			})

			parser.Parse()
			// fmt.Printf("Parse() = %v, want %v", result, tt.expected)

			if !compareArray(result, tt.expected) {
				t.Errorf("Parse() = %v\n, want %v", result, tt.expected)
			}
		})
	}
}
func compareArray(result []any, expected []any) bool {
	// Check if lengths are equal
	if len(result) != len(expected) {
		fmt.Println("Arrays are of different lengths.")
		return false
	}

	// Compare each element
	for i := range result {
		if fmt.Sprintf("%v", result[i]) != fmt.Sprintf("%v", expected[i]) {
			fmt.Printf("Difference at index %d: got %v, want %v\n", i, result[i], expected[i])
			return false
		}
	}

	return true
}
