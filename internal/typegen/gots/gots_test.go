package gots

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite testdata/*.ts golden files")

func TestConverter_Convert(t *testing.T) {
	inputs, err := filepath.Glob("testdata/*.go")
	if err != nil {
		t.Fatal(err)
	}

	if len(inputs) == 0 {
		t.Fatal("no testdata/*.go fixtures found")
	}

	for _, in := range inputs {
		name := strings.TrimSuffix(filepath.Base(in), ".go")
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			src, err := os.ReadFile(in)
			if err != nil {
				t.Fatal(err)
			}

			got, err := New().Convert(string(src))
			if err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			golden := strings.TrimSuffix(in, ".go") + ".ts"
			if *update {
				if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			expected, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("missing golden, run: go test ./internal/typegen/gots -update: %v", err)
			}

			if got != string(expected) {
				t.Errorf("output mismatch\n--- got\n%s\n--- want\n%s", got, expected)
			}
		})
	}
}

func TestConverter_ConvertInvalid(t *testing.T) {
	if _, err := New().Convert("type Broken struct {"); err == nil {
		t.Fatal("Convert() error = nil, want error for invalid go")
	}
}
