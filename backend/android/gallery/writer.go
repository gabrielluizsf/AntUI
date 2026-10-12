//go:build android

package gallery

import (
	"os"
	"sync"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Writer is something being written into the device's library.
//
// Nothing is visible to the gallery until [Writer.Close]: on Android 10 and
// up the entry is marked pending while it is being written, so a half-copied
// video never appears in anyone's photo roll. Closing is what publishes it,
// and a Writer that is dropped without closing leaves an entry nothing can
// see and nothing will clean up — so close it, on the failure path too.
type Writer struct {
	file *os.File
	uri  jni.Object
	// uriText is the content:// address, which is what other apps are given.
	uriText string

	once   sync.Once
	closed bool
	err    error
}
