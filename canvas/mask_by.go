package canvas

// MaskBy takes the alpha of every pixel of cv down to what the picture in m
// keeps of it, which is a multiply and not a cut: where m keeps half, what cv
// holds comes out half, and where m keeps nothing nothing of it is left. That
// is what makes a mask different from [Canvas.MaskShapes], which decides
// whether a pixel is there at all — a mask says how much of it is, and a
// picture painted white over half a shape leaves that half whole and takes the
// other half away rather than putting a hard edge between them.
//
// How much is kept is the mode's business: how bright the pixel of m is, by
// the same coefficients [Canvas.FilterRegion] grays a colour with, and how much
// of that pixel is there on top of it — a clear pixel keeps nothing however
// bright the colour behind its zero alpha happens to be. In MaskAlpha only the
// second half counts.
//
// A pixel of cv that m has no say about — one beside a picture smaller than cv
// — is under no mask at all and so keeps nothing, the same as the union of no
// shapes empties the clip in [Canvas.MaskShapes]. Only the alpha moves: the
// colours are as they were and only how much of each there is changes, which
// is what makes it a call to make on a layer. The narrow formats have no
// per-pixel alpha to multiply and are left alone.
func (cv *Canvas) MaskBy(m *Canvas, mode MaskMode) {
	if cv.Pixels == nil || m == nil {
		return
	}
	for y := cv.Clip.Y; y < cv.Clip.Y+cv.Clip.Height; y++ {
		for x := cv.Clip.X; x < cv.Clip.X+cv.Clip.Width; x++ {
			keep := maskKeep(m.At(x, y), mode)
			c := cv.At(x, y)
			a := int(float64(c.A())*keep + 0.5)
			if a == int(c.A()) {
				continue
			}
			cv.Put(x, y, Fade(c, a))
		}
	}
}
