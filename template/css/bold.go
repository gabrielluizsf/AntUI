package css

// Bold reports whether the file already carries the heavy weight, so a style
// that asked for bold has no fake bold to draw over it.
func (ff *FontFace) Bold() bool { return ff.Weight >= 600 }
