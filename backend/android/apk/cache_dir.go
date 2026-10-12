package apk

import (
	"os"
	"path/filepath"
)

func cacheDir() (string, error) {
	root, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(root, "antuiapk", "res")
	return dir, os.MkdirAll(dir, 0o755)
}
