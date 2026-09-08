package canvas

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/gabrielluizsf/antui/font"
)

// Face is a font at a size, ready to draw. Two kinds exist:
//
//	canvas.BuiltinFace()          the 8x16 bitmap, exactly as before
//	canvas.SystemFace(15)         whatever this platform writes its own
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

	mu      sync.Mutex
	glyphs  map[rune]*font.Mask
	derived map[int]*Face
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
		file:    f,
		size:    pixels,
		scale:   scale,
		ascent:  int(float64(f.Ascent)*scale + 0.5),
		descent: int(-float64(f.Descent)*scale + 0.5),
		height:  f.LineHeight(pixels),
		glyphs:  map[rune]*font.Mask{},
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

func (f *Face) advance(r rune) float64 {
	g := f.file.Cmap.Glyph(r)
	return float64(f.file.Advance(g)) * f.scale
}

func (f *Face) glyph(r rune) *font.Mask {
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
		m := f.glyph(r)
		if m.W > 0 {
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
	f.derived[n] = d
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
