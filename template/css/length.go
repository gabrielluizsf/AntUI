package css

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// unit is the measuring stick behind a Length.
type unit uint8

const (
	unitPx unit = iota
	unitPct
	unitAuto
	unitNone
	unitEm
	unitRem
	unitVw
	unitVh
	unitVmin
	unitVmax
	unitCh
	unitEx
	unitCm
	unitMm
	unitIn
	unitPt
	unitPc
	unitQ
)

// The fixed physical units resolve against the CSS reference density of 96
// pixels per inch: 1in = 96px, 1cm = 96/2.54px and so on. These are reference
// pixels, scaled with the window like every other fixed length at draw time.
const (
	pxPerIn = 96.0
	pxPerCm = pxPerIn / 2.54
	pxPerMm = pxPerIn / 25.4
	pxPerQ  = pxPerMm / 4
	pxPerPt = pxPerIn / 72
	pxPerPc = 16.0
)

// DefaultFontSize is the reference-pixel font the engine assumes when nothing
// sets font-size; rem and the font-relative units use it.
const DefaultFontSize = 16

// Units is the measurement context a Length needs to become pixels: the
// containing window (vw/vh/vmin/vmax), the containing width (percentages) and
// the font sizes (em/ch/ex, rem). Zero fields fall back to sensible defaults.
type Units struct {
	Width, Height int // the window, for viewport units and %-of-width
	Font          int // the element's font size, for em/ch/ex
	Root          int // the document font size, for rem
}

// base is the fallback font when a context has not said anything.
func (u Units) font() int { return max(u.Font, DefaultFontSize) }

func (u Units) root() int { return max(u.Root, DefaultFontSize) }

// A Length is a size a template can draw with: fixed reference pixels, a
// percentage of the surrounding measure, a length in any CSS unit, or a
// keyword such as auto.
type Length struct {
	u     unit
	value float64
}

// Fixed builds a pixel length.
func Fixed(px float64) Length { return Length{u: unitPx, value: px} }

// Pct builds a percentage length.
func Pct(p float64) Length { return Length{u: unitPct, value: p} }

// Auto is the "auto" keyword, used by width, height and margins.
func Auto() Length { return Length{u: unitAuto} }

// Zero is a zero pixel length.
func Zero() Length { return Fixed(0) }

// IsPct reports whether the length is a percentage.
func (l Length) IsPct() bool { return l.u == unitPct }

// Auto reports whether the length is the auto keyword.
func (l Length) Auto() bool { return l.u == unitAuto }

// None reports whether the length represents the none keyword.
func (l Length) None() bool { return l.u == unitNone }

// Unit returns the length's unit name as CSS writes it, or "" for auto/none.
func (l Length) Unit() string { return unitName[l.u] }

// unitName maps every unit to the CSS suffix that spells it.
var unitName = map[unit]string{
	unitPx:   "px",
	unitPct:  "%",
	unitEm:   "em",
	unitRem:  "rem",
	unitVw:   "vw",
	unitVh:   "vh",
	unitVmin: "vmin",
	unitVmax: "vmax",
	unitCh:   "ch",
	unitEx:   "ex",
	unitCm:   "cm",
	unitMm:   "mm",
	unitIn:   "in",
	unitPt:   "pt",
	unitPc:   "pc",
	unitQ:    "q",
	unitAuto: "",
	unitNone: "",
}

// ref converts the length into reference pixels under a context, the unit
// arithmetic every resolution is built on. Percentages are of the containing
// width, viewport units of the window, font units of the font sizes; fixed and
// physical lengths are their own number of pixels.
func (l Length) ref(ctx Units) float64 {
	switch l.u {
	case unitPx:
		return l.value
	case unitPct:
		return l.value * float64(ctx.Width) / 100
	case unitEm:
		return l.value * float64(ctx.font())
	case unitRem:
		return l.value * float64(ctx.root())
	case unitVw:
		return l.value * float64(ctx.Width) / 100
	case unitVh:
		return l.value * float64(ctx.Height) / 100
	case unitVmin:
		return l.value * float64(min(ctx.Width, ctx.Height)) / 100
	case unitVmax:
		return l.value * float64(max(ctx.Width, ctx.Height)) / 100
	// The advance of "0" and the x-height both settle around half an em for
	// the built-in face; the engine reads both as 0.5em.
	case unitCh:
		return l.value * 0.5 * float64(ctx.font())
	case unitEx:
		return l.value * 0.5 * float64(ctx.font())
	case unitCm:
		return l.value * pxPerCm
	case unitMm:
		return l.value * pxPerMm
	case unitIn:
		return l.value * pxPerIn
	case unitPt:
		return l.value * pxPerPt
	case unitPc:
		return l.value * pxPerPc
	case unitQ:
		return l.value * pxPerQ
	}
	return 0
}

// Resolve turns the length into reference pixels under a context. Keywords
// (auto, none) resolve to zero; a caller that needs to tell them apart asks
// [Length.Auto] or [Length.None] first.
func (l Length) Resolve(ctx Units) int {
	return int(math.Round(l.ref(ctx)))
}

