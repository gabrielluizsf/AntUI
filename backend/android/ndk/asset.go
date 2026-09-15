//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/asset_manager.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"io"
	"unsafe"
)

// AssetManager reads the files packed into the APK.
//
// Assets are not files on disk: they are entries inside the package, and the
// package is a zip that is never unpacked. Reading one is reading out of the
// APK, which is why there is an API for it at all and why no path leads to
// one.
type AssetManager struct {
	ptr *C.AAssetManager
}

// WrapAssets takes the manager the platform handed the activity.
func WrapAssets(p unsafe.Pointer) *AssetManager {
	if p == nil {
		return nil
	}
	return &AssetManager{ptr: (*C.AAssetManager)(p)}
}

// How an asset is meant to be read. It is a hint: the platform uses it to
// decide whether to map the file or to buffer it.
const (
	// AssetRandom is for something read out of order.
	AssetRandom = int(C.AASSET_MODE_RANDOM)
	// AssetStreaming is for something read start to end once.
	AssetStreaming = int(C.AASSET_MODE_STREAMING)
	// AssetBuffer asks for the whole thing at once, which is what makes
	// [Asset.Bytes] able to hand back the file with no copy.
	AssetBuffer = int(C.AASSET_MODE_BUFFER)
)

// ErrNoAsset is a name that is not in the package.
var ErrNoAsset = errors.New("ndk: no such asset")

// Asset is one file inside the package, open for reading.
type Asset struct {
	ptr  *C.AAsset
	size int64
}

// Open finds an asset by its path inside the package, with no leading slash:
// "levels/one.json".
func (m *AssetManager) Open(name string, mode int) (*Asset, error) {
	if m == nil || m.ptr == nil {
		return nil, errors.New("ndk: there is no asset manager")
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	a := C.AAssetManager_open(m.ptr, cname, C.int(mode))
	if a == nil {
		return nil, ErrNoAsset
	}
	return &Asset{ptr: a, size: int64(C.AAsset_getLength64(a))}, nil
}

// Size is how many bytes the asset is.
func (a *Asset) Size() int64 { return a.size }

// Read fills p, as an io.Reader does.
func (a *Asset) Read(p []byte) (int, error) {
	if a.ptr == nil {
		return 0, errors.New("ndk: the asset is closed")
	}
	if len(p) == 0 {
		return 0, nil
	}
	n := int(C.AAsset_read(a.ptr, unsafe.Pointer(&p[0]), C.size_t(len(p))))
	switch {
	case n > 0:
		return n, nil
	case n == 0:
		return 0, io.EOF
	}
	return 0, errors.New("ndk: reading the asset failed")
}

// Seek moves the read position, as an io.Seeker does.
func (a *Asset) Seek(offset int64, whence int) (int64, error) {
	if a.ptr == nil {
		return 0, errors.New("ndk: the asset is closed")
	}
	at := int64(C.AAsset_seek64(a.ptr, C.off64_t(offset), C.int(whence)))
	if at < 0 {
		return 0, errors.New("ndk: cannot seek in a compressed asset")
	}
	return at, nil
}

// Bytes is the whole asset, without a copy where the platform can manage it.
//
// It works when the asset is stored uncompressed in the package, because
// then it can simply be mapped; a compressed one has to be inflated and
// there is nothing to point at. **The slice is the mapping**, not Go memory:
// it stops being valid when the asset is closed, and writing to it is
// writing to a read-only page.
//
// It reports false when the asset was compressed, and [Asset.ReadAll] is
// what to use then — or, better, the asset should be stored rather than
// compressed. See the -0 flag antuiapk passes to aapt2.
func (a *Asset) Bytes() ([]byte, bool) {
	if a.ptr == nil {
		return nil, false
	}
	p := C.AAsset_getBuffer(a.ptr)
	if p == nil || a.size <= 0 {
		return nil, false
	}
	return unsafe.Slice((*byte)(p), a.size), true
}

// ReadAll is the whole asset as Go memory, mapped where it can be and
// inflated where it cannot.
func (a *Asset) ReadAll() ([]byte, error) {
	if mapped, ok := a.Bytes(); ok {
		out := make([]byte, len(mapped))
		copy(out, mapped)
		return out, nil
	}
	if a.size <= 0 {
		return nil, nil
	}
	out := make([]byte, a.size)
	read := 0
	for read < len(out) {
		n, err := a.Read(out[read:])
		if n > 0 {
			read += n
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return out[:read], nil
}

// Close gives the asset back. A mapping handed out by [Asset.Bytes] stops
// being valid here.
func (a *Asset) Close() error {
	if a.ptr != nil {
		C.AAsset_close(a.ptr)
		a.ptr = nil
	}
	return nil
}

// Files lists the assets directly inside a directory.
//
// **It lists files and not directories**, which is a limitation of this API
// and not of the package: there is no NDK call that reports a subdirectory.
// The assets package above this one fills that in through Java, where the
// same list does include them.
func (m *AssetManager) Files(dir string) ([]string, error) {
	if m == nil || m.ptr == nil {
		return nil, errors.New("ndk: there is no asset manager")
	}
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	d := C.AAssetManager_openDir(m.ptr, cdir)
	if d == nil {
		return nil, ErrNoAsset
	}
	defer C.AAssetDir_close(d)

	var out []string
	for {
		name := C.AAssetDir_getNextFileName(d)
		if name == nil {
			break
		}
		out = append(out, C.GoString(name))
	}
	return out, nil
}
