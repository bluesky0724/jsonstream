package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var JSONNumber = &JSONValueType{}

// JSONNumber represents a JSON number value type that:
// * Begins with a digit or minus sign '-'
// * May contain digits (0-9)
// * May contain a decimal point '.'
// * May contain an exponent indicator ('e' or 'E') followed by an optional sign
// * Must be a valid numeric value parsable to float64
// * Stores a primitive numeric value
func init() {
	JSONNumber.ParseValue = func(p *JSONParser) error {
		// Check if the current character is a valid number start (digit or minus sign)
		if !unicode.IsDigit(rune(p.buffer[p.pos])) && p.buffer[p.pos] != '-' {
			return fmt.Errorf("unexpected character '%c' at position %d", p.buffer[p.pos], p.pos)
		}

		start := p.pos
		// Continue parsing while characters are valid number components (digits, signs, exponents, or decimal point)
		// loose validation check since we parse float the value later
		for p.pos < len(p.buffer) && (unicode.IsDigit(rune(p.buffer[p.pos])) || strings.ContainsRune("-+eE.", rune(p.buffer[p.pos]))) {
			p.incrementPos()
		}

		number, err := strconv.ParseFloat(p.buffer[start:p.pos], 64)
		if err != nil {
			return fmt.Errorf("invalid number")
		}

		// Handle the parsed number value
		err = p.parseHandler(number)
		if err != nil {
			return err
		}

		err = p.consume()

		if err != nil {
			return err
		}
		return nil
	}
}
