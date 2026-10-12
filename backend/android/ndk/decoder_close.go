//go:build android

package ndk

/*
#cgo LDFLAGS: -lmediandk -landroid

#include "media_cgo.h"
*/
import "C"

// Close stops the decoder and gives everything back.
func (d *Decoder) Close() error {
	if d.closed {
		return nil
	}
	d.closed = true
	C.mw_media_close(&d.media)
	if d.file != nil {
		d.file.Close()
		d.file = nil
	}
	return nil
}
