package svg

import "github.com/gabrielluizsf/antui/canvas"

// strokeBand widens a path into the band a stroke of it covers, in the
// coordinates of the canvas. The widening is done in the coordinates the shape
// was written in, where a stroke is one width all along it, and the transform
// that moved the shape carries the band back — so a transform that stretches one
// axis stretches the stroke with it, rather than the width being read off the
// average of the two and coming out between them.
//
// `non-scaling-stroke` is the other reading and the one exception: a width in
// the pixels the drawing comes out at, with nothing scaled at all, which is why
// it is widened after the transform with the supersampled width. A transform
// with no way back cannot be undone to widen in, so the width falls back to the
// average there rather than to nothing.
func strokeBand(path *canvas.Path, st Style, m canvas.Matrix) *canvas.Path {
	style := canvas.StrokeStyle{Width: st.Width, Cap: st.Cap, Join: st.Join, MiterLimit: st.MiterLimit}
	if st.VectorEffect.NonScalingStroke {
		scaled := clonePath(path)
		scaled.Transform(m)
		style.Width = st.Width * supersample
		return canvas.StrokeOutline(scaled, style)
	}
	inv, ok := st.Transform.Inverse()
	if !ok {
		scaled := clonePath(path)
		scaled.Transform(m)
		style.Width = st.Width * scaleOf(m.Mul(st.Transform))
		return canvas.StrokeOutline(scaled, style)
	}
	local := clonePath(path)
	local.Transform(inv)
	band := canvas.StrokeOutline(local, style)
	if band == nil {
		return nil
	}
	band.Transform(st.Transform)
	band.Transform(m)
	return band
}
