package parser

import "fmt"

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
	JSONObject.ParseValue = func(p *JSONParser) error {
		// Move past opening brace '{'
		if err := p.incrementPos(); err != nil {
			return err
		}
		if err := p.consume(); err != nil {
			return err
		}

		for {
			// Check for end of object: for the empty object
			if p.buffer[p.pos] == '}' {
				if err := p.incrementPos(); err != nil {
					return err
				}
				if err := p.consume(); err != nil {
					return err
				}
				break
			}

			// Parse the key string and ":"
			key, err := parseKey(p)
			if err != nil {
				return err
			}

			// Navigate to correct path and parse value
			p.goForward(key)
			if err := p.Parse(); err != nil {
				return err
			}
			p.goBackward(key)

			// Handle comma separator or end of object
			if p.buffer[p.pos] == ',' {
				if err := p.incrementPos(); err != nil {
					return err
				}
				if err := p.consume(); err != nil {
					return err
				}
			} else if p.buffer[p.pos] == '}' {
				if err := p.incrementPos(); err != nil {
					return err
				}
				if err := p.consume(); err != nil {
					return err
				}
				break
			} else {
				return fmt.Errorf("expected ',' or '}' to be a valid object")
			}
		}

		// Call parse handler with nil value
		// Objects and Arrays are considered as they have no significant data
		// The target data is always the primitive values like string, number, boolean and null
		if err := p.parseHandler(nil); err != nil {
			return err
		}
		return nil
	}
}

// parseKey parses a JSON object key which must be a string
func parseKey(p *JSONParser) (string, error) {
	// Ensure key starts with a quote
	if p.buffer[p.pos] != '"' {
		return "", fmt.Errorf("expected string for the object key")
	}
	if err := p.incrementPos(); err != nil {
		return "", err
	}

	// Find the end of the string
	start := p.pos
	for p.buffer[p.pos] != '"' {
		// Handle escaped characters
		if p.buffer[p.pos] == '\\' {
			if err := p.incrementPos(); err != nil {
				return "", err
			}
		}
		if err := p.incrementPos(); err != nil {
			return "", err
		}
	}
	result := p.buffer[start:p.pos]

	// Move past closing quote and whitespace
	if err := p.incrementPos(); err != nil {
		return "", err
	}
	if err := p.consume(); err != nil {
		return "", err
	}

	// Ensure key is followed by colon
	if p.buffer[p.pos] != ':' {
		return "", fmt.Errorf("expected ':' after key string")
	}

	// Move past colon and whitespace
	if err := p.incrementPos(); err != nil {
		return "", err
	}
	if err := p.consume(); err != nil {
		return "", err
	}

	return result, nil
}
