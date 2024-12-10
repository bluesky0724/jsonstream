package parser

import "fmt"

var JSONArray = &JSONValueType{}

// JSONArray represents a JSON array value type that:
// * Begins with an opening square bracket '['
// * Ends with a closing square bracket ']'
// * May contain zero elements (empty array)
// * May contain multiple elements separated by commas ','
// * Each element must be a valid JSON value
// * Does not store a primitive value itself, as it's a container type
func init() {
	JSONArray.ParseValue = func(p *JSONParser) error {
		// Move past the opening '['
		// No validation logic 'coz we only call when initializer matches
		if p.pos >= len(p.buffer) {
			return fmt.Errorf("unexpected end of input while parsing array")
		}
		p.incrementPos()
		p.consume()

		for {
			// Check if we've reached the end of the array
			// Need this logic for the empty array
			if p.pos >= len(p.buffer) {
				return fmt.Errorf("unexpected end of input: array was not closed")
			}
			if p.buffer[p.pos] == ']' {
				p.incrementPos()
				p.consume()
				break
			}

			// Parse the next value in the array
			if err := p.Parse(); err != nil {
				return err
			}

			if p.pos >= len(p.buffer) {
				return fmt.Errorf("unexpected end of input: array was not closed")
			}
			if p.buffer[p.pos] == ',' {
				p.incrementPos()
				p.consume()
			} else if p.buffer[p.pos] == ']' {
				p.incrementPos()
				p.consume()
				break
			} else {
				return fmt.Errorf("expected ',' or ']' to be a valid array")
			}
		}

		// Call parse handler with nil value since array has no value
		return p.parseHandler(nil)
	}
}
