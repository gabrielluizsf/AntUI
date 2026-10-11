package font

// HasColorGlyphs reports whether the file carries colour bitmaps, which is
// what a face asks before it looks for one on any rune.
func (f *TTF) HasColorGlyphs() bool { return len(f.colorStrikes) > 0 }
