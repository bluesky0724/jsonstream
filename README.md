# JSON CSV Extractor

## Purpose

This module was developed to extract necessary data from an online JSON file and output it in CSV file format.

## Assumptions

1. File Processing

   The corresponding JSON file can be several gigabytes in size. Therefore, file reading was done by streaming.
2. Data Processing

   Although the JSON object is very large, we assume that a unique group of field values ​​that extract only the field values ​​we want can be processed in memory. (< 2 gigabytes)

3. Target field value processing
   The target field values may or may not exist in the JSON file. We assume that:

   - If a target field value is empty or stored as an object, it will be considered empty.
   - If the value is a string, boolean, number, or null, it will be returned as is.
   - If the value is an array, each element will be printed in a separate row.


## How to use the module

### Install the module

```bash
go get github.com/bluesky0724/jsonstream
```

### JSON2CSV function explanation

```go
func JSON2CSV(fileType string, input string, output string, base string, fields []string)
```

JSON2CSV converts JSON data from a file or URL to CSV format
Parameters:

- fileType: "file" for local files or "url" for online data

- input: file path or URL of the JSON data, e.g: "https://open.gsa.gov/data.json"

- output: destination CSV filename, e.g: "output.csv"

- base: base field path where target data is stored: we assume this field is an array, e.g: ".dataset"

- fields: array of field names to extract relative to base path, e.g: "modified" means the absolute path of target field is ".dataset.modified"

Example:
- Local JSON file

```Go
import "github.com/bluesky0724/jsonstream/extractor"

func main() {
	JSON2CSV("file", "data-8.json", "result-8.csv", ".dataset", []string{"modified", "publisher.name", "publisher.subOrganizationOf.name", "contactPoint.fn", "keyword"})
}
```

- Online JSON file

```Go
import "github.com/bluesky0724/jsonstream/extractor"

func main() {
	JSON2CSV("url", "https://open.gsa.gov/data.json", "result.csv", ".dataset", []string{"modified", "publisher.name", "publisher.subOrganizationOf.name", "contactPoint.fn", "keyword"})
}
```

### Test the project

```bash
go test .
```

