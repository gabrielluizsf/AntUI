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
)

// DecodeVideo starts decoding a picture track into frames.
//
// rotate is how far to turn each frame, in degrees clockwise. Passing the
// track's own [Track.Rotation] is what makes a video shot on a phone held
// sideways come out the right way up; passing 0 gives it as stored.
//
// It needs API 24, because the frames arrive through an image reader.
func (m *Media) DecodeVideo(track Track, rotate int) (*Decoder, error) {
	if !track.Video() {
		return nil, fmt.Errorf("ndk: track %d is %s, not a picture", track.Index, track.Mime)
	}
	if track.Width <= 0 || track.Height <= 0 {
		return nil, fmt.Errorf("ndk: track %d does not say how big it is", track.Index)
	}
	d := &Decoder{track: track, rotate: rotate}
	if err := d.open(m.path, track); err != nil {
		return nil, err
	}
	// The reader is the surface the codec draws into, at the video's own
	// size — scaling here would be scaling before anything has been seen.
	if r := C.mw_image_reader_new(C.int32_t(track.Width), C.int32_t(track.Height),
		C.int32_t(FormatYUV420), 3, &d.media.reader); r != C.AMEDIA_OK {
		d.Close()
		return nil, ErrNoImageReader
	}
	if r := C.mw_image_reader_window(d.media.reader, &d.media.window); r != C.AMEDIA_OK {
		d.Close()
		return nil, errors.New("ndk: the image reader has no surface")
	}
	C.ANativeWindow_acquire(d.media.window)

	if err := d.configure(d.media.window); err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}
