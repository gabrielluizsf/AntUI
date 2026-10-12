//go:build android

package camera

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// makeCanvas makes the canvas frames are converted into.
//
// The size is known in advance and does not have to be discovered: the image
// reader was created at exactly this size, so that is what every frame is.
// A quarter turn swaps it.
func (p *Preview) makeCanvas(w, h int) error {
	if p.upright == 90 || p.upright == 270 {
		w, h = h, w
	}
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		return err
	}
	p.canvas = cv
	return nil
}
