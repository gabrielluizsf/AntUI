package canvas

import (
	"fmt"
	"math"
	"os"
	"sync"
	"sync/atomic"

	"github.com/gabrielluizsf/antui/font"
)

// Face is a font at a size, ready to draw. Two kinds exist:
//
//	canvas.BuiltinFace()          the 8x16 bitmap, exactly as before
//	antui.SystemFace(15)          whatever this platform writes its own
//	                               interface in, at 15 pixels
//	canvas.LoadFace(path, 15)     a TrueType file of your own
//
// Set a face for the whole program:
//
//	canvas.SetDefaultFace(face)   // every Text, TextWidth and widget
type Face struct {
	// built-in, when font is nil; bitmap is how many times over each of its
	// pixels is drawn, which is the only size a bitmap font has.
	bitmap int
	file   *font.TTF
	size   float64 // pixels per em
	scale  float64 // font units to pixels

	ascent, descent, height int

	// next is the face handed a rune this one has no shape for — the rest
	// of a font-family list, then the face the program draws with by
	// default. It is set when the chain is made and never after, and is
	// read only by a face that has a file behind it: the bitmap font has no
	// outlines to fall back from, it is the end of every chain.
	next *Face

	// colors is the colour picture a rune draws here — nil for the ones
	// drawn as an outline — read off this file or off a face further down
	// the chain, and asked for once per rune.
	colors   map[rune]*font.ColorGlyph
	hasColor bool

	mu      sync.Mutex
	glyphs  map[rune]*font.Mask
	derived map[int]*Face
	// sized is the same face at sizes asked for outright rather than as a
	// multiple, keyed by the whole pixel it was made at. See [Face.AtSize].
	sized map[int]*Face
}

// builtin is the 8x16 bitmap, as a Face.
var builtin = &Face{
	bitmap:  1,
	ascent:  12,
	descent: FontHeight - 12,
	height:  FontHeight,
}

// BuiltinFace is the 8x16 bitmap font that has always been here.
func BuiltinFace() *Face { return builtin }

// BuiltinScaled is the built-in font with each of its pixels drawn n times
// over, which is the only size a bitmap font has.
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
		bitmap:  n,
		ascent:  builtin.ascent * n,
		descent: builtin.descent * n,
		height:  FontHeight * n,
	}
	builtinScaled[n] = f
	return f
}

var (
	builtinMu     sync.Mutex
	builtinScaled = map[int]*Face{}
)

// defaultFace is what Text and TextWidth use when nothing says otherwise.
var defaultFace atomic.Pointer[Face]

// SetDefaultFace makes every Text, TextWidth and widget in the program use
// this face. nil goes back to the built-in.
func SetDefaultFace(f *Face) { defaultFace.Store(f) }

// DefaultFace is what is being used now, and is never nil.
func DefaultFace() *Face {
	if f := defaultFace.Load(); f != nil {
		return f
	}
	return builtin
}

// LoadFace reads a TrueType file at a size in pixels per em.
func LoadFace(path string, pixels float64) (*Face, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("antui: %w", err)
	}
	return ParseFace(body, pixels)
}

// ParseFace is [LoadFace] for a font already in memory.
func ParseFace(data []byte, pixels float64) (*Face, error) {
	if pixels <= 0 {
		return nil, fmt.Errorf("antui: %v is not a size", pixels)
	}
	f, err := font.ParseTTF(data)
	if err != nil {
		return nil, err
	}
	scale := f.Scaled(pixels)
	return &Face{
		file:     f,
		size:     pixels,
		scale:    scale,
		ascent:   int(float64(f.Ascent)*scale + 0.5),
		descent:  int(-float64(f.Descent)*scale + 0.5),
		height:   f.LineHeight(pixels),
		hasColor: f.HasColorGlyphs(),
		glyphs:   map[rune]*font.Mask{},
	}, nil
}

// Size is the face's size in pixels per em, and 0 for the built-in.
func (f *Face) Size() float64 {
	if f == nil || f.file == nil {
		return 0
	}
	return f.size
}

// Height is how far apart two lines of this face are, in pixels.
func (f *Face) Height() int {
	if f == nil {
		return builtin.height
	}
	return f.height
}

// Ascent is how far the tallest letters reach above the baseline, and
// Descent how far the tails go below it.
func (f *Face) Ascent() int {
	if f == nil {
		return builtin.ascent
	}
	return f.ascent
}

