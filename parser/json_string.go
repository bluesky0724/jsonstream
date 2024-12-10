package parser

// JSONString represents a JSON string value type that:
// * Begins with an opening double quote '"'
// * Ends with a closing double quote '"'
// * May contain any Unicode characters
// * May contain escape sequences like \", \\, \n, \r, \t, etc.
// * Stores the raw string value between the quotes
var JSONString = &JSONValueType{}

func init() {
	JSONString.ParseValue = func(p *JSONParser) error {
		// Skip opening quote '"'
		if err := p.incrementPos(); err != nil {
			return err
		}

		// Track start position of string content to extract the string
		start := p.pos
		// Continue until closing quote is found
		for p.buffer[p.pos] != '"' {
			// Handle escape sequences
			if p.buffer[p.pos] == '\\' {
				if err := p.incrementPos(); err != nil {
					return err
				}
			}
			if err := p.incrementPos(); err != nil {
				return err
			}
		}
		result := p.buffer[start:p.pos]

		// Process the parsed string value
		if err := p.parseHandler(result); err != nil {
			return err
		}

		// Skip closing quote and consume any whitespace
		if err := p.incrementPos(); err != nil {
			return err
		}
		if err := p.consume(); err != nil {
			return err
		}

		return nil
	}
}
