//go:build windows

package windows

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// iconFor makes an HICON out of whichever picture is nearest a size.
func iconFor(images []*canvas.Canvas, want int) uintptr {
	best := images[0]
	for _, image := range images {
		if awayFrom(image.Width, want) < awayFrom(best.Width, want) {
			best = image
		}
	}
	return hicon(best)
}
