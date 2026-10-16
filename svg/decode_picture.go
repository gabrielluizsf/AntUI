package svg

import (
	"bytes"
	"fmt"
	"image"

	"github.com/gabrielluizsf/antui/canvas"
)

// decodePicture is a PNG, a JPEG or a GIF turned into the frames it draws —
// one for a picture that does not move, every frame of the loop for a GIF that
// animates, so a paint picks by the clock which one is showing. The channels
// come out of the decoder premultiplied by their alpha and a canvas holds them
// straight, which is what [pictureCanvas] undoes.
func decodePicture(data []byte) (*pictureFrames, error) {
	src, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("it is not a PNG, a JPEG or a GIF: %w", err)
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
	return &pictureFrames{frames: []*canvas.Canvas{cv}}, nil
}
