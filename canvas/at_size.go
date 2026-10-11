package canvas

import "math"

// AtSize is this face at a size in pixels per em, made once and kept. Where
// [Face.Scaled] asks for a multiple of the size the face already has, this one
// asks for the size itself, which is what a caller measuring in its own units
// wants: writing twelve pixels high and writing thirty-six both take the same
// letters, drawn at the size asked for.
//
// The size is rounded to a whole pixel, because a glyph is rasterised once and
// the cache it is kept in is keyed by the size it was made at. The built-in
// font is one shape cut at eight by sixteen, so its sizes are not steps any
// more: the shape is drawn with its cells at whatever size comes out, which is
// what makes a twelve pixel font twelve pixels high rather than sixteen.
func (f *Face) AtSize(pixels float64) *Face {
	if pixels <= 0 {
		return f
	}
	if f == nil || f.file == nil {
		return builtinAt(pixels)
	}
	n := int(math.Round(pixels))
	if n == int(math.Round(f.size)) {
		return f
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.sized == nil {
		f.sized = map[int]*Face{}
	}
	if d, ok := f.sized[n]; ok {
		return d
	}
	d, err := ParseFace(f.file.Data, float64(n))
	if err != nil {
		return f
	}
	if f.next != nil {
		d.next = f.next.AtSize(float64(n))
	}
	f.sized[n] = d
	return d
}
