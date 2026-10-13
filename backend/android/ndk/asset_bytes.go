//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/asset_manager.h>
#include <stdlib.h>
*/
import "C"

import "unsafe"

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
