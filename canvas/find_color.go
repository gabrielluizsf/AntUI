package canvas

import "github.com/gabrielluizsf/antui/font"

// findColor is the colour picture for a rune without the cache: this file's
// own when it holds a glyph for the rune, and otherwise the picture on the
// first face down the chain that holds one. The rune stops at the first face
// with a glyph for it, the same as an outline does — the names after that one
// in a font-family are for the runes none of the ones before them hold.
func (f *Face) findColor(r rune) *font.ColorGlyph {
	if gid := f.file.Cmap.Glyph(r); gid != 0 {
		if f.hasColor {
			return f.file.ColorGlyph(gid, f.size)
		}
		return nil
	}
	for i, n := 0, f.next; n != nil && i < 8; i, n = i+1, n.next {
		if n.file == nil {
			continue
		}
		if gid := n.file.Cmap.Glyph(r); gid != 0 {
			return n.colorGlyph(r)
		}
	}
	return nil
}
