//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/asset_manager.h>
#include <stdlib.h>
*/
import "C"

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
