package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// rectShape is a rectangle, which may have its corners rounded. The corners are
// rounded by the arcs SVG describes them with, and a radius bigger than half
// the shorter side is pulled back to it, so a shape never turns inside out.
func rectShape(e *element, warn func(string, ...any)) *canvas.Path {
	x := readLength(warn, e, "x", "x")
	y := readLength(warn, e, "y", "y")
	w, okW := parseLength(e.attr("width"))
	h, okH := parseLength(e.attr("height"))
	if !okW || !okH || w <= 0 || h <= 0 {
		if e.hasAttr("width") && !okW {
			warn("the width %q is not a length, so there is no rectangle to draw", e.attr("width"))
		}
		if e.hasAttr("height") && !okH {
			warn("the height %q is not a length, so there is no rectangle to draw", e.attr("height"))
		}
		// A rectangle with no area is not a shape, which is what SVG says of
		// one with a zero or negative side.
		return nil
	}
	rx, hasRx := readLength(warn, e, "rx", "rx"), e.hasAttr("rx")
	ry, hasRy := readLength(warn, e, "ry", "ry"), e.hasAttr("ry")
	if !hasRx && !hasRy {
		p := canvas.NewPath()
		p.AddRect(x, y, w, h)
		return p
	}
	// One radius given is the radius on both corners; each is capped at half
	// the side it curves along.
	if !hasRx {
		rx = ry
	}
	if !hasRy {
		ry = rx
	}
	rx = min(rx, w/2)
	ry = min(ry, h/2)
	p := canvas.NewPath()
	// AddRoundRect pulls a radius that is bigger than half its side back to half,
	// which is the rule SVG has, and leaves a radius of no size as a plain
	// rectangle rather than as a shape with no corners at all.
	p.AddRoundRect(x, y, w, h, rx, ry)
	return p
}

// circleShape is a circle, which is an ellipse with one radius.
func circleShape(e *element, warn func(string, ...any)) *canvas.Path {
	cx, _ := parseLength(e.attr("cx"))
	cy, _ := parseLength(e.attr("cy"))
	r, ok := parseLength(e.attr("r"))
	if !ok || r <= 0 {
		return nil
	}
	p := canvas.NewPath()
	p.AddCircle(cx, cy, r)
	return p
}

// ellipseShape is an ellipse about its own centre.
func ellipseShape(e *element, warn func(string, ...any)) *canvas.Path {
	cx, _ := parseLength(e.attr("cx"))
	cy, _ := parseLength(e.attr("cy"))
	rx, okX := parseLength(e.attr("rx"))
	ry, okY := parseLength(e.attr("ry"))
	if !okX || !okY || rx <= 0 || ry <= 0 {
		return nil
	}
	p := canvas.NewPath()
	p.AddEllipse(cx, cy, rx, ry)
	return p
}

// lineShape is a straight line from one point to the other. It is drawn only by
// its stroke: SVG fills nothing, because a line encloses no area, so a line
// with a fill on it shows its fill nowhere.
func lineShape(e *element) *canvas.Path {
	x1, ok1 := parseLength(e.attr("x1"))
	y1, ok1b := parseLength(e.attr("y1"))
	x2, ok2 := parseLength(e.attr("x2"))
	y2, ok2b := parseLength(e.attr("y2"))
	if !ok1 || !ok1b || !ok2 || !ok2b {
		return nil
	}
	p := canvas.NewPath()
	p.MoveTo(x1, y1)
	p.LineTo(x2, y2)
	return p
}

// pointsShape is a polyline or a polygon: a list of points, which a polygon
// joins back to where it started and a polyline leaves open.
func pointsShape(s string, close bool, warn func(string, ...any)) *canvas.Path {
	fields := splitArgs(s)
	// The points may be written as a flat list of numbers, two to a point, or
	// as a list of pairs separated by anything else.
	nums := make([]float64, 0, len(fields))
	for _, f := range fields {
		v, err := parseNumber(f)
		if err != nil {
			warn("points needs a list of numbers, and %q is not one", f)
			return nil
		}
		nums = append(nums, v)
	}
	if len(nums)%2 != 0 {
		warn("points needs two numbers to a point, and there are %d of them", len(nums))
		return nil
	}
	if len(nums) < 4 {
		// A single point is not a line: there is no line to draw and no area to
		// fill, so there is nothing here to paint.
		return nil
	}
	pts := make([]canvas.Point, 0, len(nums)/2)
	for i := 0; i < len(nums); i += 2 {
		pts = append(pts, canvas.Point{X: nums[i], Y: nums[i+1]})
	}
	p := canvas.NewPath()
	p.AddPolyline(pts, close)
	return p
}

