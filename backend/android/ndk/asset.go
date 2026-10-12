//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/asset_manager.h>
#include <stdlib.h>
*/
import "C"

// Asset is one file inside the package, open for reading.
type Asset struct {
	ptr  *C.AAsset
	size int64
}
