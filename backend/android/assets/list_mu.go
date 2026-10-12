//go:build android

package assets

import (
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
