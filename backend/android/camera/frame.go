//go:build android

package camera

import (
	"errors"
	"unsafe"

	"github.com/gabrielluizsf/antui/canvas"
)

// Frame is the newest picture from the camera, or false when none has
// arrived since the last call.
//
// The canvas belongs to the preview and is reused: a 1080p frame is eight
// megabytes, and making a new one thirty times a second would do nothing but
// keep the collector busy. Copy anything that has to outlive the next call.
func (p *Preview) Frame() (*canvas.Canvas, bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil, false, errors.New("camera: the preview is closed")
	}
	if p.canvas == nil {
		if err := p.makeCanvas(p.askedW, p.askedH); err != nil {
			return nil, false, err
		}
	}
	cv := p.canvas
	// canvas.Color is a uint32 underneath, and the converter writes words.
	words := unsafe.Slice((*uint32)(unsafe.Pointer(&cv.Pixels[0])), len(cv.Pixels))
	w, h, ok, err := p.device.Frame(words, cv.Stride, p.upright)
	if err != nil || !ok {
		return nil, false, err
	}
	if w != cv.Width || h != cv.Height {
		// The camera settled on something other than what was asked for.
		// That frame is lost — it was written into a canvas the wrong shape
		// — and the next one lands in a canvas that fits.
		p.askedW, p.askedH = w, h
		p.canvas = nil
		return nil, false, nil
	}
	return cv, true, nil
}
