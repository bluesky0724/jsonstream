package parser

import (
	"bufio"
	"errors"
	"unicode"
)

// ChunkSize defines the size of data chunks to read: 1kb by default
const ChunkSize = 1024

// JSONParser represents a JSON parser with buffered reading capabilities
type JSONParser struct {
	reader       *bufio.Reader
	buffer       string
	pos          int             // the position of the parser pointer
	NowField     string          // the current field parser is checking
	parseHandler func(any) error // the logic the parser handles after parsing
}

// JSONValueType defines a type to check in JSON format
type JSONValueType struct {
	// In JSON file, there are multiple values
	// The object, array, string, number
	// and the strict values like true, false, and null
	// The parsing logic for different value types here
	ParseValue func(p *JSONParser) error
}

// NewJSONParser creates a new JSON parser instance
func NewJSONParser(reader *bufio.Reader, parseHandler func(any) error) *JSONParser {
	return &JSONParser{
		reader:       reader,
		buffer:       "", // initially empty string
		pos:          0,  // the position of the pointer is 0
		NowField:     "", // no field is detected in the beginning
		parseHandler: parseHandler,
	}
}

// SetParseHandler sets the parse handler function, made to set this private parseHandler
func (p *JSONParser) SetParseHandler(parseHandler func(any) error) {
	p.parseHandler = parseHandler
}

// StreamData reads data chunks from the reader into the buffer
func (p *JSONParser) StreamData() error {
	chunk := make([]byte, ChunkSize) // stream data by chunk size

	n, err := p.reader.Read(chunk)
	if n > 0 {
		p.buffer += string(chunk[:n])
	}

	if err != nil {
		if err.Error() == "EOF" {
			return nil // End of file reached
		}
		return errors.New("error loading more data")
	}
	return nil
}

// goForward appends a field or "." to the current field path
func (p *JSONParser) goForward(field string) {
	p.NowField += field
}

// goBackward removes a field from the current field path
func (p *JSONParser) goBackward(field string) {
	n := len(p.NowField) - len(field)
	p.NowField = p.NowField[:n]
}

// Parse is the main function to parse the JSON data
func (p *JSONParser) Parse() error {
	p.consume() // Skip whitespace for the first time..

	// Determine the JSONValue type by comparing the initializer with the current buffer
	switch p.buffer[p.pos] {
	// The JSONObject and JSONArray are composite types
	// ParseValue function has no result return but calls the main Parse function inside
	case '{':
		// Parsing object: append "." and remove it before and after parsing
		p.goForward(".")
		if err := JSONObject.ParseValue(p); err != nil {
			return err
		}
		p.goBackward(".")
	case '[':
		if err := JSONArray.ParseValue(p); err != nil {
			return err
		}
	// The other types are primitive types
	// These ParseValue functions only move the pointer and call parseHandler with the result taken
	case '"':
		if err := JSONString.ParseValue(p); err != nil {
			return err
		}
	case 't':
		if err := JSONTrue.ParseValue(p); err != nil {
			return err
		}
	case 'f':
		if err := JSONFalse.ParseValue(p); err != nil {
			return err
		}
	case 'n':
		if err := JSONNull.ParseValue(p); err != nil {
			return err
		}
	default: // If no initializer is matching, we can assume the value is number
		if err := JSONNumber.ParseValue(p); err != nil {
			return err
		}
	}
	return nil
}

// In conclusion, if this Parse function is called with the NowField as an empty string (""),
// it assumes that the pointer is at the beginning of a JSONObject or JSONArray and will
// continue parsing until the end of this composite data is reached
// To use this Parse function properly, you should validate the JSON file has the valid JSONObject
// or JSONArray in the beginning (e.g. Just a simple '[]' in the beginning will finish the process)
// and should initialize the JSONParser with NewJSONParser function.

// incrementPos increments the buffer position and loads more data if needed
func (p *JSONParser) incrementPos() error {
	p.pos++
	if p.pos >= len(p.buffer) {
		if err := p.StreamData(); err != nil { // If buffer limit is reached, load more data
			return err
		}
	}
	// After calling this function, when we find that p.pos is smaller than buffer length
	// We can ensure that the parser reached the end of the JSON file
	return nil
}

// skipWhitespace skips any whitespace characters in the buffer
func (p *JSONParser) skipWhitespace() error {
	for p.pos < len(p.buffer) && unicode.IsSpace(rune(p.buffer[p.pos])) {
		if err := p.incrementPos(); err != nil {
			return err
		}
	}
	return nil
}

// subtractBuffer removes processed data from the buffer
func (p *JSONParser) subtractBuffer() {
	p.buffer = p.buffer[p.pos:]
	p.pos = 0
}

// consume skips whitespace and removes processed data from the buffer
// This function is called whenever the parsing is finished
// After processing the value (string, number...) and the symbol ('{}' or ':'...)
// skip all white spaces and removes all processed data
func (p *JSONParser) consume() error {
	if err := p.skipWhitespace(); err != nil {
		return err
	}
	p.subtractBuffer()
	return nil
}
