//go:build unit

package databaseshared

import "testing"

func TestCreateDatabaseEngineSlug(t *testing.T) {
	tests := []struct {
		name    string
		engine  string
		version string
		want    string
	}{
		{name: "mysql major version", engine: "mysql", version: "8.0.36", want: "mysql/8"},
		{name: "mysql minor version", engine: "mysql", version: "8.4.0", want: "mysql/8.4"},
		{name: "postgres major version", engine: "postgresql", version: "16.4", want: "postgresql/16"},
		{name: "valkey minor version", engine: "valkey", version: "8.1.0", want: "valkey/8.1"},
		{name: "valkey major version", engine: "valkey", version: "8.0.0", want: "valkey/8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateDatabaseEngineSlug(tt.engine, tt.version); got != tt.want {
				t.Errorf("CreateDatabaseEngineSlug(%q, %q) = %q, want %q", tt.engine, tt.version, got, tt.want)
			}
		})
	}
}

func TestDatabaseEngineSlugRoundTrip(t *testing.T) {
	tests := []struct {
		engine  string
		version string
		want    string
	}{
		{engine: "mysql", version: "8.0.36", want: "mysql/8"},
		{engine: "postgresql", version: "16.4", want: "postgresql/16"},
		{engine: "valkey", version: "8.1.0", want: "valkey/8.1"},
	}

	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			engineID := CreateDatabaseEngineSlug(test.engine, test.version)
			if engineID != test.want {
				t.Fatalf("CreateDatabaseEngineSlug() = %q, want %q", engineID, test.want)
			}

			engine, version, err := ParseDatabaseEngineSlug(engineID)
			if err != nil {
				t.Fatalf("ParseDatabaseEngineSlug() error = %v", err)
			}
			if engine != test.engine || version != test.want[len(test.engine)+1:] {
				t.Fatalf("ParseDatabaseEngineSlug() = (%q, %q), want (%q, %q)", engine, version, test.engine, test.want[len(test.engine)+1:])
			}
		})
	}
}
