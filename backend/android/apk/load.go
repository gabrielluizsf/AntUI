package apk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Load reads a config file. Paths inside it are resolved against the
// directory the file is in.
//
// The signing passwords are not read from the file — they are not written to
// one either — and come from ANTUIAPK_STOREPASS and ANTUIAPK_KEYPASS if
// those are set. A caller that has them from somewhere better should set
// them on the result afterwards.
func Load(path string) (*Config, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("apk: %w", err)
	}
	var cfg Config
	dec := json.NewDecoder(bytes.NewReader(body))
	// A misspelt key is the whole reason a config file goes wrong quietly:
	// "permission" instead of "permissions" is not an error anywhere, it is
	// simply an app without the permission.
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("apk: reading %s: %w", path, err)
	}

	base := filepath.Dir(path)
	for _, p := range []*string{&cfg.Dir, &cfg.Icon, &cfg.Assets, &cfg.Out,
		&cfg.Symbols, &cfg.Keystore.Path} {
		if *p != "" && !filepath.IsAbs(*p) {
			*p = filepath.Join(base, *p)
		}
	}
	if v := os.Getenv("ANTUIAPK_STOREPASS"); v != "" {
		cfg.Keystore.StorePass = v
	}
	if v := os.Getenv("ANTUIAPK_KEYPASS"); v != "" {
		cfg.Keystore.KeyPass = v
	}
	return &cfg, nil
}
