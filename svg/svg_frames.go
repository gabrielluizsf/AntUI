package svg

import (
	"fmt"
	"math"
	"slices"

	"github.com/gabrielluizsf/antui/canvas"
)

// maxPictureDepth is how many drawings may be read through one `<image>` after
// another before the reading is refused. A circle is stopped by the chain
// itself — the drawing whose address is already being read through is not read
// again — but an address that answers with a new one every time, which a
// server is free to do, has no circle to stop and only a bound to keep.
const maxPictureDepth = 16

// svgFrames is a drawing read as a picture: the bytes are read whole as the
// document they name — with the address it came from as the base its own
// relative addresses resolve against — and painted at the size the document
// says it is, which is the size a picture of it has. What goes wrong inside
// it stays inside it: a picture it cannot read is a hole in it rather than a
// warning in the drawing pointing at it, and the drawing is painted all the
// same.
func svgFrames(chain []string, target string, data []byte) (*pictureFrames, error) {
	if slices.Contains(chain, target) {
		return nil, fmt.Errorf("it is one of the drawings it is read through")
	}
	if len(chain) >= maxPictureDepth {
		return nil, fmt.Errorf("it is read through more than %d drawings deep", maxPictureDepth)
	}
	img, err := parseText(string(data), target, chain)
	if err != nil {
		return nil, err
	}
	w, h := max(1, int(math.Round(img.Width))), max(1, int(math.Round(img.Height)))
	cv := img.Render(w, h)
	if cv == nil {
		return nil, fmt.Errorf("it has no size to it")
	}
	return &pictureFrames{
		frames:    []*canvas.Canvas{cv},
		ownW:      img.Width,
		ownH:      img.Height,
		docFit:    img.docFit,
		docAlign:  img.docAlign,
		hasDocFit: img.hasDocFit,
	}, nil
}
