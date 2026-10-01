package importer

import (
	"context"
	"fmt"

	"github.com/abisheklwagun/mongo-easy/internal/parser"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DetectUpsertKey(documents []parser.Document) (string, error) {
	candidates := []string{
		"_id",
		"id",
		"customer_id",
		"product_id",
		"user_id",
		"sku",
		"email",
	}

	// Use the first field that exists in every document and has unique values.
	for _, key := range candidates {
		seen := make(map[string]bool)
		valid := true

		for _, document := range documents {
			value, exists := document[key]
			if !exists {
				valid = false
				break
			}

			valueString := fmt.Sprintf("%v", value)
			if seen[valueString] {
				valid = false
				break
			}

			seen[valueString] = true
		}

		if valid {
			return key, nil
		}
	}

	return "", fmt.Errorf("could not find a safe unique key for upsert")
}

func UpsertDocuments(
	database *mongo.Database,
	collectionName string,
	documents []parser.Document,
	key string,
) (int, int, int, error) {
	collection := database.Collection(collectionName)

	// Make sure every document has the key before starting the import.
	for _, document := range documents {
		if _, exists := document[key]; !exists {
			return 0, 0, 0, fmt.Errorf("document is missing upsert key: %s", key)
		}
	}

	insertedCount := 0
	updatedCount := 0
	matchedCount := 0

	for _, document := range documents {
		value := document[key]
		filter := map[string]interface{}{key: value}
		update := map[string]interface{}{"$set": document}

		result, err := collection.UpdateOne(
			context.Background(),
			filter,
			update,
			options.UpdateOne().SetUpsert(true),
		)
		if err != nil {
			return insertedCount, updatedCount, matchedCount, err
		}

		if result.UpsertedCount > 0 {
			insertedCount++
		} else {
			matchedCount++
			if result.ModifiedCount > 0 {
				updatedCount++
			}
		}
	}

	return insertedCount, updatedCount, matchedCount, nil
}
