package cli

import "testing"

func TestParseImportMode(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    ImportMode
		wantErr bool
	}{
		{
			name: "default",
			args: nil,
			want: ImportModeAsk,
		},
		{
			name: "append",
			args: []string{"--mode", "append"},
			want: ImportModeAppend,
		},
		{
			name: "replace",
			args: []string{"--mode", "replace"},
			want: ImportModeReplace,
		},
		{
			name: "upsert",
			args: []string{"--mode", "upsert"},
			want: ImportModeUpsert,
		},
		{
			name: "ask",
			args: []string{"--mode", "ask"},
			want: ImportModeAsk,
		},
		{
			name:    "invalid mode",
			args:    []string{"--mode", "banana"},
			wantErr: true,
		},
		{
			name:    "missing value",
			args:    []string{"--mode"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseImportMode(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseCollectionMode(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    CollectionMode
		wantErr bool
	}{
		{
			name: "default",
			args: nil,
			want: CollectionModeAsk,
		},
		{
			name: "ask",
			args: []string{"--collection-mode", "ask"},
			want: CollectionModeAsk,
		},
		{
			name: "combine",
			args: []string{"--collection-mode", "combine"},
			want: CollectionModeCombine,
		},
		{
			name: "separate",
			args: []string{"--collection-mode", "separate"},
			want: CollectionModeSeparate,
		},
		{
			name: "cancel",
			args: []string{"--collection-mode", "cancel"},
			want: CollectionModeCancel,
		},
		{
			name:    "invalid mode",
			args:    []string{"--collection-mode", "banana"},
			wantErr: true,
		},
		{
			name:    "missing value",
			args:    []string{"--collection-mode"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCollectionMode(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseUpsertKey(t *testing.T) {
	got := ParseUpsertKey([]string{
		"--mode", "upsert",
		"--key", "customer_id",
	})

	if got != "customer_id" {
		t.Fatalf("got %q, want customer_id", got)
	}
}

func TestParseDatabase(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "default",
			args: nil,
			want: "mongo-easy",
		},
		{
			name: "database option",
			args: []string{"--database", "analytics"},
			want: "analytics",
		},
		{
			name: "short database option",
			args: []string{"-db", "analytics"},
			want: "analytics",
		},
		{
			name:    "missing value",
			args:    []string{"--database"},
			wantErr: true,
		},
		{
			name:    "missing value short option",
			args:    []string{"-db"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDatabase(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
