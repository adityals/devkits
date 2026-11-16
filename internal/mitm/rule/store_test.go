package rule

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileStoreRename(t *testing.T) {
	dir := t.TempDir()
	store := NewFileStore(dir)

	for _, r := range []Rule{
		{ID: "old", Enabled: false, Script: "a"},
		{ID: "other", Enabled: true, Script: "b"},
	} {
		if err := store.Save(r); err != nil {
			t.Fatalf("Save %s: %v", r.ID, err)
		}
	}

	if err := store.Rename("old", "new", false); err != nil {
		t.Fatalf("Rename: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "_new.js")); err != nil {
		t.Fatalf("disabled rule lost its prefix: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "_old.js")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("old file still there: %v", err)
	}

	// The target is enabled while the source is disabled, so only the cross-state check catches it.
	if err := store.Rename("new", "other", false); !errors.Is(err, ErrExists) {
		t.Fatalf("Rename onto existing rule err = %v, want ErrExists", err)
	}
}
