//go:build android

package share

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// The one directory the provider serves, under the app's cache. It matches
// what the Java side opens; the two have to agree and there is nowhere else
// for them to agree in.
const dirName = "antui-share"

var (
	pkgOnce sync.Once
	pkgName string
)

func packageName() (string, error) {
	var err error
	pkgOnce.Do(func() {
		err = jni.Do(func(e *jni.Env) error {
			var e2 error
			pkgName, e2 = e.InvokeString(app.Context(), "getPackageName", jni.Sig(jni.TString))
			return e2
		})
	})
	if err != nil {
		return "", err
	}
	if pkgName == "" {
		return "", errors.New("share: the app does not know its own name")
	}
	return pkgName, nil
}

// Dir is where files to be handed to other apps live.
//
// It is inside the cache, which means the system may empty it when the
// device runs short — and that is right: a file put here is a copy made to
// be given away, and if it goes, it goes.
func Dir() (string, error) {
	a := app.Current()
	if a == nil {
		return "", app.ErrNoUIThread
	}
	cache := a.Info().Cache
	if cache == "" {
		return "", errors.New("share: the app does not know its cache directory")
	}
	dir := filepath.Join(cache, dirName)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

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

// Write makes a file with these bytes in it and gives back its address.
func Write(name string, data []byte) (string, error) {
	f, uri, err := Create(name)
	if err != nil {
		return "", err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return uri, nil
}

// URI is the address of a file already in [Dir].
//
// It is a content:// address and not a path. **Since Android 7 a path is not
// something one app may give another**: an intent carrying a file:// address
// throws, deliberately, because a path is meaningless in another app's world
// and used to be a way to hand out things that were never meant to leave.
func URI(name string) (string, error) {
	pkg, err := packageName()
	if err != nil {
		return "", err
	}
	return "content://" + pkg + shim.ProviderSuffix + "/" + url.PathEscape(name), nil
}

// Clear empties the directory. Worth doing when the app is done handing
// things out: the copies are of no further use and the system will otherwise
// keep them until it needs the space.
func Clear() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		os.RemoveAll(filepath.Join(dir, e.Name()))
	}
	return nil
}
