package apk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// A build described in a file rather than in twenty flags.
//
// Everything in [Config] can be written down, and the file is JSON because
// that is what the standard library reads and because an editor can check
// it. The one thing missing from it is the signing passwords, which are the
// one thing that must never be in a file next to the source.
//
//	{
//	  "package": "com.example.hello",
//	  "label": "Hello",
//	  "versionCode": 3,
//	  "versionName": "1.2",
//	  "icon": "icon.png",
//	  "iconBackground": "#3E63DD",
//	  "permissions": ["android.permission.INTERNET"],
//	  "keystore": { "path": "release.jks", "alias": "release" }
//	}
//
// Every path in it — the icon, the assets, the keystore, the output — is
// relative to the file itself rather than to wherever the command was run.
// A config that means different things from different directories is not a
// config.

// ConfigName is the file [Load] looks for when it is not told.
const ConfigName = "antuiapk.json"

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
