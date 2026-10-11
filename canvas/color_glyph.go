package canvas

import "github.com/gabrielluizsf/antui/font"

// colorGlyph is the colour picture drawn for a rune here: this file's own
// when it has one, the picture on the face the rune falls to when this file
// has no glyph for the rune at all, and nil when the rune is drawn as an
// outline — or when no face in the chain carries a picture for it. The answer
// is kept, so a frame asks once per rune however many times it draws it.
func (f *Face) colorGlyph(r rune) *font.ColorGlyph {
	if f == nil || f.file == nil {
		return nil
	}
	f.mu.Lock()
	cg, seen := f.colors[r]
	f.mu.Unlock()
	if seen {
		return cg
	}
	cg = f.findColor(r)
	f.mu.Lock()
	if f.colors == nil {
		f.colors = map[rune]*font.ColorGlyph{}
	}
	f.colors[r] = cg
	f.mu.Unlock()
	return cg
}
