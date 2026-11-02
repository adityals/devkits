package rule

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const disabledPrefix = "_"

// idPattern keeps ids safe as file names and blocks path traversal.
var idPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

type Store interface {
	Load() ([]Rule, error)
	Save(Rule) error
	Delete(id string) error
}

// FileStore keeps one rule per file: <id>.js is enabled, _<id>.js is disabled.
type FileStore struct {
	dir string
}

func NewFileStore(dir string) *FileStore {
	return &FileStore{dir: dir}
}

func ValidateID(id string) error {
	if !idPattern.MatchString(id) {
		return fmt.Errorf("%w: id %q must match %s", ErrInvalidRule, id, idPattern)
	}

	return nil
}

func (s *FileStore) Load() ([]Rule, error) {
	entries, err := os.ReadDir(s.dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("list rule dir: %w", err)
	}

	var rules []Rule
	for _, e := range entries {
		name, ok := strings.CutSuffix(e.Name(), ".js")
		if e.IsDir() || !ok {
			continue
		}

		id, disabled := strings.CutPrefix(name, disabledPrefix)
		if err := ValidateID(id); err != nil {
			return nil, fmt.Errorf("rule file %s: %w", e.Name(), err)
		}

		if slices.ContainsFunc(rules, func(r Rule) bool { return r.ID == id }) {
			return nil, fmt.Errorf("%w: both %s.js and _%s.js exist", ErrInvalidRule, id, id)
		}

		script, err := os.ReadFile(filepath.Join(s.dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read rule %s: %w", e.Name(), err)
		}

		rules = append(rules, Rule{
			ID:      id,
			Enabled: !disabled,
			Script:  string(script),
		})
	}

	if len(rules) > MaxRules {
		return nil, fmt.Errorf("%w: %d rule files, max %d", ErrMaxRules, len(rules), MaxRules)
	}

	slices.SortFunc(rules, func(a, b Rule) int {
		return strings.Compare(a.ID, b.ID)
	})

	return rules, nil
}

func (s *FileStore) Save(r Rule) error {
	if err := ValidateID(r.ID); err != nil {
		return err
	}

	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return fmt.Errorf("create rule dir: %w", err)
	}

	if err := s.writeAtomic(s.path(r.ID, r.Enabled), []byte(r.Script)); err != nil {
		return fmt.Errorf("write rule %s: %w", r.ID, err)
	}

	return s.remove(s.path(r.ID, !r.Enabled))
}

func (s *FileStore) Delete(id string) error {
	if err := ValidateID(id); err != nil {
		return err
	}

	return errors.Join(
		s.remove(s.path(id, true)),
		s.remove(s.path(id, false)),
	)
}

func (s *FileStore) path(id string, enabled bool) string {
	if !enabled {
		id = disabledPrefix + id
	}

	return filepath.Join(s.dir, id+".js")
}

func (s *FileStore) remove(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove rule file: %w", err)
	}

	return nil
}

// writeAtomic writes via temp file + rename so a crash never leaves a half-written file.
func (s *FileStore) writeAtomic(path string, data []byte) error {
	tmp, err := os.CreateTemp(s.dir, ".tmp-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), path)
}
