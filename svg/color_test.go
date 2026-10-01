package svg

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// A drawing writes its colours in whichever of the ways CSS has, and every one
// of them has to come out as the colour it names.

func TestParseColorHex(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want canvas.Color
	}{
		{"#f00", canvas.RGB(255, 0, 0)},
		{"#ff0000", canvas.RGB(255, 0, 0)},
		{"#0f08", canvas.RGBA(0, 255, 0, 136)},
		{"#00ff0088", canvas.RGBA(0, 255, 0, 136)},
		{"#abc", canvas.RGB(170, 187, 204)},
		{"#FFFFFF", canvas.RGB(255, 255, 255)},
		{"#3e63dd", canvas.RGB(0x3E, 0x63, 0xDD)},
		{"3e63dd", canvas.RGB(0x3E, 0x63, 0xDD)},
		{"#3e63dd80", canvas.RGBA(0x3E, 0x63, 0xDD, 0x80)},
	} {
		got, ok := parseColor(tc.in)
		if !ok {
			t.Errorf("parseColor(%q) did not read it", tc.in)
			continue
		}
		if got != tc.want {
			t.Errorf("parseColor(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestParseColorHexRejectsWhatIsNotHex(t *testing.T) {
	for _, s := range []string{"#", "#f", "#ff", "#fffff", "#gggggg", "#12345"} {
		if c, ok := parseColor(s); ok {
			t.Errorf("parseColor(%q) = %v, want it not read at all", s, c)
		}
	}
}

func TestParseColorFunctions(t *testing.T) {
	for name, c := range map[string]canvas.Color{
		"rgb(255, 0, 0)":          canvas.RGB(255, 0, 0),
		"rgb(255 0 0)":            canvas.RGB(255, 0, 0),
		"rgba(255, 0, 0, 1)":      canvas.RGBA(255, 0, 0, 255),
		"rgba(255, 0, 0, 0.5)":    canvas.RGBA(255, 0, 0, 128),
		"rgb(255 0 0 / 50%)":      canvas.RGBA(255, 0, 0, 128),
		"rgb(100%, 0%, 0%)":       canvas.RGB(255, 0, 0),
		"rgb(300, -20, 0)":        canvas.RGB(255, 0, 0), // clamped, not wrapped
		"hsl(0, 100%, 50%)":       canvas.RGB(255, 0, 0),
		"hsl(120, 100%, 50%)":     canvas.RGB(0, 255, 0),
		"hsl(240, 100%, 50%)":     canvas.RGB(0, 0, 255),
		"hsl(0, 0%, 100%)":        canvas.RGB(255, 255, 255),
		"hsl(0, 0%, 0%)":          canvas.RGB(0, 0, 0),
		"hsl(0.5turn, 100%, 50%)": canvas.RGB(0, 255, 255),
		"hsla(0, 100%, 50%, 0.5)": canvas.RGBA(255, 0, 0, 128),
		"rgb(255,0,0)":            canvas.RGB(255, 0, 0),
		"RGB(255, 0, 0)":          canvas.RGB(255, 0, 0),
	} {
		got, ok := parseColor(name)
		if !ok {
			t.Errorf("parseColor(%q) did not read it", name)
			continue
		}
		if got != c {
			t.Errorf("parseColor(%q) = %v, want %v", name, got, c)
		}
	}
}

func TestParseColorRejectsWhatItCannotRead(t *testing.T) {
	for _, s := range []string{
		"rgb(255, 0)",
		"rgb(255, 0, 0, 0, 0)",
		"rgb(a, b, c)",
		"hsl(0, 100)",
		"hsl(0, 100, 50%)",
		"lab(50% 40 59)",
		"",
		"   ",
		"notacolor",
	} {
		if c, ok := parseColor(s); ok {
			t.Errorf("parseColor(%q) = %v, want it not read", s, c)
		}
	}
}

func TestParseColorNames(t *testing.T) {
	// A name is matched without regard to case, and a handful of them are two
	// spellings of one colour.
	for name, want := range map[string]canvas.Color{
		"red":           canvas.RGB(255, 0, 0),
		"RED":           canvas.RGB(255, 0, 0),
		"Red":           canvas.RGB(255, 0, 0),
		"blue":          canvas.RGB(0, 0, 255),
		"rebeccapurple": canvas.RGB(0x66, 0x33, 0x99),
		"gray":          canvas.RGB(128, 128, 128),
		"grey":          canvas.RGB(128, 128, 128),
		"aqua":          canvas.RGB(0, 255, 255),
		"cyan":          canvas.RGB(0, 255, 255),
		"transparent":   canvas.RGBA(0, 0, 0, 0),
	} {
		got, ok := parseColor(name)
		if !ok {
			t.Errorf("parseColor(%q) did not read it", name)
			continue
		}
		if got != want {
			t.Errorf("parseColor(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestEveryNamedColourIsARealColour(t *testing.T) {
	// The table of names is written by hand, so every one of them is checked to
	// be a colour that a drawing could have meant.
	if len(namedColors) < 140 {
		t.Errorf("%d named colours, want the whole CSS list", len(namedColors))
	}
	for name, c := range namedColors {
		if name == "transparent" {
			if c.A() != 0 {
				t.Errorf("transparent has alpha %d, want 0", c.A())
			}
			continue
		}
		if c.A() != 0xFF {
			t.Errorf("%s has alpha %d, want it opaque", name, c.A())
		}
	}
}

func TestParsePaint(t *testing.T) {
	// A fill is either a colour to paint with or nothing at all, and the two are
	// told apart by the bool rather than by the colour being zero.
	if c, ok := parsePaint("red"); !ok || c != canvas.RGB(255, 0, 0) {
		t.Errorf("parsePaint(red) = %v, %v", c, ok)
	}
	for _, s := range []string{"none", "", "url(#grad)", "  "} {
		if c, ok := parsePaint(s); ok {
			t.Errorf("parsePaint(%q) = %v, want no paint at all", s, c)
		}
	}
	// A gradient answers no to being a colour, which is what lets the style keep
	// the reference itself and follow it when the whole drawing is at hand. It
	// says so without an error, so a drawing full of them still reads.
	if _, ok := parsePaint("url(#linearGradient-1)"); ok {
		t.Error("a gradient came back as a colour")
	}
}

func TestHslToRGBAtItsCorners(t *testing.T) {
	// A hue of zero is red, a third of the way round is green and half way is
	// blue; no saturation is a grey, and no lightness is black or white.
	for _, tc := range []struct {
		name    string
		h, s, l float64
		r, g, b float64
	}{
		{"red", 0, 1, 0.5, 1, 0, 0},
		{"green", 120, 1, 0.5, 0, 1, 0},
		{"blue", 240, 1, 0.5, 0, 0, 1},
		{"a quarter of the way round is cyan", 180, 1, 0.5, 0, 1, 1},
		{"a hue past the end is the one it wraps to", 480, 1, 0.5, 0, 1, 0},
		{"white", 0, 0, 1, 1, 1, 1},
		{"black", 0, 0, 0, 0, 0, 0},
		{"grey", 0, 0, 0.5, 0.5, 0.5, 0.5},
	} {
		r, g, b := hslToRGB(tc.h, tc.s, tc.l)
		if !nearFloat(r, tc.r, 0.01) || !nearFloat(g, tc.g, 0.01) || !nearFloat(b, tc.b, 0.01) {
			t.Errorf("hsl(%g,%g,%g) for %s = %v,%v,%v, want %v,%v,%v",
				tc.h, tc.s, tc.l, tc.name, r, g, b, tc.r, tc.g, tc.b)
		}
	}
}

// nearFloat says whether two numbers are within a hair of each other, which is
// all that matters for a colour turned out of a hue.
func nearFloat(a, b, tol float64) bool { return math.Abs(a-b) <= tol }
