package canvas

import "fmt"

// Area is a rectangle, used for clipping and for dirty regions.
type Area struct {
	X, Y, Width, Height int
}

// Format is how one pixel is stored in a canvas.
type Format int

// The pixel formats.
const (
	ARGB32 Format = iota // four bytes, the default
	RGB565               // two bytes: 5 red, 6 green, 5 blue, no alpha
	Pal8                 // one byte, an index into a 256-entry palette
)

// Bytes is how much one pixel takes in this format.
func (f Format) Bytes() int {
	switch f {
	case RGB565:
		return 2
	case Pal8:
		return 1
	default:
		return 4
	}
}

// String names the format.
func (f Format) String() string {
	switch f {
	case ARGB32:
		return "ARGB32"
	case RGB565:
		return "RGB565"
	case Pal8:
		return "PAL8"
	default:
		return fmt.Sprintf("Format(%d)", int(f))
	}
}

// Canvas is an ARGB pixel surface that everything is drawn into — the
// window's front buffer is one, and NewCanvas makes free-standing ones that
// need no window at all, which is what lets drawing be tested in a terminal.
//
// On an ARGB32 canvas Pixels is the buffer, row by row, and may be read and
// written directly. On the narrow formats Pixels is deliberately nil: code
// that reaches straight into it then draws nothing rather than writing four
// bytes into a two-byte pixel. Everything going through the methods works on
// all three.
type Canvas struct {
	// font is this canvas's own face, when it has one. See Canvas.SetFace:
	// the default is the program's, and hasFont is what tells "nothing set"
	// apart from "set to the built-in on purpose".
	font    *Face
	hasFont bool

	Pixels []Color // width*height pixels, 0xAARRGGBB; nil on narrow formats
	Width  int
	Height int
	Stride int  // pixels per row (>= Width)
	Clip   Area // the active clip region

	format  Format
	narrow  []byte  // the 16- or 8-bit pixels
	palette []Color // 256 entries, on Pal8
	nearest []byte  // 32K: an RGB555 colour to its palette index
	dither  bool

	// Premultiplied records whether each channel has already been multiplied
	// by its own alpha. The flag lives beside the pixels rather than beside
	// the image so that converting one twice — which darkens it by the alpha
	// again and looks exactly like the first conversion — cannot happen.
	Premultiplied bool
}

// NewCanvas makes an ARGB32 canvas with its pixels zeroed.
func NewCanvas(width, height int) (*Canvas, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("antui: canvas size %dx%d is not valid", width, height)
	}
	cv := &Canvas{
		Pixels: make([]Color, width*height),
		Width:  width,
		Height: height,
		Stride: width,
		Clip:   Area{0, 0, width, height},
		format: ARGB32,
	}
	return cv, nil
}

// NewCanvasFormat makes a canvas holding pixels in the given format. ARGB32
// is exactly NewCanvas. A Pal8 canvas starts on a grey ramp, so a program
// that forgets to set a palette gets a grey picture rather than a black one
// and can see what it forgot.
func NewCanvasFormat(width, height int, format Format) (*Canvas, error) {
	if format == ARGB32 {
		return NewCanvas(width, height)
	}
	if format != RGB565 && format != Pal8 {
		return nil, fmt.Errorf("antui: there is no format %d", int(format))
	}
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("antui: canvas size %dx%d is not valid", width, height)
	}
	cv := &Canvas{
		Width:  width,
		Height: height,
		Stride: width,
		Clip:   Area{0, 0, width, height},
		format: format,
		narrow: make([]byte, width*height*format.Bytes()),
		dither: true,
	}
	if format == Pal8 {
		ramp := make([]Color, 256)
		for i := range ramp {
			ramp[i] = RGB(uint8(i), uint8(i), uint8(i))
		}
		cv.palette = ramp
		_ = cv.SetPalette(ramp)
	}
	return cv, nil
}

// Format reports how this canvas stores a pixel.
func (cv *Canvas) Format() Format { return cv.format }

// SetDither turns ordered dithering on or off for a narrow canvas. It is on
// by default, which is where it costs nothing and is any use.
func (cv *Canvas) SetDither(on bool) { cv.dither = on }

