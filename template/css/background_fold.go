package css

import (
	"github.com/gabrielluizsf/antui/canvas"
)

// backgroundFold accumulates the per-layer pieces of a background shorthand
// before they are stored on a style.
type backgroundFold struct {
	images   []BackImage
	poss     []BackPos
	sizes    []BackSize
	repeats  []BackRepeat
	clips    []uint8
	origins  []uint8
	attaches []uint8
	color    canvas.Color
	colorSet bool
}

// store writes the folded background onto the style.
func (f *backgroundFold) store(st *Style) {
	if f.colorSet {
		st.Background = f.color
	}
	st.BackgroundImages = f.images
	st.BackgroundPos = f.poss
	st.BackgroundSize = f.sizes
	st.BackgroundRepeat = f.repeats
	st.BackgroundClip = f.clips
	st.BackgroundOrigin = f.origins
	st.BackgroundAttach = f.attaches
}
