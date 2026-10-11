package canvas

// BuiltinScaled is the built-in font with each of its pixels drawn n times
// over, which is a size the bitmap font comes in.
func BuiltinScaled(n int) *Face {
	if n <= 1 {
		return builtin
	}
	n = min(n, 16)
	builtinMu.Lock()
	defer builtinMu.Unlock()
	if f, ok := builtinScaled[n]; ok {
		return f
	}
	f := &Face{
		unit:    float64(n),
		ascent:  builtin.ascent * n,
		descent: builtin.descent * n,
		height:  FontHeight * n,
	}
	builtinScaled[n] = f
	return f
}
