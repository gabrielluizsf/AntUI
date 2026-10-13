//go:build android

package ndk

/*
#cgo LDFLAGS: -lmediandk -landroid

#include "media_cgo.h"
*/
import "C"

import "fmt"

// open makes an extractor of this decoder's own and selects its track. Each
// decoder gets one: an extractor has a single position, and two tracks read
// at different rates cannot share it.
func (d *Decoder) open(path string, track Track) error {
	ex, f, err := newExtractor(path)
	if err != nil {
		return err
	}
	d.media.ex = ex
	d.file = f
	if r := C.AMediaExtractor_selectTrack(d.media.ex, C.size_t(track.Index)); r != C.AMEDIA_OK {
		return fmt.Errorf("ndk: cannot select track %d (%d)", track.Index, int(r))
	}
	return nil
}
