package svg

import (
	"bytes"
	"fmt"
	"image"
)

// checkPictureSize reads only the header of a picture and refuses one that
// says it is bigger than [maxPicturePixels], before the decoder is given the
// chance to allocate what it says. frames is how many pictures of that size
// the file holds, one for anything that does not animate.
func checkPictureSize(data []byte, frames int) error {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("it is not a drawing, a PNG, a JPEG or a GIF: %w", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return fmt.Errorf("it has no size to it")
	}
	if cfg.Width > maxPicturePixels/cfg.Height/max(1, frames) {
		return fmt.Errorf("it says it is %dx%d pixels, more than the %d a picture may be", cfg.Width, cfg.Height, maxPicturePixels)
	}
	return nil
}
