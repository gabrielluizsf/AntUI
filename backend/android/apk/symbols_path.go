package apk

import (
	"path/filepath"
	"strings"
)

// symbolsPath is where the archive goes when the config does not say: beside
// the package, named after it.
//
//	hello.apk  →  hello-symbols.zip
//	hello.aab  →  hello-symbols.zip
func symbolsPath(cfg *Config) string {
	if cfg.Symbols != "" {
		return cfg.Symbols
	}
	out := cfg.Out
	return strings.TrimSuffix(out, filepath.Ext(out)) + "-symbols.zip"
}
