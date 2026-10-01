package importer

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/abisheklwagun/mongo-easy/internal/cli"
	"github.com/abisheklwagun/mongo-easy/internal/parser"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func importZIP(database *mongo.Database, path string, options Options) error {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Errorf("could not open ZIP file: %w", err)
	}
	defer reader.Close()

	collections := make(map[string][]parser.Document)
	cancelled := make(map[string]bool)
	fileCount := 0

	fmt.Println("Files in archive:")

	for _, file := range reader.File {
		if !parser.IsSupportedDataFile(file.Name) {
			continue
		}

		collection := parser.CollectionName(file.Name)
		extension := filepath.Ext(file.Name)

		content, err := file.Open()
		if err != nil {
			return fmt.Errorf("could not open file: %s", file.Name)
		}

		data, readErr := io.ReadAll(content)
		closeErr := content.Close()
		if readErr != nil {
			return fmt.Errorf("could not read file: %s", file.Name)
		}
		if closeErr != nil {
			return fmt.Errorf("could not close file: %s", file.Name)
		}

		documents, err := parser.Parse(data, extension)
		if err != nil {
			return fmt.Errorf("%s: %w", file.Name, err)
		}

		fileCount++

		if existing, exists := collections[collection]; exists {
			// Decide what to do when multiple files use the same collection name.
			switch options.CollectionMode {
			case cli.CollectionModeCombine:
				// Keep the original collection and add the new documents below.

			case cli.CollectionModeCancel:
				fmt.Println("Import cancelled for collection:", collection)
				cancelled[collection] = true
				continue

			case cli.CollectionModeAsk:
				fmt.Println("Import cancelled.")
				fmt.Println("Explicit collection mode required:")
				fmt.Println("  --collection-mode combine")
				fmt.Println("  --collection-mode separate")
				fmt.Println("  --collection-mode cancel")
				cancelled[collection] = true
				continue

			case cli.CollectionModeSeparate:
				separateCollection := collection + "_" + strings.TrimPrefix(extension, ".")
				fmt.Println("Creating separate collection:", separateCollection)
				collections[separateCollection] = append(
					collections[separateCollection],
					documents...,
				)
				continue

			default:
				return fmt.Errorf("unsupported collection mode: %s", options.CollectionMode)
			}

			_ = existing
		}

		collections[collection] = append(collections[collection], documents...)
	}

	fmt.Println("Total files:", fileCount)

	for collectionName, documents := range collections {
		if cancelled[collectionName] {
			fmt.Println("Skipping cancelled collection:", collectionName)
			continue
		}

		if err := importCollection(database, collectionName, documents, options); err != nil {
			fmt.Println("Error:", err)
		}
	}

	return nil
}
