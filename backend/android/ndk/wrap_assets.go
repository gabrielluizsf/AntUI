//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/asset_manager.h>
#include <stdlib.h>
*/
import "C"

import "unsafe"

// WrapAssets takes the manager the platform handed the activity.
func WrapAssets(p unsafe.Pointer) *AssetManager {
	if p == nil {
		return nil
	}
	return &AssetManager{ptr: (*C.AAssetManager)(p)}
}
