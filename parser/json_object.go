package parser

var JSONObject = &JSONValueType{}

// JSONObject represents a JSON object value type that:
// * Begins with an opening curly brace '{'
// * Ends with a closing curly brace '}'
// * May contain zero key-value pairs (empty object)
// * May contain one or multiple key-value pairs separated by commas ','
// * Each key must be a string enclosed in double quotes
// * Each value must be a valid JSON value (string, number, boolean, null, object, array)
// * Does not store a primitive value itself, as it's a container type
func init() {
	JSONObject.ParseValue = func(p *JSONParser) {
		// Move past opening brace '{'
		p.incrementPos()
		p.consume()

		for {
			// Check for end of object: for the empty object
			if p.buffer[p.pos] == '}' {
				p.incrementPos()
				p.consume()
				break
			}

			// Parse the key string and ":"
			key := parseKey(p)

			// Navigate to correct path and parse value
			p.goForward(key)
			p.Parse()
			p.goBackward(key)

			// Handle comma separator or end of object
			if p.buffer[p.pos] == ',' {
				p.incrementPos()
				p.consume()
			} else if p.buffer[p.pos] == '}' {
				p.incrementPos()
				p.consume()
				break
			} else {
				panic("expected ',' or '}' to be a valid object")
			}
		}

		// Call parse handler with nil value
		// Objects and Arrays are considered as they have no significant data
		// The target data is always the primitive values like string, number, boolean and null
		p.parseHandler(nil)

	}
}

// parseKey parses a JSON object key which must be a string
func parseKey(p *JSONParser) string {
	// Ensure key starts with a quote
	if p.buffer[p.pos] != '"' {
		panic("Expected string for the object key")
	}
	p.incrementPos()

	// Find the end of the string
	start := p.pos
	for p.buffer[p.pos] != '"' {
		// Handle escaped characters
		if p.buffer[p.pos] == '\\' {
			p.incrementPos()
		}
		p.incrementPos()
	}
	result := p.buffer[start:p.pos]

	// Move past closing quote and whitespace
	p.incrementPos()
	p.consume()

	// Ensure key is followed by colon
	if p.buffer[p.pos] != ':' {
		panic("expected ':' after key string")
	}

	// Move past colon and whitespace
	p.incrementPos()
	p.consume()

	return result
}
