//go:build android

package share

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Create makes a file other apps can be given, and returns it along with the
// address to give them.
//
// The name may not have a directory in it: everything shared is one flat
// folder, because a provider that serves a tree is a provider that has to be
// argued with about what is inside it.
func Create(name string) (*os.File, string, error) {
	if name == "" || strings.ContainsAny(name, `/\`) || name == "." || name == ".." {
		return nil, "", errors.New("share: a shared file's name may not have a path in it")
	}
	dir, err := Dir()
	if err != nil {
		return nil, "", err
	}
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return nil, "", err
	}
	uri, err := URI(name)
	if err != nil {
		f.Close()
		return nil, "", err
	}
	return f, uri, nil
}