func (f *Face) Descent() int {
	if f == nil {
		return builtin.descent
	}
	return f.descent
}

// Fixed reports whether every character is the same width.
func (f *Face) Fixed() bool { return f == nil || f.file == nil }

func (f *Face) scaleOf() int {
	if f == nil || f.bitmap <= 0 {
		return 1
	}
	return f.bitmap
}

// Width is how wide a string is, in pixels, measuring the widest line.
func (f *Face) Width(text string) int {
	if f == nil || f.file == nil {
		return font.TextWidth(text) * f.scaleOf()
	}
	widest, line := 0, 0.0
	for _, r := range text {
		switch r {
		case '\n':
			widest = max(widest, int(line+0.5))
			line = 0
		case '\r':
		case '\t':
			line += f.advance(' ') * 4
		default:
			line += f.advance(r)
		}
	}
	return max(widest, int(line+0.5))
}

// advance is how far one rune moves the pen, in this face's pixels: its own
// advance when it holds the rune, and the advance of the face in the chain
// that does when this one does not.
func (f *Face) advance(r rune) float64 {
	if f == nil || f.file == nil {
		return float64(font.TextWidth(string(r))) * float64(f.scaleOf())
	}
	if cg := f.colorGlyph(r); cg != nil {
		return float64(cg.Advance) * f.colorScale(cg)
	}
	g := f.file.Cmap.Glyph(r)
	if g == 0 {
		if n := f.fallbackFor(r); n != nil {
			return n.advance(r)
		}
	}
	return float64(f.file.Advance(g)) * f.scale
}

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

// colorScale is how much a picture cut at its strike's size is enlarged by to
// draw at this face's own size.
func (f *Face) colorScale(cg *font.ColorGlyph) float64 {
	if cg == nil || cg.PPEM <= 0 {
		return 1
	}
	return f.size / float64(cg.PPEM)
}

// blitColor draws one colour picture into a canvas with the pen at x and the
// baseline at y: scaled from the strike it was cut at to the size this face
// draws at, with its left edge bearingX to the right of the pen and its top
// bearingY above the baseline. The weight and slant of a style lean and
// double the picture the same way they do an outline.
func (f *Face) blitColor(cv *Canvas, cg *font.ColorGlyph, x, y float64, o TextStyle) {
	k := f.colorScale(cg)
	bounds := cg.Image.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	dstW, dstH := max(1, int(float64(srcW)*k+0.5)), max(1, int(float64(srcH)*k+0.5))
	left := int(math.Round(x + float64(cg.BearingX)*k))
	top := int(math.Round(y - float64(cg.BearingY)*k))
	step := max(o.Scale, 1)
	for row := range dstH {
		sy := bounds.Min.Y + row*srcH/dstH
		shift := 0
		if o.Italic {
			shift = (dstH - 1 - row) / 4
		}
		for col := range dstW {
			r, g, blue, a := cg.Image.At(bounds.Min.X+col*srcW/dstW, sy).RGBA()
			if a == 0 {
				continue
			}
			// A decoder hands the channels back already multiplied by the
			// alpha; the canvas takes them plain.
			c := RGBA(uint8(r*255/a), uint8(g*255/a), uint8(blue*255/a), uint8(a>>8))
			px := left + col + shift
			cv.Pixel(px, top+row, c)
			if o.Bold {
				cv.Pixel(px+step, top+row, c)
			}
		}
	}
}

// fallbackFor is the face of the chain that has a shape for r, past this one,
// and nil when none of them does. Eight is further down than any chain is
// built, so one written round in a circle stops rather than running forever.
func (f *Face) fallbackFor(r rune) *Face {
	for i, n := 0, f.next; n != nil && i < 8; i, n = i+1, n.next {
		if n.file != nil && n.file.Cmap.Glyph(r) != 0 {
			return n
		}
	}
	return nil
}

