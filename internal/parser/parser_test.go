package parser

import "testing"

func TestParseJSON(t *testing.T) {
	documents, err := ParseJSON([]byte(`[{"id":1},{"id":2}]`))
	if err != nil {
		t.Fatal(err)
	}
	if len(documents) != 2 {
		t.Fatalf("got %d documents, want 2", len(documents))
	}
}

func TestParseJSONSingleDocument(t *testing.T) {
	documents, err := ParseJSON([]byte(`{"product":"Laptop"}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(documents) != 1 {
		t.Fatalf("got %d documents, want 1", len(documents))
	}
}

func TestParseJSONL(t *testing.T) {
	documents, err := ParseJSONL([]byte("{\"id\":1}\n{\"id\":2}\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(documents) != 2 {
		t.Fatalf("got %d documents, want 2", len(documents))
	}
}

func TestParseCSV(t *testing.T) {
	documents, err := ParseCSV([]byte("id,name\n1,A\n2,B\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(documents) != 2 {
		t.Fatalf("got %d documents, want 2", len(documents))
	}
}
