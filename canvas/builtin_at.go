package canvas

import "math"

// builtinAt is the built-in font at a size in pixels, made once and kept. The
// sizes that are a whole multiple of its own height are the faces
// [BuiltinScaled] already hands out — the same shape drawn n times over — and
// any other size is that shape with its cells at the size asked for, so twelve
// pixels of writing come out twelve pixels high rather than being rounded up
// to the sixteen the shape was cut at.
func builtinAt(pixels float64) *Face {
	n := int(math.Round(pixels))
	if n <= 0 {
		return builtin
	}
	if n >= FontHeight && n%FontHeight == 0 && n/FontHeight <= 16 {
		return BuiltinScaled(n / FontHeight)
	}
	builtinMu.Lock()
	defer builtinMu.Unlock()
	if f, ok := builtinSized[n]; ok {
		return f
	}
	unit := float64(n) / float64(FontHeight)
	f := &Face{
		unit:    unit,
		ascent:  int(math.Round(float64(builtin.ascent) * unit)),
		descent: int(math.Round(float64(builtin.descent) * unit)),
		height:  n,
	}
	builtinSized[n] = f
	return f
}
