package svg

import "github.com/gabrielluizsf/antui/canvas"

// clipRegionOf is the shapes of one clipPath moved into the place of the
// element that named it: with the element's own transform put on them, so
// that the cut follows the element it cuts, or left as they were written
// where the shapes are written in fractions of the element's own box and it
// is [clipUnder] that puts them where the box is. Every copy carries the
// clips cut into the shapes along with it — see [mapMaskShapes].
func (img *Image) clipRegionOf(cp *clipPath, st Style) *clipRegion {
	if cp.boxUnits {
		// The clipPath uses objectBoundingBox units: the shapes are fractions
		// of the element's box and are moved by clipUnder, once the element is
		// built and the box is known — not by the element's own transform,
		// which would put them somewhere before they were ever mapped.
		return &clipRegion{shapes: mapMaskShapes(cp.shapes, canvas.Identity()), boxUnits: true}
	}
	m := canvas.Identity()
	if st.HasTransform && st.Transform != (canvas.Matrix{}) {
		// The element's own transform is on the shapes now, the same way it is
		// put on the element's own path while it is built, so that what cuts
		// the element moves with the element. Everything that points at one
		// clipPath takes its own copy for exactly this reason: the drawing is
		// read once and may be painted at many sizes and in many places.
		m = st.Transform
	}
	return &clipRegion{shapes: mapMaskShapes(cp.shapes, m)}
}
