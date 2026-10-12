//go:build android

package ndk

/*
#cgo LDFLAGS: -lmediandk -landroid

#include "media_cgo.h"
*/
import "C"

import (
	"os"
)

// Decoder is one track being decoded.
type Decoder struct {
	media C.mw_media
	// file is the one the extractor reads from, kept open for its life.
	file   *os.File
	track  Track
	rotate int
	closed bool
	// ended is set once the codec has produced its last output.
	ended bool
}
