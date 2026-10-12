package parser

import (
	"fmt"
	"strings"
	"testing"
)

func genJSONL(n int) []byte {
	var b strings.Builder
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "{\"customer_id\":%d,\"name\":\"User %d\",\"email\":\"user%d@example.com\",\"age\":%d}\n", i, i, i, 20+i%50)
	}
	return []byte(b.String())
}

func genCSV(n int) []byte {
	var b strings.Builder
	b.WriteString("customer_id,name,email,age\n")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%d,User %d,user%d@example.com,%d\n", i, i, i, 20+i%50)
	}
	return []byte(b.String())
}

func BenchmarkParseJSONL100k(b *testing.B) {
	data := genJSONL(100000)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		docs, err := ParseJSONL(data)
		if err != nil || len(docs) != 100000 {
			b.Fatal(err, len(docs))
		}
	}
}

func BenchmarkParseCSV100k(b *testing.B) {
	data := genCSV(100000)
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		docs, err := ParseCSV(data)
		if err != nil || len(docs) != 100000 {
			b.Fatal(err, len(docs))
		}
	}
}
