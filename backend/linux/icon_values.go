//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import "github.com/gabrielluizsf/antui/canvas"

// iconValues flattens the icon images into the CARDINAL list _NET_WM_ICON
// wants, skipping any image that would push the list past what fits in one
// request.
func iconValues(images []*canvas.Canvas, capHint, most int) []uint32 {
	values := make([]uint32, 0, capHint)
	for _, image := range images {
		need := 2 + image.Width*image.Height
		if len(values)+need > most {
			continue
		}
		values = append(values, uint32(image.Width), uint32(image.Height))
		for y := range image.Height {
			for x := range image.Width {
				values = append(values, uint32(image.At(x, y)))
			}
		}
	}
	return values
}
