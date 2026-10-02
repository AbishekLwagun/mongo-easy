package importer

import (
	"fmt"
	"os"

	"github.com/abisheklwagun/mongo-easy/internal/cli"
	"github.com/abisheklwagun/mongo-easy/internal/mongodb"
	"github.com/abisheklwagun/mongo-easy/internal/parser"
	"github.com/abisheklwagun/mongo-easy/internal/ui"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Options struct {
	Mode           cli.ImportMode
	CollectionMode cli.CollectionMode
	UpsertKey      string
}

func Run(database *mongo.Database, path, extension string, options Options) (bool, error) {
	// ZIP files need to be handled differently because they can contain multiple files.
	if extension == ".zip" {
		return importZIP(database, path, options)
	}

	// Normal files are handled by the regular file importer.
	return importFile(database, path, extension, options)
}
func importFile(database *mongo.Database, path, extension string, options Options) (bool, error) {
	// Read the dataset file.
	data, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("could not read file: %s", path)
	}

	// Parse the file into MongoDB documents.
	documents, err := parser.Parse(data, extension)
	if err != nil {
		return false, err
	}

	// Use the file name as the collection name.
	collectionName := parser.CollectionName(path)
	ui.Info(fmt.Sprintf("Collection: %s", collectionName))
	ui.Info(fmt.Sprintf("Documents ready for import: %d", len(documents)))

	// Import the documents and return whether the import happened.
	imported, err := importCollection(database, collectionName, documents, options)
	return imported, err
}

func importCollection(
	database *mongo.Database,
	collectionName string,
	documents []parser.Document,
	options Options,
) (bool, error) {
	if options.Mode == cli.ImportModeAsk {
		// Check the collection first so we do not add data without the user choosing a mode.
		hasDocuments, err := mongodb.HasDocuments(database, collectionName)
		if err != nil {
			return false, fmt.Errorf("could not check collection: %w", err)
		}

		if hasDocuments {
			ui.Warning(fmt.Sprintf("Collection already contains data: %s", collectionName))
			fmt.Println("Import cancelled.")
			fmt.Println("Explicit import mode required:")
			fmt.Println("  --mode append")
			fmt.Println("  --mode replace")
			fmt.Println("  --mode upsert")
			return false, nil
		}
	}

	switch options.Mode {
	case cli.ImportModeReplace:
		ui.Info(fmt.Sprintf("Replacing collection: %s", collectionName))

		if err := mongodb.ReplaceCollection(database, collectionName, documents); err != nil {
			return false, fmt.Errorf("could not import documents: %w", err)
		}

		ui.Success(fmt.Sprintf("Imported %d documents", len(documents)))

	case cli.ImportModeUpsert:
		ui.Info(fmt.Sprintf("Upserting collection: %s", collectionName))

		key := options.UpsertKey
		if key == "" {
			// If no key was given, try to find a usable key from the documents.
			var err error
			key, err = DetectUpsertKey(documents)
			if err != nil {
				return false, err
			}
		}

		ui.Info(fmt.Sprintf("Upsert key: %s", key))

		spinner := ui.StartSpinner(fmt.Sprintf("Processing %d documents", len(documents)))

		inserted, updated, matched, err := UpsertDocuments(
			database,
			collectionName,
			documents,
			key,
		)
		if err != nil {
			spinner.StopError("Import failed")
			return false, fmt.Errorf("could not import documents: %w", err)
		}

		spinner.StopSuccess(
			fmt.Sprintf("Inserted %d, updated %d, already existed %d", inserted, updated, matched),
		)

	case cli.ImportModeAppend:
		ui.Info(fmt.Sprintf("Appending to collection: %s", collectionName))

		if err := mongodb.InsertDocuments(database, collectionName, documents); err != nil {
			return false, fmt.Errorf("could not import documents: %w", err)
		}

		ui.Success(fmt.Sprintf("Appended %d documents", len(documents)))

	case cli.ImportModeAsk:
		if err := mongodb.InsertDocuments(database, collectionName, documents); err != nil {
			return false, fmt.Errorf("could not import documents: %w", err)
		}

		ui.Success(fmt.Sprintf("Imported %d documents", len(documents)))

	default:
		return false, fmt.Errorf("unsupported import mode: %s", options.Mode)
	}

	return true, nil
}
