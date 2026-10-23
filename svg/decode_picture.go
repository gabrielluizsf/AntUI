package svg

import (
	"bytes"
	"fmt"
	"image"

	"github.com/gabrielluizsf/antui/canvas"
)

// decodePicture is what an address answered with: the frames it draws — one
// for a picture that does not move, every frame of the loop for a GIF that
// animates — either a photograph read by the decoders of the standard library
// or a drawing read by this package's own. chain is the drawings being read
// through one `<image>` after another to get here, which is what keeps a
// circle of them from reading forever, and target is the address they were
// read from.
func decodePicture(chain []string, target string, data []byte) (*pictureFrames, error) {
	if sniffSVG(data) {
		return svgFrames(chain, target, data)
	}
	if err := checkPictureSize(data, 1); err != nil {
		return nil, err
	}
	src, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("it is not a drawing, a PNG, a JPEG or a GIF: %w", err)
	}
	if format == "gif" {
		// The first frame is all `image.Decode` answers; the animation is the
		// whole of what a GIF is, so the bytes go through the GIF reader that
		// takes every frame of it.
		return gifFrames(data)
	}
	cv, err := pictureCanvas(src)
	if err != nil {
		return nil, err
	}
	return &pictureFrames{frames: []*canvas.Canvas{cv}, ownW: float64(cv.Width), ownH: float64(cv.Height)}, nil
}
