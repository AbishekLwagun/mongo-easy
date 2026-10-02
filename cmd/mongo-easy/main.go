package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/abisheklwagun/mongo-easy/internal/cli"
	"github.com/abisheklwagun/mongo-easy/internal/importer"
	"github.com/abisheklwagun/mongo-easy/internal/mongodb"
	"github.com/abisheklwagun/mongo-easy/internal/ui"
)

func main() {
	// Show the help message when no command or a help flag is provided.
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		cli.PrintHelp()
		return
	}

	// Mongo Easy currently supports the import command.
	if os.Args[1] != "import" {
		fmt.Println("Error: Unknown command", os.Args[1])
		return
	}

	// The import command needs a file or ZIP path to work.
	if len(os.Args) < 3 {
		fmt.Println("Error: Dataset path is required")
		return
	}

	path := os.Args[2]

	// Check the path before trying to connect to MongoDB or read the file.
	if err := validateDataset(path); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// The remaining arguments are the options for the import command.
	args := os.Args[3:]

	importMode, err := cli.ParseImportMode(args)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	collectionMode, err := cli.ParseCollectionMode(args)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	upsertKey := cli.ParseUpsertKey(args)

	databaseName, err := cli.ParseDatabase(args)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Connect to MongoDB before starting the import.
	client, err := mongodb.Connect()
	if err != nil {
		fmt.Println("Error: Could not connect to MongoDB:", err)
		return
	}
	defer client.Disconnect(context.Background())

	// Ping MongoDB so we know the server is actually reachable.
	if err := client.Ping(context.Background(), nil); err != nil {
		fmt.Println("Error: Could not reach MongoDB:", err)
		return
	}

	ui.Header()

	// Show the upsert key when upsert mode is being used.
	// If no key was provided, the importer will try to find one automatically.
	if importMode == cli.ImportModeUpsert {
		if upsertKey != "" {
			fmt.Println("Upsert key:", upsertKey)
		} else {
			fmt.Println("Upsert key: automatic detection")
		}
	}

	database := client.Database(databaseName)

	// Keep the import settings together so the importer does not need
	// to know anything about command-line arguments.
	options := importer.Options{
		Mode:           importMode,
		CollectionMode: collectionMode,
		UpsertKey:      upsertKey,
	}

	// The importer handles the actual parsing and MongoDB import.
	// The file extension tells it which type of file it is dealing with.
	imported, err := importer.Run(
		database,
		path,
		filepath.Ext(path),
		options,
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Show the final message only when at least one collection was imported.
	if imported {
		ui.Done()
	} else {
		ui.Warning("No collections were imported.")
	}

}

// validateDataset checks the input before the import starts.
// Keeping this check here gives the user a simple error instead of
// failing later while trying to open or parse the file.
func validateDataset(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("File not found: %s", path)
	}

	if info.IsDir() {
		return fmt.Errorf("Dataset is a directory %s", path)
	}

	// These are the file formats currently supported by Mongo Easy.
	switch filepath.Ext(path) {
	case ".zip", ".json", ".jsonl", ".csv":
		return nil
	default:
		return fmt.Errorf("Unsupported file type: %s", filepath.Ext(path))
	}
}
