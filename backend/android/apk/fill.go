package apk

import (
	"os"
	"path/filepath"
)

func (k *Keystore) fill() {
	if k.Path == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			k.Path = filepath.Join(home, ".android", "debug.keystore")
		}
		if k.Alias == "" {
			k.Alias = debugAlias
		}
		if k.StorePass == "" {
			k.StorePass = debugPass
		}
	}
	if k.KeyPass == "" {
		k.KeyPass = k.StorePass
	}
}
