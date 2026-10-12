package apk

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Save writes a config file. It writes what is set and leaves out what is
// not, so that the file says what the app decided rather than restating
// every default — a default that is written down is a default that stops
// following the library when it changes.
func (c Config) Save(path string) error {
	// Relative to the file, the same way Load reads them.
	base, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	out := c
	for _, p := range []*string{&out.Dir, &out.Icon, &out.Assets, &out.Out,
		&out.Symbols, &out.Keystore.Path} {
		if *p == "" {
			continue
		}
		if abs, err := filepath.Abs(*p); err == nil {
			// Relative when relative is shorter, which is the same as
			// saying "when the file and what it points at live together".
			// A config beside its app says "." and moves with the
			// repository; a config pointing across the disk says so plainly
			// rather than climbing six levels to get there.
			if rel, err := filepath.Rel(base, abs); err == nil && len(rel) < len(abs) {
				*p = filepath.ToSlash(rel)
			}
		}
	}

	body, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	return os.WriteFile(path, append(body, '\n'), 0o644)
}
