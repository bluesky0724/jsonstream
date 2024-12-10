package extractor

import (
	"bufio"
	"encoding/csv"
	"fmt"

	"github.com/bluesky0724/jsonstream/parser"
)

// JSONExtractor represents a structure for extracting JSON data and converting it to CSV format
type JSONExtractor struct {
	parser  *parser.JSONParser // JSON parser instance
	writer  *csv.Writer        // CSV writer instance
	base    string             // Base field path for target data
	targets []string           // Target field names to extract
	values  map[string][]any   // Map to store extracted values
}

// NewJSONExtractor creates a new JSONExtractor instance
func NewJSONExtractor(reader *bufio.Reader, writer *csv.Writer, baseField string, fields []string) (*JSONExtractor, error) {
	parser, err := parser.NewJSONParser(reader, nil)

	if err != nil {
		return nil, err
	}

	// Initialize map with the absolute paths of target fields
	targetValues := make(map[string][]any)
	for _, field := range fields {
		absolutePath := getAbsolutePath(baseField, field)
		targetValues[absolutePath] = []any{}
	}

	extractor := &JSONExtractor{
		parser:  parser,
		writer:  writer,
		base:    baseField,
		targets: fields,
		values:  targetValues,
	}

	// The logic to extract and export the target data is passed to parser as a parseHandler
	parser.SetParseHandler(extractor.parseHandler)

	return extractor, nil
}

// composeCSV writes the collected values to CSV and reinitializes the values map
func (e *JSONExtractor) composeCSV() error {
	if err := e.writeCSV(e.targets, e.values); err != nil {
		return fmt.Errorf("failed to write CSV: %w", err)
	}
	e.initValues()
	return nil
}

// getAbsolutePath combines base path and field name to create absolute field path
func getAbsolutePath(base string, field string) string {
	return base + "." + field
}

// parseHandler processes parsed JSON values and triggers CSV composition
// When we define this parseHandler properly, we can handle several tasks as we need
// This handler updates the current values when parsing the target fields
// and compose CSV when handling one element is finished
// JSONParser is just parsing and validating the JSON object
// and by passing this parseHandler to JSONParser we can do multiple jobs
// For example, when we need to export some data from JSON to CSV,
// and integrate with database, we can implement this logic here
// The base and targetValues are only attached in this JSONExtractor structure,
// so actually we can even define the new JSONProcessor to handle the brand new job
// just creating and passing parseHandler to JSONParser
func (e *JSONExtractor) parseHandler(value any) error {
	nowField := e.parser.NowField

	if nowField == e.base+"." { // This means the parser is parsing an element in the base array
		if err := e.composeCSV(); err != nil {
			return fmt.Errorf("failed to compose CSV: %w", err)
		}
	} else if e.shouldUpdate(nowField) { // This means the parser parsed the target field
		e.updateValues(nowField, value)
	}
	return nil
}

// initValues reinitializes the values map with empty arrays
func (e *JSONExtractor) initValues() {
	for _, field := range e.targets {
		absolutePath := getAbsolutePath(e.base, field)
		e.values[absolutePath] = []any{}
	}
}

// updateValues adds a new value to the specified field in the values map
func (e *JSONExtractor) updateValues(nowField string, value any) {
	// If we find several values that matches with the target field,
	// we can be sure of this value is in array format
	// To handle all this data, we append all values.
	e.values[nowField] = append(e.values[nowField], value)
}

// shouldUpdate checks if the current field should be updated based on target fields
func (e *JSONExtractor) shouldUpdate(field string) bool {
	for _, target := range e.targets {
		// field is the absolute path of current pointer
		// and target is the relative path of target field
		if field == getAbsolutePath(e.base, target) {
			return true
		}
	}
	return false
}

// writeCSV writes the collected values to the CSV file using backtracking
func (e *JSONExtractor) writeCSV(fields []string, values map[string][]any) error {
	absolutePaths := make([]string, len(fields))
	for i, field := range fields {
		absolutePaths[i] = getAbsolutePath(e.base, field)
	}
	return e.backtrack(absolutePaths, values, 0, []string{})
}

// backtrack generates all possible combinations of field values for CSV rows
func (e *JSONExtractor) backtrack(keys []string, obj map[string][]any, index int, current []string) error {
	if index == len(keys) {
		if err := e.writer.Write(current); err != nil {
			return fmt.Errorf("error writing field values: %w", err)
		}
		return nil
	}

	// if target field value is empty, we use ""
	if len(obj[keys[index]]) == 0 {
		obj[keys[index]] = append(obj[keys[index]], "")
	}

	for _, value := range obj[keys[index]] {
		stringValue := fmt.Sprintf("%v", value)
		current = append(current, stringValue)
		if err := e.backtrack(keys, obj, index+1, current); err != nil {
			return err
		}
		current = current[:len(current)-1]
	}
	return nil
}

// Extract starts the JSON extraction process and writes data to CSV
func (e *JSONExtractor) Extract() error {
	if err := e.writer.Write(e.targets); err != nil {
		return fmt.Errorf("error writing target fields: %w", err)
	}
	if err := e.parser.Parse(); err != nil { // Start to parse the data
		return fmt.Errorf("error parsing data: %w", err)
	}
	return nil
}
