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
	"unsafe"
)

func (d *Decoder) configure(window *C.ANativeWindow) error {
	mime := C.CString(d.track.Mime)
	defer C.free(unsafe.Pointer(mime))
	d.media.codec = C.AMediaCodec_createDecoderByType(mime)
	if d.media.codec == nil {
		return fmt.Errorf("ndk: this device has no decoder for %s", d.track.Mime)
	}
	format := C.AMediaExtractor_getTrackFormat(d.media.ex, C.size_t(d.track.Index))
	if format == nil {
		return errors.New("ndk: the track has no format")
	}
	defer C.AMediaFormat_delete(format)

	if r := C.AMediaCodec_configure(d.media.codec, format, window, nil, 0); r != C.AMEDIA_OK {
		return fmt.Errorf("ndk: cannot set up the %s decoder (%d)", d.track.Mime, int(r))
	}
	if r := C.AMediaCodec_start(d.media.codec); r != C.AMEDIA_OK {
		return fmt.Errorf("ndk: cannot start the %s decoder (%d)", d.track.Mime, int(r))
	}
	return nil
}
