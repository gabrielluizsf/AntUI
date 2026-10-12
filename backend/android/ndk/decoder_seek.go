//go:build android

package ndk

/*
#cgo LDFLAGS: -lmediandk -landroid

#include "media_cgo.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"time"
)

// Seek moves to a position. The decoder jumps to the nearest key frame at or
// before it, which for video may be a second earlier — there is no other
// kind of seek in a compressed stream.
func (d *Decoder) Seek(to time.Duration) error {
	if d.closed {
		return errors.New("ndk: the decoder is closed")
	}
	if r := C.AMediaExtractor_seekTo(d.media.ex, C.int64_t(to.Microseconds()),
		C.AMEDIAEXTRACTOR_SEEK_PREVIOUS_SYNC); r != C.AMEDIA_OK {
		return fmt.Errorf("ndk: cannot seek to %v (%d)", to, int(r))
	}
	C.AMediaCodec_flush(d.media.codec)
	d.media.eof = 0
	d.media.drained = 0
	d.ended = false
	return nil
}
