package canvas

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
	Dirty  Area // the bbox of every write since the last ResetDirty

	// blurFront and blurBack are the two float buffers a blur passes its
	// picture between. They are the largest thing a canvas ever allocates
	// besides its pixels, and a window keeps them for as long as it lives
	// rather than asking for new ones on every frame that blurs something.
	blurFront, blurBack []float64

	format  Format
	narrow  []byte  // the 16- or 8-bit pixels
	palette []Color // 256 entries, on Pal8
	nearest []byte  // 32K: an RGB555 colour to its palette index
	dither  bool

	// base is the frame that was last presented, which every write is
	// checked against; spans, rows and change record what really moved as it
	// happens. Nil until Compare is called, and then a window is the only
	// thing that has any use for it. See change.go.
	base   []Color
	spans  []rowSpan
	rows   []int
	change Area

	// Premultiplied records whether each channel has already been multiplied
	// by its own alpha. The flag lives beside the pixels rather than beside the
	// image so that converting one twice — which darkens it by the alpha
	// again and looks exactly like the first conversion — cannot happen.
	Premultiplied bool

	// over is whether what is drawn into this canvas goes *over* what is
	// already there, keeping the alpha of both, rather than over a surface
	// that is assumed to be solid. See NewLayer.
	over bool
}