func (f *Face) glyph(r rune) *font.Mask {
	if f.file.Cmap.Glyph(r) == 0 {
		if n := f.fallbackFor(r); n != nil {
			return n.glyph(r)
		}
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if m, ok := f.glyphs[r]; ok {
		return m
	}
	g := f.file.Cmap.Glyph(r)
	m := font.Render(f.file.GlyphContours(g, 0), f.scale)
	f.glyphs[r] = &m
	return &m
}

// WithFallback is this face with next handed every rune this one has no shape
// for, the way a font-family list hands a character down its names to the face
// that has it. The result draws differently from what this one draws, so it
// cannot share the glyph cache; it does share the font file, and the caches of
// sizes at which either of them is met.
func (f *Face) WithFallback(next *Face) *Face {
	if f == nil {
		return next
	}
	if next == nil {
		return f
	}
	return &Face{
		bitmap:   f.bitmap,
		file:     f.file,
		size:     f.size,
		scale:    f.scale,
		ascent:   f.ascent,
		descent:  f.descent,
		height:   f.height,
		hasColor: f.hasColor,
		next:     next,
		glyphs:   map[rune]*font.Mask{},
	}
}

// Draw writes a string into a canvas with the baseline of the first line at
// y, and returns the width of the widest line.
func (f *Face) Draw(cv *Canvas, x, y int, text string, c Color) int {
	if f == nil || f.file == nil {
		return cv.textBitmap(x, y-f.Ascent(), text, c, f.scaleOf())
	}
	if cv == nil {
		return f.Width(text)
	}
	penX, penY := float64(x), y
	widest := 0
	for _, r := range text {
		switch r {
		case '\n':
			widest = max(widest, int(penX-float64(x)+0.5))
			penX, penY = float64(x), penY+f.height
			continue
		case '\r':
			continue
		case '\t':
			penX += f.advance(' ') * 4
			continue
		case ' ':
			penX += f.advance(' ')
			continue
		}
		if cg := f.colorGlyph(r); cg != nil {
			f.blitColor(cv, cg, penX, float64(penY), TextStyle{})
		} else if m := f.glyph(r); m.W > 0 {
			blitMask(cv, m, int(penX+0.5), penY, c)
		}
		penX += f.advance(r)
	}
	return max(widest, int(penX-float64(x)+0.5))
}

// blitMask composites one glyph's coverage onto a canvas in a colour.
func blitMask(cv *Canvas, m *font.Mask, x, baseline int, c Color) {
	alpha := int(c.A())
	if alpha == 0 {
		return
	}
	top := baseline - m.Top
	for row := range m.H {
		for col := range m.W {
			v := int(m.Alpha[row*m.Stride+col])
			if v == 0 {
				continue
			}
			cv.Pixel(x+m.Left+col, top+row, RGBA(c.R(), c.G(), c.B(),
				uint8(v*alpha/255)))
		}
	}
}

// Scaled is the same font at a multiple of its size, made once and kept.
func (f *Face) Scaled(n int) *Face {
	if n <= 1 {
		return f
	}
	if f == nil {
		return BuiltinScaled(n)
	}
	if f.file == nil {
		return BuiltinScaled(f.scaleOf() * n)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.derived == nil {
		f.derived = map[int]*Face{}
	}
	if d, ok := f.derived[n]; ok {
		return d
	}
	d, err := ParseFace(f.file.Data, f.size*float64(n))
	if err != nil {
		return f
	}
	if f.next != nil {
		d.next = f.next.Scaled(n)
	}
	f.derived[n] = d
	return d
}

// AtSize is this face at a size in pixels per em, made once and kept. Where
// [Face.Scaled] asks for a multiple of the size the face already has, this one
// asks for the size itself, which is what a caller measuring in its own units
// wants: writing twelve pixels high and writing thirty-six both take the same
// letters, drawn at the size asked for.
//
// The size is rounded to a whole pixel, because a glyph is rasterised once and
// the cache it is kept in is keyed by the size it was made at. The built-in
// font has one shape drawn n times over, so its sizes come out in steps of its
// own height and anything below it asks for the one it has.
func (f *Face) AtSize(pixels float64) *Face {
	if pixels <= 0 {
		return f
	}
	if f == nil || f.file == nil {
		return BuiltinScaled(max(1, int(math.Round(pixels/float64(FontHeight)))))
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

// SetFace makes this one canvas draw text in a face of its own, whatever the
// program's default is. nil goes back to the default.
func (cv *Canvas) SetFace(f *Face) {
	if cv != nil {
		cv.font = f
		cv.hasFont = true
		if f == nil {
			cv.hasFont = false
		}
	}
}

// Face is what this canvas draws text in.
func (cv *Canvas) Face() *Face { return cv.face() }

// face is the face to use, never nil.
func (cv *Canvas) face() *Face {
	if cv != nil && cv.hasFont && cv.font != nil {
		return cv.font
	}
	return DefaultFace()
}
