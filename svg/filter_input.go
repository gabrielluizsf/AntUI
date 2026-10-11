package svg

import "github.com/gabrielluizsf/antui/canvas"

// filterInput is the picture a primitive's `in` or `in2` names: nothing at
// all means what the primitive stands beside — the picture so far for `in`,
// the picture the element drew for `in2` — and SourceGraphic, SourceAlpha and
// the results an earlier primitive left are the rest. A name nothing holds
// was said while the filter was being read; the blank picture is what it was
// said to mean.
func filterInput(name string, deflt, src *canvas.Canvas, withAlpha **canvas.Canvas, named map[string]*canvas.Canvas) *canvas.Canvas {
	switch {
	case name == "":
		return deflt
	case name == "SourceGraphic":
		return src
	case name == "SourceAlpha":
		if *withAlpha == nil {
			*withAlpha = alphaPicture(src)
		}
		return *withAlpha
	case name == "BackgroundImage" || name == "BackgroundAlpha":
		return blankPicture(src.Width, src.Height)
	default:
		if c := named[name]; c != nil {
			return c
		}
		return blankPicture(src.Width, src.Height)
	}
}
