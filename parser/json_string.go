package parser

// JSONString represents a JSON string value type that:
// * Begins with an opening double quote '"'
// * Ends with a closing double quote '"'
// * May contain any Unicode characters
// * May contain escape sequences like \", \\, \n, \r, \t, etc.
// * Stores the raw string value between the quotes
var JSONString = &JSONValueType{}

func init() {
	JSONString.ParseValue = func(p *JSONParser) {
		// Skip opening quote '"'
		p.incrementPos()

		// Track start position of string content to extract the string
		start := p.pos
		// Continue until closing quote is found
		for p.buffer[p.pos] != '"' {
			// Handle escape sequences
			// By ignoring \ in the string. we can extract the correct string which possibly include \"
			if p.buffer[p.pos] == '\\' {
				p.incrementPos() // ignore \
			}
			p.incrementPos() // check next byte
		}
		result := p.buffer[start:p.pos]

		// Process the parsed string value
		p.parseHandler(result)

		// Skip closing quote and consume any whitespace
		p.incrementPos()
		p.consume()
	}
}
