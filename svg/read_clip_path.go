package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// readClipPath reads one `<clipPath>`: how it said its shapes were to be
// measured, and the picture they are to be read out of. The picture is built
// in the coordinates it was written in — the transform on the clipPath and on
// anything inside it counts, since that is part of how it was written — and
// the shapes themselves are worked out afterwards, once every clipPath in the
// drawing is named; see [Image.readClipPaths]. The style it is built from
// starts at the drawing's own rather than at the style of the place the
// clipPath stands, because a clipPath takes nothing from the elements around
// it except the clip-rule it was written under.
//
// The children are built rather than left out, which is what forceKids says:
// what is inside a definition is not part of the picture where it stands —
// that is what being a definition means — but here it is not being painted
// either, it is being read as the shape of a cut.
func (img *Image) readClipPath(e *element, inherited canvas.FillRule) *clipPath {
	base := img.rootStyle()
	base.ClipRule = inherited
	units := strings.ToLower(strings.TrimSpace(e.attr("clipPathUnits")))
	return &clipPath{
		boxUnits: units != "" && units != "userspaceonuse",
		content:  img.build(e, base, true),
	}
}