// SetPalette gives a paletted canvas its colours and builds the table that
// maps a colour back to the nearest index.
func (cv *Canvas) SetPalette(entries []Color) error {
	if cv.format != Pal8 {
		return fmt.Errorf("antui: SetPalette on a %s canvas, not a paletted one", cv.format)
	}
	if len(entries) == 0 {
		return fmt.Errorf("antui: SetPalette was given no colours")
	}
	count := min(len(entries), 256)
	if cv.palette == nil {
		cv.palette = make([]Color, 256)
	}
	for i := range 256 {
		cv.palette[i] = entries[min(i, count-1)] | 0xFF000000
	}

	if cv.nearest == nil {
		cv.nearest = make([]byte, 32768)
	}
	for r := range 32 {
		for g := range 32 {
			for b := range 32 {
				red := r<<3 | r>>2
				green := g<<3 | g>>2
				blue := b<<3 | b>>2
				best, bestDist := 0, -1
				for i, p := range cv.palette {
					dr := red - int(p.R())
					dg := green - int(p.G())
					db := blue - int(p.B())
					d := dr*dr + dg*dg + db*db
					if bestDist < 0 || d < bestDist {
						bestDist, best = d, i
					}
				}
				cv.nearest[r<<10|g<<5|b] = byte(best)
			}
		}
	}
	return nil
}

// PaletteIndex is the nearest entry to a colour, or -1 when the canvas has
// no palette.
func (cv *Canvas) PaletteIndex(c Color) int {
	if cv.palette == nil {
		return -1
	}
	if cv.nearest != nil {
		return int(cv.nearest[rgb555(c)])
	}
	best, bestDist := 0, -1
	for i, p := range cv.palette {
		dr := int(c.R()) - int(p.R())
		dg := int(c.G()) - int(p.G())
		db := int(c.B()) - int(p.B())
		d := dr*dr + dg*dg + db*db
		if bestDist < 0 || d < bestDist {
			bestDist, best = d, i
		}
	}
	return best
}

// SetClip narrows drawing to a rectangle, intersected with the canvas.
func (cv *Canvas) SetClip(x, y, w, h int) {
	x0 := clampInt(x, 0, cv.Width)
	y0 := clampInt(y, 0, cv.Height)
	x1 := clampInt(x+w, 0, cv.Width)
	y1 := clampInt(y+h, 0, cv.Height)
	cv.Clip = Area{x0, y0, max(0, x1-x0), max(0, y1-y0)}
}

// ResetClip opens drawing back up to the whole canvas.
func (cv *Canvas) ResetClip() {
	cv.Clip = Area{0, 0, cv.Width, cv.Height}
}

func (cv *Canvas) insideClip(x, y int) bool {
	return x >= cv.Clip.X && y >= cv.Clip.Y &&
		x < cv.Clip.X+cv.Clip.Width && y < cv.Clip.Y+cv.Clip.Height
}

// ---------------------------------------------------------------------------
// Narrow colour
// ---------------------------------------------------------------------------

// To565 packs a colour into 16 bits. It rounds rather than truncating.
func To565(c Color) uint16 {
	r5 := (uint32(c.R())*31 + 127) / 255
	g6 := (uint32(c.G())*63 + 127) / 255
	b5 := (uint32(c.B())*31 + 127) / 255
	return uint16(r5<<11 | g6<<5 | b5)
}

// From565 unpacks a 16-bit colour, repeating the top bits into the bottom.
func From565(packed uint16) Color {
	r5 := uint32(packed>>11) & 0x1F
	g6 := uint32(packed>>5) & 0x3F
	b5 := uint32(packed) & 0x1F
	r := r5<<3 | r5>>2
	g := g6<<2 | g6>>4
	b := b5<<3 | b5>>2
	return Color(0xFF000000 | r<<16 | g<<8 | b)
}

// bayer is the classic 4x4 threshold matrix.
var bayer = [16]uint32{
	0, 8, 2, 10,
	12, 4, 14, 6,
	3, 11, 1, 9,
	15, 7, 13, 5,
}

