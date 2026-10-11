package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

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
