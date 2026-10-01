package cli

import "fmt"

type ImportMode string

const (
	ImportModeAsk     ImportMode = "ask"
	ImportModeAppend  ImportMode = "append"
	ImportModeReplace ImportMode = "replace"
	ImportModeUpsert  ImportMode = "upsert"
)

type CollectionMode string

const (
	CollectionModeAsk      CollectionMode = "ask"
	CollectionModeCombine  CollectionMode = "combine"
	CollectionModeSeparate CollectionMode = "separate"
	CollectionModeCancel   CollectionMode = "cancel"
)

func ParseImportMode(args []string) (ImportMode, error) {
	mode := ImportModeAsk

	// Parse the import mode from the command-line arguments.
	for i := 0; i < len(args); i++ {
		if args[i] != "--mode" {
			continue
		}

		if i+1 >= len(args) {
			return "", fmt.Errorf("--mode requires a value")
		}

		switch args[i+1] {
		case "ask":
			mode = ImportModeAsk
		case "append":
			mode = ImportModeAppend
		case "replace":
			mode = ImportModeReplace
		case "upsert":
			mode = ImportModeUpsert
		default:
			return "", fmt.Errorf("unsupported import mode: %s", args[i+1])
		}

		i++
	}

	return mode, nil
}

func ParseCollectionMode(args []string) (CollectionMode, error) {
	mode := CollectionModeAsk

	// Parse how duplicate collections in a ZIP should be handled.
	for i := 0; i < len(args); i++ {
		if args[i] != "--collection-mode" {
			continue
		}

		if i+1 >= len(args) {
			return "", fmt.Errorf("--collection-mode requires a value")
		}

		switch args[i+1] {
		case "ask":
			mode = CollectionModeAsk
		case "combine":
			mode = CollectionModeCombine
		case "separate":
			mode = CollectionModeSeparate
		case "cancel":
			mode = CollectionModeCancel
		default:
			return "", fmt.Errorf("unsupported collection mode: %s", args[i+1])
		}

		i++
	}

	return mode, nil
}

func ParseUpsertKey(args []string) string {
	// Use the key provided with --key, if one was specified.
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "--key" {
			return args[i+1]
		}
	}

	return ""
}
