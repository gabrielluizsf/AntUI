package apk

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// overlayFor writes the injected file and the overlay that points at it,
// returning the overlay's path.
//
// An overlay maps a path the go command would read to a different file on
// disk. Mapping a path that does not exist is how a file is *added* to a
// package without touching the caller's directory — which matters, because
// building someone's program should not write into their source tree.
func overlayFor(cfg *Config, work string) (string, error) {
	src := filepath.Join(work, injectedName)
	if err := os.WriteFile(src, []byte(injected), 0o644); err != nil {
		return "", fmt.Errorf("apk: %w", err)
	}
	doc := struct{ Replace map[string]string }{
		Replace: map[string]string{filepath.Join(cfg.Dir, injectedName): src},
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("apk: %w", err)
	}
	path := filepath.Join(work, "overlay.json")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return "", fmt.Errorf("apk: %w", err)
	}
	return path, nil
}
