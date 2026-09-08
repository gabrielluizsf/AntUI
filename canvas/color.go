package canvas

import (
	"fmt"
	"strconv"
	"strings"
)

// Color is a pixel in 0xAARRGGBB order: alpha in the top eight bits, then
// red, green and blue. It is the one colour representation the whole library
// speaks, on every surface format — the narrow formats convert on the way in
// and on the way out, so callers never see a 565 word or a palette index.
type Color uint32

// RGB builds an opaque colour.
func RGB(r, g, b uint8) Color {
	return Color(0xFF000000 | uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

// RGBA builds a colour with an explicit alpha.
func RGBA(r, g, b, a uint8) Color {
	return Color(uint32(a)<<24 | uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

// A returns the alpha channel.
func (c Color) A() uint8 { return uint8(c >> 24) }

// R returns the red channel.
func (c Color) R() uint8 { return uint8(c >> 16) }

// G returns the green channel.
func (c Color) G() uint8 { return uint8(c >> 8) }

// B returns the blue channel.
func (c Color) B() uint8 { return uint8(c) }

// RGBA implements image/color.Color, so a canvas.Color can be handed to
// anything in the standard library that draws. The channels are
// alpha-premultiplied and scaled to 16 bits, as that interface requires.
func (c Color) RGBA() (r, g, b, a uint32) {
	al := uint32(c.A())
	al |= al << 8
	r = uint32(c.R()) * al / 0xFF
	g = uint32(c.G()) * al / 0xFF
	b = uint32(c.B()) * al / 0xFF
	return r, g, b, al
}

// The built-in palette.
const (
	Transparent Color = 0x00000000
	Black       Color = 0xFF000000
	White       Color = 0xFFFFFFFF
	Red         Color = 0xFFE5484D
	Green       Color = 0xFF30A46C
	Blue        Color = 0xFF3E63DD
	Yellow      Color = 0xFFFFC53D
	Orange      Color = 0xFFF76B15
	Purple      Color = 0xFF8E4EC6
	Cyan        Color = 0xFF00A2C7
	Magenta     Color = 0xFFD6409F
	Gray        Color = 0xFF8B8D98
	LightGray   Color = 0xFFD9D9E0
	DarkGray    Color = 0xFF3B3C42
)

// Blend composites src over dst, honouring the alpha of src. The result is
// always opaque: this is what lands in a framebuffer, and a framebuffer has
// nothing behind it to show through.
func Blend(dst, src Color) Color {
	a := uint32(src>>24) & 0xFF
	if a == 255 {
		return src | 0xFF000000
	}
	if a == 0 {
		return dst
	}
	inv := 255 - a
	r := ((uint32(src>>16)&0xFF)*a + (uint32(dst>>16)&0xFF)*inv + 127) / 255
	g := ((uint32(src>>8)&0xFF)*a + (uint32(dst>>8)&0xFF)*inv + 127) / 255
	b := ((uint32(src)&0xFF)*a + (uint32(dst)&0xFF)*inv + 127) / 255
	return Color(0xFF000000 | r<<16 | g<<8 | b)
}

// Mix interpolates between two colours, t running from 0 (all a) to 1 (all
// b). Every channel is interpolated, alpha included.
func Mix(a, b Color, t float32) Color {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	k := uint32(t*255 + 0.5)
	inv := 255 - k
	ch := func(shift uint) uint32 {
		return ((uint32(a>>shift)&0xFF)*inv + (uint32(b>>shift)&0xFF)*k) / 255
	}
	return Color(ch(24)<<24 | ch(16)<<16 | ch(8)<<8 | ch(0))
}

// Shade lightens the colour for a positive amount and darkens it for a
// negative one, from -1 to 1. The alpha is left alone.
func Shade(c Color, amount float32) Color {
	var target Color
	if amount >= 0 {
		target = (c & 0xFF000000) | 0x00FFFFFF
	} else {
		target = c & 0xFF000000
		amount = -amount
	}
	return Mix(c, target, amount)
}

// Fade replaces the colour's alpha, keeping its RGB.
func Fade(c Color, alpha int) Color {
	return (c & 0x00FFFFFF) | Color(uint32(clampInt(alpha, 0, 255))<<24)
}

// String writes a colour the way it is written everywhere else: "#RRGGBB",
// or "#AARRGGBB" when it is not opaque.
func (c Color) String() string {
	if c.A() == 0xFF {
		return fmt.Sprintf("#%02X%02X%02X", c.R(), c.G(), c.B())
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", c.A(), c.R(), c.G(), c.B())
}

// ParseColor reads "#RRGGBB" or "#AARRGGBB". The hash is optional.
func ParseColor(s string) (Color, error) {
	body := strings.TrimPrefix(s, "#")
	if len(body) != 6 && len(body) != 8 {
		return 0, fmt.Errorf("antui: %q is not a colour; it looks like \"#3E63DD\"", s)
	}
	v, err := strconv.ParseUint(body, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("antui: %q is not a colour; it looks like \"#3E63DD\"", s)
	}
	if len(body) == 6 {
		v |= 0xFF000000
	}
	return Color(v), nil
}

// MarshalText and UnmarshalText let a colour be written down — in JSON, in a
// config file, in anything that goes through encoding.TextMarshaler — as
// "#3E63DD" rather than as the number 4283372509.
func (c Color) MarshalText() ([]byte, error) { return []byte(c.String()), nil }

func (c *Color) UnmarshalText(text []byte) error {
	v, err := ParseColor(string(text))
	if err != nil {
		return err
	}
	*c = v
	return nil
}
