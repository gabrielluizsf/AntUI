//go:build android

package media

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// prepareCanvases makes the buffers a video track is drawn through.
func prepareCanvases(p *Player) error {
	if p.video == nil {
		return nil
	}
	cv, err := canvas.NewCanvas(p.info.Width, p.info.Height)
	if err != nil {
		p.Close()
		return err
	}
	p.cur = cv
	if cv, err = canvas.NewCanvas(p.info.Width, p.info.Height); err == nil {
		p.next = cv
	}
	return nil
}
