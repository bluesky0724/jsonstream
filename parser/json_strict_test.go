package parser

import (
	"testing"
)

func TestJSONTrue(t *testing.T) {
	p := &JSONParser{reader: nil, buffer: ("true"), pos: 0, NowField: "", parseHandler: func(v interface{}) {
		if v != true {
			t.Errorf("Expected true, got %v", v)
		}
	}}
	JSONTrue.ParseValue(p)
}

func TestJSONFalse(t *testing.T) {
	p := &JSONParser{buffer: ("false"), parseHandler: func(v interface{}) {
		if v != false {
			t.Errorf("Expected false, got %v", v)
		}
	}}
	JSONFalse.ParseValue(p)
}

func TestJSONNull(t *testing.T) {
	p := &JSONParser{buffer: ("null"), parseHandler: func(v interface{}) {
		if v != nil {
			t.Errorf("Expected nil, got %v", v)
		}
	}}
	JSONNull.ParseValue(p)
}

func TestStrictCheckPanic(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"tru", "true"},
		{"fals", "false"},
		{"nul", "null"},
		{"truex", "true"},
		{"falsex", "false"},
		{"nullx", "null"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for input %s", tc.input)
				}
			}()
			p := &JSONParser{buffer: (tc.input)}
			strictCheck(p, tc.expected)
		})
	}
}
