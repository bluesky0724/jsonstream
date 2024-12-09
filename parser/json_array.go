package parser

var JSONArray = &JSONValueType{}

// JSONArray represents a JSON array value type that:
// * Begins with an opening square bracket '['
// * Ends with a closing square bracket ']'
// * May contain zero elements (empty array)
// * May contain multiple elements separated by commas ','
// * Each element must be a valid JSON value
// * Does not store a primitive value itself, as it's a container type

func init() {
	JSONArray.ParseValue = func(p *JSONParser) {
		// Move past the opening '['
		// No validation logic 'coz we only call when initializer matches
		p.incrementPos()
		p.consume()

		for {
			// Check if we've reached the end of the array
			// Need this logic for the empty array
			if p.buffer[p.pos] == ']' {
				p.incrementPos()
				p.consume()
				break
			}

			// Parse the next value in the array
			p.Parse()

			if p.buffer[p.pos] == ',' {
				p.incrementPos()
				p.consume()
			} else if p.buffer[p.pos] == ']' {
				p.incrementPos()
				p.consume()
				break
			} else {
				panic("expected ',' or ']' to be a valid array")
			}
		}

		// Call parse handler with nil value since array has no value
		p.parseHandler(nil)
	}
}
