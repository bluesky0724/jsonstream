package jsonstream

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"

	"github.com/bluesky0724/jsonstream/extractor"
)

// JSON2CSV converts JSON data from a file or URL to CSV format
// Parameters:
//
//	fileType: "file" for local files or "url" for online data
//	input: file path or URL of the JSON data, e.g: "https://open.gsa.gov/data.json"
//	output: destination CSV filename, e.g: "output.csv"
//	base: base field path where target data is stored: we assume this field is an array, e.g: ".dataset"
//	fields: array of field names to extract relative to base path, e.g: "modified" means the absolute path of target field is ".dataset.modified"
func JSON2CSV(fileType string, input string, output string, base string, fields []string) error {
	var reader *bufio.Reader

	// Handle local file input
	if fileType == "file" {
		file, err := os.Open(input)
		if err != nil {

			return fmt.Errorf("error opening file: %w", err)
		}
		defer file.Close()

		reader = bufio.NewReader(file)
	} else if fileType == "url" {
		// Handle URL input
		jsonURL := input

		resp, err := http.Get(jsonURL)
		if err != nil {

			return fmt.Errorf("error fetching URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {

			return fmt.Errorf("failed to fetch data: HTTP status %d", resp.StatusCode)
		}

		reader = bufio.NewReader(resp.Body)
	} else {
		return fmt.Errorf("invalid fileType: must be 'file' or 'url'")
	}

	// Create output CSV file
	csvFilename := output

	csvFile, err := os.Create(csvFilename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer csvFile.Close()

	// Initialize CSV writer
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Create and run the JSON extractor
	extractor, err := extractor.NewJSONExtractor(reader, writer, base, fields)

	if err := extractor.Extract(); err != nil {
		return fmt.Errorf("error extracting JSON: %w", err)
	}

	return nil
}
