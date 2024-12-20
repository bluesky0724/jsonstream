package parser

import "fmt"

// This file holds the strict JSON types like JSONTrue, JSONFalse and JSONNull

var JSONTrue = &JSONValueType{}
var JSONFalse = &JSONValueType{}
var JSONNull = &JSONValueType{}

// init initializes the parse functions for JSON boolean and null values
func init() {
	// ParseValue verifies the input file and then returns the strict value
	JSONTrue.ParseValue = func(p *JSONParser) error {
		err := strictCheck(p, "true")
		if err != nil {
			return err
		}

		return p.parseHandler(true)
	}

	JSONFalse.ParseValue = func(p *JSONParser) error {
		err := strictCheck(p, "false")
		if err != nil {
			return err
		}

		return p.parseHandler(false)
	}

	JSONNull.ParseValue = func(p *JSONParser) error {
		err := strictCheck(p, "null")
		if err != nil {
			return err
		}

		return p.parseHandler(nil)
	}
}

// strictCheck verifies that the input matches the expected string exactly
// It panics if there is a mismatch
func strictCheck(p *JSONParser, expected string) error {
	position := 0
	for position < len(expected) {
		if p.pos >= len(p.buffer) || p.buffer[p.pos] != expected[position] {
			panic(fmt.Sprintf("expected '%s'", expected))
		}
		if err := p.incrementPos(); err != nil { // If buffer limit is reached, load more data
			return err
		}
		position++
	}
	return p.consume()
}
