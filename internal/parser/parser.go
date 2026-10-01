package parser

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"path/filepath"
)

type Document map[string]interface{}

// IsSupportedDataFile checks whether the file is a format Mongo Easy can import.
func IsSupportedDataFile(name string) bool {
	switch filepath.Ext(name) {
	case ".json", ".jsonl", ".csv":
		return true
	default:
		return false
	}
}

// CollectionName uses the file name without its extension as the collection name.
func CollectionName(filename string) string {
	baseName := filepath.Base(filename)
	extension := filepath.Ext(baseName)

	return baseName[:len(baseName)-len(extension)]
}

// Parse sends the file data to the parser for its specific format.
func Parse(data []byte, extension string) ([]Document, error) {
	switch extension {
	case ".json":
		return ParseJSON(data)
	case ".jsonl":
		return ParseJSONL(data)
	case ".csv":
		return ParseCSV(data)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", extension)
	}
}

// ParseJSON handles both a single JSON document and an array of documents.
func ParseJSON(data []byte) ([]Document, error) {
	var value interface{}

	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}

	switch value := value.(type) {
	case map[string]interface{}:
		return []Document{Document(value)}, nil

	case []interface{}:
		documents := make([]Document, 0, len(value))

		for _, item := range value {
			document, ok := item.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("array contains a non-document value")
			}

			documents = append(documents, Document(document))
		}

		return documents, nil

	default:
		return nil, fmt.Errorf("unsupported JSON structure")
	}
}

// ParseJSONL reads one JSON document from each non-empty line.
func ParseJSONL(data []byte) ([]Document, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var documents []Document

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())

		if len(line) == 0 {
			continue
		}

		var document Document
		if err := json.Unmarshal(line, &document); err != nil {
			return nil, err
		}

		documents = append(documents, document)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}

// ParseCSV uses the first row as the field names for each document.
func ParseCSV(data []byte) ([]Document, error) {
	reader := csv.NewReader(bytes.NewReader(data))

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("CSV must contain a header and at least one row")
	}

	headers := rows[0]
	documents := make([]Document, 0, len(rows)-1)

	for _, row := range rows[1:] {
		if len(row) != len(headers) {
			return nil, fmt.Errorf("CSV row has different number of columns")
		}

		document := Document{}

		for i, header := range headers {
			document[header] = row[i]
		}

		documents = append(documents, document)
	}

	return documents, nil
}