// div255 divides by 255 without dividing: exact for everything below 65535.
func div255(t uint32) uint32 { return (t + (t >> 8) + 1) >> 8 }

func to565Dithered(c Color, x, y int) uint16 {
	t := bayer[(y&3)*4+(x&3)] * 16
	r5 := min(div255(uint32(c.R())*31+t), 31)
	g6 := min(div255(uint32(c.G())*63+t), 63)
	b5 := min(div255(uint32(c.B())*31+t), 31)
	return uint16(r5<<11 | g6<<5 | b5)
}

func rgb555(c Color) uint32 {
	return (uint32(c)>>9)&0x7C00 | (uint32(c)>>6)&0x03E0 | (uint32(c)>>3)&0x001F
}

// narrowRead and narrowWrite are what every narrow path is built out of.
// Neither checks the clip: the callers do.
func (cv *Canvas) narrowRead(x, y int) Color {
	at := y*cv.Stride + x
	switch {
	case cv.format == RGB565:
		lo := uint16(cv.narrow[at*2])
		hi := uint16(cv.narrow[at*2+1])
		return From565(lo | hi<<8)
	case cv.format == Pal8 && cv.palette != nil:
		return cv.palette[cv.narrow[at]]
	}
	return 0
}

func (cv *Canvas) narrowWrite(x, y int, c Color) {
	at := y*cv.Stride + x
	switch {
	case cv.format == RGB565:
		var packed uint16
		if cv.dither {
			packed = to565Dithered(c, x, y)
		} else {
			packed = To565(c)
		}
		cv.narrow[at*2] = byte(packed)
		cv.narrow[at*2+1] = byte(packed >> 8)
	case cv.format == Pal8 && cv.nearest != nil:
		cv.narrow[at] = cv.nearest[rgb555(c)]
	case cv.format == Pal8:
		cv.narrow[at] = byte(max(cv.PaletteIndex(c), 0))
	}
}

// narrowRun writes a whole run in one format, without a call per pixel.
func (cv *Canvas) narrowRun(x0, x1, y int, c Color) {
	switch cv.format {
	case RGB565:
		at := (y*cv.Stride + x0) * 2
		if cv.dither {
			for x := x0; x < x1; x++ {
				packed := to565Dithered(c, x, y)
				cv.narrow[at] = byte(packed)
				cv.narrow[at+1] = byte(packed >> 8)
				at += 2
			}
			return
		}
		packed := To565(c)
		lo, hi := byte(packed), byte(packed>>8)
		for x := x0; x < x1; x++ {
			cv.narrow[at] = lo
			cv.narrow[at+1] = hi
			at += 2
		}
	case Pal8:
		at := y*cv.Stride + x0
		var index byte
		if cv.nearest != nil {
			index = cv.nearest[rgb555(c)]
		} else {
			index = byte(max(cv.PaletteIndex(c), 0))
		}
		row := cv.narrow[at : at+(x1-x0)]
		for i := range row {
			row[i] = index
		}
	}
}

// View returns a canvas onto part of this one, sharing its pixels rather
// than copying them. It keeps the parent's stride, so its rows step over the
// rest of the parent — which is what makes a sprite sheet need no special
// case: a frame of it is just a canvas.
//
// Narrow-format canvases have no view: their pixels are not Colors.
func (cv *Canvas) View(x, y, w, h int) *Canvas {
	if cv.Pixels == nil {
		return nil
	}
	x0 := clampInt(x, 0, cv.Width)
	y0 := clampInt(y, 0, cv.Height)
	x1 := clampInt(x+w, x0, cv.Width)
	y1 := clampInt(y+h, y0, cv.Height)

	width, height := x1-x0, y1-y0
	view := &Canvas{
		Width:  width,
		Height: height,
		Stride: cv.Stride,
		Clip:   Area{0, 0, width, height},
		format: ARGB32,
	}
	if width > 0 && height > 0 {
		view.Pixels = cv.Pixels[y0*cv.Stride+x0:]
	}
	return view
}
