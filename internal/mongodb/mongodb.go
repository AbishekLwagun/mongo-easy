package mongodb

import (
	"context"

	"github.com/abisheklwagun/mongo-easy/internal/parser"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Connect creates a MongoDB client using the local MongoDB server.
func Connect() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetServerAPIOptions(serverAPI)

	return mongo.Connect(opts)
}

// HasDocuments checks if the collection already has at least one document.
// This is used by the default "ask" mode before importing new data.
func HasDocuments(database *mongo.Database, collectionName string) (bool, error) {
	var document parser.Document

	err := database.Collection(collectionName).
		FindOne(context.Background(), map[string]interface{}{}).
		Decode(&document)

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// InsertDocuments adds all parsed documents to the selected collection.
func InsertDocuments(
	database *mongo.Database,
	collectionName string,
	documents []parser.Document,
) error {
	values := make([]interface{}, len(documents))

	for i, document := range documents {
		values[i] = document
	}

	_, err := database.Collection(collectionName).
		InsertMany(context.Background(), values)

	return err
}

// ReplaceCollection removes the existing collection and then imports
// the new documents into a fresh collection with the same name.
func ReplaceCollection(
	database *mongo.Database,
	collectionName string,
	documents []parser.Document,
) error {
	if err := database.Collection(collectionName).Drop(context.Background()); err != nil {
		return err
	}

	return InsertDocuments(database, collectionName, documents)
}
