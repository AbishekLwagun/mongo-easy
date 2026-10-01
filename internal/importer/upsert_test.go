package importer

import (
	"testing"

	"github.com/abisheklwagun/mongo-easy/internal/parser"
)

func TestDetectUpsertKey(t *testing.T) {
	documents := []parser.Document{
		{"customer_id": 1, "name": "A"},
		{"customer_id": 2, "name": "B"},
	}

	got, err := DetectUpsertKey(documents)
	if err != nil {
		t.Fatal(err)
	}
	if got != "customer_id" {
		t.Fatalf("got %q, want customer_id", got)
	}
}

func TestDetectUpsertKeyRejectsDuplicate(t *testing.T) {
	documents := []parser.Document{
		{"email": "same@example.com"},
		{"email": "same@example.com"},
	}

	if _, err := DetectUpsertKey(documents); err == nil {
		t.Fatal("expected duplicate key to be rejected")
	}
}

func TestDetectUpsertKeyRejectsNoCandidate(t *testing.T) {
	documents := []parser.Document{
		{"product": "Laptop"},
	}

	if _, err := DetectUpsertKey(documents); err == nil {
		t.Fatal("expected no safe key error")
	}
}