// readLength is the length written on one attribute, answering zero where there
// is none, and saying so where there is one that cannot be read — a number that
// is not a number is something the drawing asked for and this could not do,
// which is a warning, while an attribute that is not there at all is just a
// default.
func readLength(warn func(string, ...any), e *element, attr, what string) float64 {
	if !e.hasAttr(attr) {
		return 0
	}
	v, ok := parseLength(e.attr(attr))
	if !ok {
		warn("the %s %q is not a length, so it is taken as none", what, e.attr(attr))
	}
	return v
}

// pathData is the `d` of a `<path>`: a list of commands, each followed by the
// numbers it takes. The letters may be run together and the numbers written
// without spaces between them, as the grammar allows, so the whole string is
// scanned for both at once.
func pathData(d string, warn func(string, ...any)) *canvas.Path {
	if strings.TrimSpace(d) == "" {
		return nil
	}
	sc := &pathScanner{src: d}
	p := canvas.NewPath()
	var (
		cmd       byte
		cur       canvas.Point
		start     canvas.Point // where the current subpath began, for Z
		haveCtl   canvas.Point // the control point of a curve, for S and T
		prevCubic bool
		prevQuad  bool
		started   bool
	)
	for {
		// Everything read is the whole of it: a command repeated by a number
		// runs out of numbers as soon as the string does.
		if sc.eof() {
			break
		}
		// A letter is a command; after one, its own letter repeats when the
		// next thing in the string is a number.
		if sc.atLetter() {
			cmd = sc.letter()
		} else if cmd == 0 {
			warn("the path data starts with something that is not a command")
			return nil
		} else if cmd == 'M' {
			// A number after a moveto is another moveto, and after anything
			// else it is a lineto, which is what SVG calls an implicit one.
			cmd = 'L'
		} else if cmd == 'm' {
			cmd = 'l'
		} else if cmd == 'Z' || cmd == 'z' {
			cmd = 0
			sc.skipSpace()
			if sc.eof() {
				break
			}
			continue
		}
		rel := cmd >= 'a' && cmd <= 'z'
		upper := upperCmd(cmd)

		switch upper {
		case 'M':
			x, ok1 := sc.number()
			y, ok2 := sc.number()
			if !ok1 || !ok2 {
				warn("the %c command needs two numbers", cmd)
				return nil
			}
			if rel && started {
				x, y = cur.X+x, cur.Y+y
			}
			cur = canvas.Point{X: x, Y: y}
			start = cur
			p.MoveTo(x, y)
			started = true
			haveCtl, prevCubic, prevQuad = cur, false, false
			// A pair of numbers written after a move is a line, which is what
			// SVG says of one, so the command for what follows is a line and
			// not another move.
			if rel {
				cmd = 'l'
			} else {
				cmd = 'L'
			}

		case 'L':
			x, ok1 := sc.number()
			y, ok2 := sc.number()
			if !ok1 || !ok2 {
				warn("the %c command needs two numbers", cmd)
				return nil
			}
			if rel {
				x, y = cur.X+x, cur.Y+y
			}
			cur = canvas.Point{X: x, Y: y}
			p.LineTo(x, y)
			haveCtl, prevCubic, prevQuad = cur, false, false

		case 'H':
			x, ok := sc.number()
			if !ok {
				warn("the %c command needs one number", cmd)
				return nil
			}
			if rel {
				x = cur.X + x
			}
			cur = canvas.Point{X: x, Y: cur.Y}
			p.LineTo(x, cur.Y)
			haveCtl, prevCubic, prevQuad = cur, false, false

		case 'V':
			y, ok := sc.number()
			if !ok {
				warn("the %c command needs one number", cmd)
				return nil
			}
			if rel {
				y = cur.Y + y
			}
			cur = canvas.Point{X: cur.X, Y: y}
			p.LineTo(cur.X, y)
			haveCtl, prevCubic, prevQuad = cur, false, false

		case 'C':
			c1, ok1 := sc.point()
			c2, ok2 := sc.point()
			e, ok3 := sc.point()
			if !ok1 || !ok2 || !ok3 {
				warn("the %c command needs six numbers", cmd)
				return nil
			}
			if rel {
				c1, c2, e = add(cur, c1), add(cur, c2), add(cur, e)
			}
			p.CubicTo(c1.X, c1.Y, c2.X, c2.Y, e.X, e.Y)
			cur, haveCtl, prevCubic, prevQuad = e, c2, true, false

		case 'S':
			c2, ok1 := sc.point()
			e, ok2 := sc.point()
			if !ok1 || !ok2 {
				warn("the %c command needs four numbers", cmd)
				return nil
			}
			if rel {
				c2, e = add(cur, c2), add(cur, e)
			}
			// The first control is the mirror of the last one, and the first
			// control point itself when the curve before was not a curve.
			c1 := cur
			if prevCubic {
				c1 = mirror(haveCtl, cur)
			}
			p.CubicTo(c1.X, c1.Y, c2.X, c2.Y, e.X, e.Y)
			cur, haveCtl, prevCubic, prevQuad = e, c2, true, false

		case 'Q':
			c, ok1 := sc.point()
			e, ok2 := sc.point()
			if !ok1 || !ok2 {
				warn("the %c command needs four numbers", cmd)
				return nil
			}
			if rel {
				c, e = add(cur, c), add(cur, e)
			}
			p.QuadTo(c.X, c.Y, e.X, e.Y)
			cur, haveCtl, prevCubic, prevQuad = e, c, false, true

		case 'T':
			e, ok := sc.point()
			if !ok {
				warn("the %c command needs two numbers", cmd)
				return nil
			}
			if rel {
				e = add(cur, e)
			}
			c := cur
			if prevQuad {
				c = mirror(haveCtl, cur)
			}
			p.QuadTo(c.X, c.Y, e.X, e.Y)
			cur, haveCtl, prevCubic, prevQuad = e, c, false, true

		case 'A':
			rx, ok1 := sc.number()
			ry, ok2 := sc.number()
			rot, ok3 := sc.number()
			large, ok4 := sc.flag()
			sweep, ok5 := sc.flag()
			e, ok6 := sc.point()
			if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
				warn("the %c command needs seven numbers", cmd)
				return nil
			}
			if rel {
				e = add(cur, e)
			}
			p.ArcTo(rx, ry, rot, large, sweep, e.X, e.Y)
			cur, haveCtl, prevCubic, prevQuad = e, cur, false, false

		case 'Z':
			p.Close()
			cur = start
			started = true
			haveCtl, prevCubic, prevQuad = cur, false, false
			cmd = 0

		default:
			warn("%q is not a command a path has", string(cmd))
			return nil
		}
	}
	return p
}

