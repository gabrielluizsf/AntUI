package css

// Slanted reports whether the file already leans, so a style that asked for
// italic has no synthetic slant to add on top of it.
func (ff *FontFace) Slanted() bool { return ff.Style != FontStyleNormal }
