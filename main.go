package main

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
func JSON2CSV(fileType string, input string, output string, base string, fields []string) {
	var reader *bufio.Reader

	// Handle local file input
	if fileType == "file" {
		file, err := os.Open(input)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		reader = bufio.NewReader(file)
	} else if fileType == "url" {
		// Handle URL input
		jsonURL := input

		resp, err := http.Get(jsonURL)
		if err != nil {
			fmt.Println("Error fetching URL:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to fetch data: %d\n", resp.StatusCode)
			return
		}

		reader = bufio.NewReader(resp.Body)
	}

	// Create output CSV file
	csvFilename := output
	csvFile, csvErr := os.Create(csvFilename)
	if csvErr != nil {
		fmt.Println("Error creating file: ", csvErr)
		return
	}
	defer csvFile.Close()

	// Initialize CSV writer
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Create and run the JSON extractor
	extractor := extractor.NewJSONExtractor(reader, writer, base, fields)
	extractor.Extract()
}

func main() {
	// JSON2CSV("file", "data-8.json", "result-8.csv", ".dataset", []string{"modified", "publisher.name", "publisher.subOrganizationOf.name", "contactPoint.fn", "keyword"})
	JSON2CSV("url", "https://open.gsa.gov/data.json", "result.csv", ".dataset", []string{"modified", "publisher.name", "publisher.subOrganizationOf.name", "contactPoint.fn", "keyword"})

}