// Px resolves the length against a base: percentages scale the base, fixed
// lengths return their pixels, keywords return zero. It is [Length.Resolve]
// with the base standing in for the whole measurement context, which keeps
// the old one-argument call sites working.
func (l Length) Px(base int) int {
	return l.Resolve(Units{Width: base, Height: base, Font: base, Root: base})
}

// parseLengthAt reads a length with a measurement context, so % and the
// viewport/font units resolve and calc() min() max() clamp() can evaluate.
// A calc() that balances is converted out of its formula into plain px.
func parseLengthAt(raw string, ctx Units) (Length, error) {
	s := strings.TrimSpace(strings.ToLower(raw))
	if s == "" {
		return Length{}, fmt.Errorf("css: empty length")
	}
	switch s {
	case "auto":
		return Auto(), nil
	case "none":
		return Length{u: unitNone}, nil
	case "inherit", "initial", "unset":
		return Length{u: unitNone}, nil
	}
	if mathLook(s) {
		if px, ok := EvalMath(s, ctx); ok {
			return Fixed(px), nil
		}
		return Length{}, fmt.Errorf("css: %q does not evaluate", raw)
	}
	return parseLength(s)
}

// parseLength reads "12", "12px", "50%", "1.5rem", "3vw", "2in" or "auto".
func parseLength(raw string) (Length, error) {
	s := strings.TrimSpace(strings.ToLower(raw))
	if s == "" {
		return Length{}, fmt.Errorf("css: empty length")
	}
	if s == "auto" {
		return Auto(), nil
	}
	if s == "none" {
		return Length{u: unitNone}, nil
	}
	if s == "inherit" || s == "initial" || s == "unset" {
		return Length{u: unitNone}, nil
	}
	if strings.HasSuffix(s, "%") {
		return parseLenPrefix(s, unitPct, "%")
	}
	for _, suf := range unitSuffixes {
		if strings.HasSuffix(s, suf.name) {
			return parseLenPrefix(s, suf.u, suf.name)
		}
	}
	return parseLenPrefix(s, unitPx, "")
}

// unitSuffixes lists every unit suffix the parser accepts, longest first so a
// suffix is not swallowed by a shorter one: "rem" must win over "em", and
// "vmax" over "vh". "px" is a suffix like any other here.
var unitSuffixes = []struct {
	name string
	u    unit
}{
	{"vmin", unitVmin},
	{"vmax", unitVmax},
	{"rem", unitRem},
	{"cm", unitCm},
	{"mm", unitMm},
	{"in", unitIn},
	{"pt", unitPt},
	{"pc", unitPc},
	{"px", unitPx},
	{"em", unitEm},
	{"vw", unitVw},
	{"vh", unitVh},
	{"ch", unitCh},
	{"ex", unitEx},
	{"q", unitQ},
}

func parseLenPrefix(s string, u unit, suffix string) (Length, error) {
	body := strings.TrimSuffix(s, suffix)
	v, err := strconv.ParseFloat(strings.TrimSpace(body), 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return Length{}, fmt.Errorf("css: %q is not a length", s)
	}
	return Length{u: u, value: v}, nil
}

// parseFour reads up to four space-separated lengths the way CSS expands them
// across sides: 1 value for all four, 2 for vertical/horizontal, 3 for
// top/sides/bottom, and 4 for top/right/bottom/left.
func parseFour(raw string) ([4]Length, error) {
	var out [4]Length
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return out, fmt.Errorf("css: expected a length, got %q", raw)
	}
	vals := make([]Length, len(parts))
	for i, p := range parts {
		v, err := parseLength(p)
		if err != nil {
			return out, err
		}
		vals[i] = v
	}
	return expandInto(out, vals), nil
}

// parseFourAt is parseFour under a measurement context, so % and the
// viewport/font units resolve and box shorthands may hold a calc(). A box
// shorthand that is a single math formula resolves to one equal length.
func parseFourAt(raw string, ctx Units) ([4]Length, error) {
	if mathLook(raw) {
		px, ok := EvalMath(raw, ctx)
		if ok {
			return [4]Length{Fixed(px), Fixed(px), Fixed(px), Fixed(px)}, nil
		}
		return [4]Length{}, fmt.Errorf("css: %q does not evaluate", raw)
	}
	var out [4]Length
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return out, fmt.Errorf("css: expected a length, got %q", raw)
	}
	vals := make([]Length, len(parts))
	for i, p := range parts {
		v, err := parseLengthAt(p, ctx)
		if err != nil {
			return out, err
		}
		vals[i] = v
	}
	return expandInto(out, vals), nil
}

func expandInto(out [4]Length, vals []Length) [4]Length {
	switch len(vals) {
	case 1:
		return [4]Length{vals[0], vals[0], vals[0], vals[0]}
	case 2:
		return [4]Length{vals[0], vals[1], vals[0], vals[1]}
	case 3:
		return [4]Length{vals[0], vals[1], vals[2], vals[1]}
	default:
		return [4]Length{vals[0], vals[1], vals[2], vals[3]}
	}
}
