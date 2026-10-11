package antui

import (
	"fmt"
	"github.com/gabrielluizsf/antui/canvas"
)

// iconPictures turns canvases into what the formats above want, largest
// first. A size over 256 is left out: the entry that names it has one byte
// for the width, and 256 is written as a nought.
func iconPictures(images []*canvas.Canvas) ([]iconPicture, error) {
	var out []iconPicture
	for _, image := range images {
		if image == nil || image.Width <= 0 || image.Height <= 0 {
			continue
		}
		if image.Width > 256 || image.Height > 256 {
			continue
		}
		data, err := encodeIconPNG(image)
		if err != nil {
			return nil, err
		}
		out = append(out, iconPicture{image.Width, image.Height, data})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("antui: there is no picture to make an icon from")
	}
	// Largest first, which is what every reader here expects.
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].width > out[j-1].width; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out, nil
}
