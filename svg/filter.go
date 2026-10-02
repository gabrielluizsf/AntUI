package svg

import (
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// A `filter` says what an element's picture goes through on its way into the
// drawing: a list of functions — `blur(1) grayscale(1) hue-rotate(90deg)` —
// each of which changes the picture of everything to its left, all of them run
// in the order they were written. The picture is the whole of what the element
// draws, its children and its writing included, which is why a filter is not
// inherited: a `<g filter>` turns the one picture it makes of its children
// once, rather than each child being turned as it is drawn, and the same is
// true of a clip and a mask beside it.
//
// The functions here are the ones the canvas already knows how to apply —
// grayscale, sepia, invert, brightness, contrast, hue-rotate and blur, the
// same set a stylesheet's `filter` maps onto — and everything else, a
// `drop-shadow` or a reference to a `<filter>` element (which is not read in
// this package), is said out loud and left out of the list rather than taken
// to mean nothing: what the list can still do, it still does.
//
// The list reaches the canvas in [applyFilters], which runs each function over
// the picture the element made — the whole of the picture, since that is what
// the element is asking to have put through it, and a blur that reaches past
// where the element stopped has room to reach into.

// filterOp is one function of the list with what it was given: a fraction for
// grayscale, sepia, invert, brightness and contrast, degrees for hue-rotate,
// and a length in the drawing's own units for blur — the one of them that is a
// size, and so the one that is turned into pixels when it is applied.
type filterOp struct {
	kind   canvas.FilterKind
	amount float64
}

// readFilters reads the list a `filter` attribute holds: `none` is no filter
// at all and says so without complaint, and every function is read on its own,
// so a list this package cannot wholly read still does the part of it that it
// can — one function left out is a warning, not a whole filter dropped.
func readFilters(raw string, warn func(string, ...any)) []filterOp {
	v := strings.TrimSpace(raw)
	if v == "" {
		warn("the filter %q is not one this package can read, so it is left out", raw)
		return nil
	}
	if strings.EqualFold(v, "none") {
		return nil
	}
	var out []filterOp
	for _, tok := range filterTokens(v) {
		if isFilterRef(tok) {
			// A reference to a `<filter>` element elsewhere in the drawing is
			// the other way a file asks for one of these, and what is inside a
			// `<filter>` is not read here at all — so the reference is named
			// for what it is and left out, and the rest of the list goes on.
			warn("the filter reference %q is not read here, so it is left out", tok)
			continue
		}
		if op, ok := readFilter(tok); ok {
			out = append(out, op)
			continue
		}
		warn("the filter %q is not one this package can read, so it is left out", tok)
	}
	return out
}

// filterTokens breaks a filter list into its functions, keeping the arguments
// of each one with it: `drop-shadow(0 1px 2px black)` is one thing to try to
// read and not five, because a space inside the parentheses is part of what it
// says. Only a space at the top level, where nothing is open, comes between
// two functions.
func filterTokens(s string) []string {
	var out []string
	depth, start := 0, 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '(':
			depth++
		case c == ')':
			if depth > 0 {
				depth--
			}
		case depth == 0 && (c == ' ' || c == '\t' || c == '\n' || c == '\r'):
			if start < i {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

// isFilterRef says the token is a `url(...)` — the reference to a `<filter>`
// element that this package reads the name of and nothing else.
func isFilterRef(tok string) bool {
	t := strings.TrimSpace(tok)
	return strings.HasPrefix(strings.ToLower(t), "url(") ||
		strings.EqualFold(t, "url")
}

// readFilter reads one function of the list, such as `blur(2px)` or
// `grayscale(50%)`. What is not one of the functions this package can apply,
// does not open and close its arguments, or gives a number that is not one
// answers false, and the caller says so rather than guessing what was meant.
func readFilter(tok string) (filterOp, bool) {
	open := strings.IndexByte(tok, '(')
	if open < 0 || !strings.HasSuffix(tok, ")") {
		return filterOp{}, false
	}
	name := strings.ToLower(strings.TrimSpace(tok[:open]))
	args := strings.TrimSpace(tok[open+1 : len(tok)-1])
	switch name {
	case "grayscale", "sepia", "invert", "brightness", "contrast":
		v, ok := filterAmount(args)
		if !ok {
			return filterOp{}, false
		}
		kind := canvas.FilterGrayscale
		switch name {
		case "sepia":
			kind = canvas.FilterSepia
		case "invert":
			kind = canvas.FilterInvert
		case "brightness":
			kind = canvas.FilterBrightness
		case "contrast":
			kind = canvas.FilterContrast
		}
		return filterOp{kind: kind, amount: v}, true
	case "hue-rotate":
		deg, err := parseAngle(args)
		if err != nil {
			return filterOp{}, false
		}
		return filterOp{kind: canvas.FilterHueRotate, amount: deg}, true
	case "blur":
		if args == "" {
			// A blur of nothing is no blur, which is what the function with no
			// argument says in CSS, and it is read rather than refused so that
			// a list containing it is not warning about something harmless.
			return filterOp{kind: canvas.FilterBlur, amount: 0}, true
		}
		v, ok := parseLength(args)
		if !ok || v < 0 {
			return filterOp{}, false
		}
		return filterOp{kind: canvas.FilterBlur, amount: v}, true
	}
	return filterOp{}, false
}

// filterAmount reads one function's argument: a number or a percentage, where
// an empty argument is the identity of one. Something else is not a number at
// all, and the caller says so instead of the filter quietly becoming full
// strength or none.
func filterAmount(args string) (float64, bool) {
	if args == "" {
		return 1, true
	}
	if p, ok := strings.CutSuffix(args, "%"); ok {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return 0, false
		}
		return v / 100, true
	}
	v, err := strconv.ParseFloat(args, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// applyFilters runs an element's filter list over the picture it drew, one
// function after another in the order they were written: a grayscale and then
// a blur is not the same picture as a blur and then a grayscale, and the list
// says which one it wants.
//
// Every function goes over the whole picture the layer holds, which is exactly
// what the element asked to have put through it. A smaller rectangle would
// have to be what the element painted, widened before a blur by as far as three
// box passes can carry a picture past where the element stopped — and the
// canvas only remembers where it was written when it is watching for changes,
// which this one is not, so the painted rectangle is not something to be had
// here. The whole of the layer is honest about that, gives a blur room to
// reach in every direction, and costs a walk over transparent pixels that come
// out of a colour filter untouched and of a blur still carrying nothing.
//
// A blur is a size and is written in the drawing's own units, so it is turned
// into pixels here the same as a stroke's width is: one unit of the drawing is
// scale pixels of the canvas, and the scale has the element's transform in it
// because what a `<g transform="scale(2)">` draws is twice as big and its blur
// has to be too. The rest of the functions are fractions and angles, which a
// scale does not touch.
func applyFilters(cv *canvas.Canvas, filters []filterOp, scale float64) {
	for _, f := range filters {
		amount := f.amount
		if f.kind == canvas.FilterBlur {
			amount = f.amount * scale
		}
		cv.FilterRegion(0, 0, cv.Width, cv.Height, f.kind, amount)
	}
}
