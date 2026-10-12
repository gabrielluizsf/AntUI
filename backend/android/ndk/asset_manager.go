//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/asset_manager.h>
#include <stdlib.h>
*/
import "C"

// AssetManager reads the files packed into the APK.
//
// Assets are not files on disk: they are entries inside the package, and the
// package is a zip that is never unpacked. Reading one is reading out of the
// APK, which is why there is an API for it at all and why no path leads to
// one.
type AssetManager struct {
	ptr *C.AAssetManager
}
