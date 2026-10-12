//go:build android

package ndk

/*
#cgo LDFLAGS: -lmediandk -landroid

#include "media_cgo.h"
*/
import "C"

import (
	"fmt"
	"time"
)

// OpenMedia reads a file's table of contents. It decodes nothing: that is
// [Media.DecodeVideo] and [Media.DecodeAudio], each of which opens the file
// again for a track of its own.
func OpenMedia(path string) (*Media, error) {
	ex, f, err := newExtractor(path)
	if err != nil {
		return nil, err
	}
	defer C.AMediaExtractor_delete(ex)
	defer f.Close()

	n := int(C.AMediaExtractor_getTrackCount(ex))
	m := &Media{path: path, tracks: make([]Track, 0, n)}
	for i := range n {
		format := C.AMediaExtractor_getTrackFormat(ex, C.size_t(i))
		if format == nil {
			continue
		}
		t := Track{Index: i}
		t.Mime = formatString(format, "mime")
		if us, ok := formatInt64(format, "durationUs"); ok {
			t.Duration = time.Duration(us) * time.Microsecond
		}
		t.Width, _ = formatInt(format, "width")
		t.Height, _ = formatInt(format, "height")
		t.Rotation, _ = formatInt(format, "rotation-degrees")
		t.FrameRate, _ = formatInt(format, "frame-rate")
		t.SampleRate, _ = formatInt(format, "sample-rate")
		t.Channels, _ = formatInt(format, "channel-count")
		C.AMediaFormat_delete(format)
		m.tracks = append(m.tracks, t)
	}
	if len(m.tracks) == 0 {
		return nil, fmt.Errorf("ndk: %s has no tracks this device can read", path)
	}
	return m, nil
}
