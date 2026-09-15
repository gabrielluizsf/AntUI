//go:build android

package assets

import (
	"errors"
	"io"
	"io/fs"
	"path"
	"sort"
	"sync"
	"time"

	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// FS is everything packed into the APK, as an ordinary [io/fs.FS].
//
// That is the whole point of this package. A game that reads its levels with
// os.ReadFile on a desktop reads them with fs.ReadFile here, and the same
// code serves both — which is the difference between porting a program and
// rewriting it.
//
//	levels, err := fs.ReadFile(assets.FS(), "levels/one.json")
//	err = fs.WalkDir(assets.FS(), "sounds", func(p string, d fs.DirEntry, err error) error { ... })
func FS() fs.FS { return assetFS{} }

// Open is fs.FS.Open on [FS].
func Open(name string) (fs.File, error) { return assetFS{}.Open(name) }

// ReadFile is the whole of one asset.
func ReadFile(name string) ([]byte, error) { return assetFS{}.ReadFile(name) }

// ReadDir is what is directly inside a directory.
func ReadDir(name string) ([]fs.DirEntry, error) { return assetFS{}.ReadDir(name) }

// assetFS satisfies fs.FS, fs.ReadDirFS and fs.ReadFileFS.
type assetFS struct{}

var (
	_ fs.FS         = assetFS{}
	_ fs.ReadDirFS  = assetFS{}
	_ fs.ReadFileFS = assetFS{}
)

func manager() (*ndk.AssetManager, error) {
	a := app.Current()
	if a == nil {
		return nil, app.ErrNoUIThread
	}
	m := a.Assets()
	if m == nil {
		return nil, errors.New("assets: this app has no asset manager")
	}
	return m, nil
}

// inner is the name the platform wants: no leading slash, and "" for the
// root rather than fs's ".".
func inner(name string) string {
	if name == "." {
		return ""
	}
	return name
}

func (f assetFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	m, err := manager()
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}
	// A directory first: opening one as a file gives an error from the
	// platform that says nothing about it being a directory.
	if entries, err := list(name); err == nil && len(entries) > 0 {
		return &dirFile{name: path.Base(name), entries: entries}, nil
	}
	// Buffered, so that an uncompressed asset can be handed over as the
	// mapping it already is rather than copied.
	a, err := m.Open(inner(name), ndk.AssetBuffer)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return &file{name: path.Base(name), asset: a}, nil
}

func (f assetFS) ReadFile(name string) ([]byte, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrInvalid}
	}
	m, err := manager()
	if err != nil {
		return nil, &fs.PathError{Op: "read", Path: name, Err: err}
	}
	a, err := m.Open(inner(name), ndk.AssetBuffer)
	if err != nil {
		return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrNotExist}
	}
	defer a.Close()
	return a.ReadAll()
}

func (f assetFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrInvalid}
	}
	entries, err := list(name)
	if err != nil {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: err}
	}
	if len(entries) == 0 {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}
	return entries, nil
}

// file is one asset open for reading.
type file struct {
	name  string
	asset *ndk.Asset
}

func (f *file) Stat() (fs.FileInfo, error) {
	return info{name: f.name, size: f.asset.Size()}, nil
}
func (f *file) Read(p []byte) (int, error) { return f.asset.Read(p) }
func (f *file) Close() error               { return f.asset.Close() }

// Seek is here because a great many readers want it — an image decoder, a zip
// reader — and an uncompressed asset can do it.
func (f *file) Seek(offset int64, whence int) (int64, error) {
	return f.asset.Seek(offset, whence)
}

// dirFile is a directory opened as a file, which fs.WalkDir does.
type dirFile struct {
	name    string
	entries []fs.DirEntry
	at      int
}

func (d *dirFile) Stat() (fs.FileInfo, error) { return info{name: d.name, dir: true}, nil }
func (d *dirFile) Close() error               { return nil }
func (d *dirFile) Read([]byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.name, Err: errors.New("is a directory")}
}

func (d *dirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if n <= 0 {
		out := d.entries[d.at:]
		d.at = len(d.entries)
		return out, nil
	}
	if d.at >= len(d.entries) {
		return nil, io.EOF
	}
	end := min(d.at+n, len(d.entries))
	out := d.entries[d.at:end]
	d.at = end
	return out, nil
}

// info is what Stat gives back. There are no times and no owners inside a
// package, so there is nothing to report but the name, the size and whether
// it is a directory.
type info struct {
	name string
	size int64
	dir  bool
}

func (i info) Name() string { return i.name }
func (i info) Size() int64  { return i.size }
func (i info) Mode() fs.FileMode {
	if i.dir {
		return fs.ModeDir | 0o555
	}
	return 0o444
}
func (i info) ModTime() time.Time { return time.Time{} }
func (i info) IsDir() bool        { return i.dir }
func (i info) Sys() any           { return nil }

// entry is one thing in a directory.
type entry struct{ info info }

func (e entry) Name() string               { return e.info.name }
func (e entry) IsDir() bool                { return e.info.dir }
func (e entry) Type() fs.FileMode          { return e.info.Mode().Type() }
func (e entry) Info() (fs.FileInfo, error) { return e.info, nil }

// The listing goes through Java.
//
// The NDK's AAssetDir lists **files only**: there is no call in it that
// reports a subdirectory, so a walk over the tree with it would stop at the
// first level and report nothing missing. AssetManager.list in Java does
// include directories, so that is what is used — and it is the one thing in
// this package that crosses the bridge.
var (
	listMu     sync.Mutex
	javaAssets jni.Object
)

func list(dir string) ([]fs.DirEntry, error) {
	names, err := listNames(dir)
	if err != nil || len(names) == 0 {
		return nil, err
	}
	out := make([]fs.DirEntry, 0, len(names))
	for _, name := range names {
		// A name whose own listing has anything in it is a directory. There
		// is no other way to ask: the platform reports a flat list of names
		// and nothing about what they are. An empty directory cannot be told
		// from a file this way, and does not need to be — aapt2 drops empty
		// directories on the way into the package.
		child := name
		if d := inner(dir); d != "" {
			child = d + "/" + name
		}
		sub, _ := listNames(child)
		out = append(out, entry{info: info{name: name, dir: len(sub) > 0}})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out, nil
}

func listNames(dir string) ([]string, error) {
	var out []string
	err := jni.Do(func(e *jni.Env) error {
		listMu.Lock()
		held := javaAssets
		listMu.Unlock()
		if held.IsNil() {
			got, err := e.Invoke(app.Context(), "getAssets",
				jni.Sig(jni.TClass("android/content/res/AssetManager")))
			if err != nil {
				return err
			}
			held = e.Global(got)
			listMu.Lock()
			if javaAssets.IsNil() {
				javaAssets = held
			} else {
				e.DeleteGlobal(held)
				held = javaAssets
			}
			listMu.Unlock()
		}
		jdir, err := e.String(inner(dir))
		if err != nil {
			return err
		}
		arr, err := e.Invoke(held, "list",
			jni.Sig(jni.TArray(jni.TString), jni.TString), jni.Ref(jdir))
		if err != nil || arr.IsNil() {
			return err
		}
		out, err = e.GoStrings(arr)
		return err
	})
	return out, err
}