// upperCmd is the command without its case, since the case only says whether
// the numbers are measured from where the line is.
func upperCmd(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 'a' + 'A'
	}
	return c
}

// pathScanner reads the path data one thing at a time, skipping the punctuation
// between them, which is what lets the numbers be written with or without
// commas and spaces between them.
type pathScanner struct {
	src string
	pos int
}

// skipSpace passes over the spaces, commas and tabs between two things, and any
// minus or plus signs that are written as their own token.
func (s *pathScanner) skipSpace() {
	for s.pos < len(s.src) {
		switch s.src[s.pos] {
		case ' ', ',', '\t', '\n', '\r':
			s.pos++
		default:
			return
		}
	}
}

// eof reports whether everything has been read.
func (s *pathScanner) eof() bool {
	s.skipSpace()
	return s.pos >= len(s.src)
}

// atLetter reports whether the next thing is a command rather than a number.
func (s *pathScanner) atLetter() bool {
	s.skipSpace()
	if s.pos >= len(s.src) {
		return false
	}
	c := s.src[s.pos]
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// letter reads one command letter.
func (s *pathScanner) letter() byte {
	c := s.src[s.pos]
	s.pos++
	return c
}

// number reads one number, which may be written with an exponent as SVG allows.
func (s *pathScanner) number() (float64, bool) {
	s.skipSpace()
	start := s.pos
	if s.pos < len(s.src) && (s.src[s.pos] == '+' || s.src[s.pos] == '-') {
		s.pos++
	}
	seen := false
	for s.pos < len(s.src) && s.src[s.pos] >= '0' && s.src[s.pos] <= '9' {
		s.pos++
		seen = true
	}
	if s.pos < len(s.src) && s.src[s.pos] == '.' {
		s.pos++
		for s.pos < len(s.src) && s.src[s.pos] >= '0' && s.src[s.pos] <= '9' {
			s.pos++
			seen = true
		}
	}
	if seen && s.pos < len(s.src) && (s.src[s.pos] == 'e' || s.src[s.pos] == 'E') {
		mark := s.pos
		s.pos++
		if s.pos < len(s.src) && (s.src[s.pos] == '+' || s.src[s.pos] == '-') {
			s.pos++
		}
		digits := false
		for s.pos < len(s.src) && s.src[s.pos] >= '0' && s.src[s.pos] <= '9' {
			s.pos++
			digits = true
		}
		if !digits {
			s.pos = mark
		}
	}
	if !seen {
		s.pos = start
		return 0, false
	}
	v, err := parseNumber(s.src[start:s.pos])
	if err != nil {
		return 0, false
	}
	return v, true
}

// flag reads the zero or one of the arc command, which SVG allows to be written
// without any separator at all, as in `a1 1 0 011 1`.
func (s *pathScanner) flag() (bool, bool) {
	s.skipSpace()
	if s.pos >= len(s.src) {
		return false, false
	}
	switch s.src[s.pos] {
	case '0':
		s.pos++
		return false, true
	case '1':
		s.pos++
		return true, true
	}
	return false, false
}

// point reads the two numbers of a point.
func (s *pathScanner) point() (canvas.Point, bool) {
	x, ok1 := s.number()
	y, ok2 := s.number()
	return canvas.Point{X: x, Y: y}, ok1 && ok2
}
