package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// defaultFontSize is how tall the writing is when nothing in the drawing says:
// sixteen of its own units, which is the size a `font-size` of `medium` names
// and the size a file that says nothing about it is read as.
const defaultFontSize = 16

// TextAnchor is where along the pen a piece of writing hangs: at it, centred on
// it, or ending at it. The third is what lets a label sit to the left of the
// point it names rather than to the right of it.
type TextAnchor uint8

const (
	AnchorStart TextAnchor = iota
	AnchorMiddle
	AnchorEnd
)

// TextRun is one piece of writing and everything that says where it goes: the
// pen it is written at, how far it is moved from where the piece before it left
// off, and the style it is written in. A run is what a `<tspan>` turns into, and
// the writing of a `<text>` with no `<tspan>` in it at all is a single run.
type TextRun struct {
	// Text is what is written.
	Text string
	// Style is how it is written: how tall, where along the pen it hangs, the
	// colour it is filled in, and the transform that turns the drawing's own
	// coordinates into the ones it lands on.
	Style Style
	// X and Y are where the pen goes before this piece, and HasX and HasY say
	// whether the drawing put it there: a run with no `x` of its own starts
	// where the writing around it left off, while `x="0"` is a place that has
	// been said and means the origin.
	X, Y float64
	// HasX and HasY tell an `x="0"` apart from no `x` at all, which is the
	// difference between writing at the edge of the drawing and writing on
	// from the last letter.
	HasX, HasY bool
	// DX and DY are how far the pen is moved along from where it already was,
	// which is what spaces a piece away from the one before it without saying
	// anything about where on the page either of them is. They add up through
	// the elements around a piece, the way SVG reads them.
	DX, DY float64
}

// textNode builds a `<text>`: the writing flattened into runs in the order it
// was written, rather than a node for every `<tspan>` in it, because a
// `<tspan>` is not a shape of its own — it is a stretch of the writing with a
// place and a style of its own, and it only means anything beside the writing
// on either side of it.
func (img *Image) textNode(e *element, st Style, warn func(string, ...any)) *Node {
	n := &Node{Name: e.Name, Style: st}
	n.Runs = trimEdges(img.textRuns(e, st, warn))
	var stroked, shaded bool
	for _, r := range n.Runs {
		if r.Style.HasStroke && !stroked {
			warn("the writing has a stroke, and writing is drawn by filling it, so the stroke is left out")
			stroked = true
		}
		if r.Style.fillGradient != nil && !shaded {
			warn("the writing is filled with a gradient, and writing cannot be filled with one, so it is left out")
			shaded = true
		}
	}
	// The letters are drawn one bitmap at a time from where their pen lands, so
	// a drawing that turns the writing would need each letter turned as well,
	// which this does not do. Saying so is better than writing at an angle that
	// is not the one the file asked for.
	if (st.Transform.B != 0 || st.Transform.C != 0) && len(n.Runs) > 0 {
		warn("the writing is turned, and writing goes down as it stands, so it lands where the pen is rather than along the turn")
	}
	return n
}

// textRuns is everything written inside a `<text>` or a `<tspan>`, flattened
// into one list in the order it was written. The writing and the elements
// between it are taken in turn, so a `<tspan>`'s writing goes where it was
// written — after the writing before it, not at the origin.
//
// A child that is not writing at all is left out with a warning of its own: a
// drawing may put a `<rect>` inside a `<text>` for reasons of its own, and
// there is nothing here that knows what to do with it.
func (img *Image) textRuns(e *element, st Style, warn func(string, ...any)) []TextRun {
	runs := make([]TextRun, 0, len(e.Parts))
	for _, p := range e.Parts {
		if p.Elem == nil {
			if p.Text != "" {
				runs = append(runs, TextRun{Text: p.Text, Style: st})
			}
			continue
		}
		if p.Elem.Name != "tspan" {
			warn("<%s> inside the writing is not writing this package draws, so it is left out", p.Elem.Name)
			continue
		}
		// The `<tspan>` is read the way any other element is read, over the
		// style the writing around it has: a size or a colour on it is the
		// size or colour of what it holds and nothing else.
		kid := st.with(p.Elem, warn)
		img.resolvePaints(&kid, warn)
		runs = append(runs, img.textRuns(p.Elem, kid, warn)...)
	}
	if len(runs) == 0 {
		return nil
	}
	// Where the writing starts is on the element it was written on, and it
	// lands on the first piece of that writing. A piece that a `<tspan>` closer
	// to the writing has already placed keeps the place it was given: the
	// innermost of them is the one that counts, the same way an attribute on a
	// child beats the one written on its parent.
	setCoord(&runs[0], e, warn)
	return runs
}

// setCoord puts where an element says its writing starts onto the first piece
// of it, and moves the pen along by however far it asks, leaving alone any
// place a nearer element has already said.
func setCoord(r *TextRun, e *element, warn func(string, ...any)) {
	if v, ok := coord(e, "x", warn); ok && !r.HasX {
		r.X, r.HasX = v, true
	}
	if v, ok := coord(e, "y", warn); ok && !r.HasY {
		r.Y, r.HasY = v, true
	}
	if v, ok := coord(e, "dx", warn); ok {
		r.DX += v
	}
	if v, ok := coord(e, "dy", warn); ok {
		r.DY += v
	}
}

// coord is the first of the coordinates an element wrote under a name, which
// may be a place for each letter of the writing: only the first is used, since
// the writing is put down as one piece rather than a letter at a time. A number
// that is not one leaves the writing where it was, and a list of them is said
// rather than being read as though it were the first of them by accident.
func coord(e *element, name string, warn func(string, ...any)) (float64, bool) {
	raw := strings.TrimSpace(e.attr(name))
	if raw == "" {
		return 0, false
	}
	args := splitArgs(raw)
	if len(args) == 0 {
		return 0, false
	}
	v, err := parseNumber(args[0])
	if err != nil {
		warn("the %s %q is not a number, so the writing stays where it was", name, raw)
		return 0, false
	}
	if len(args) > 1 {
		warn("the %s gives a place for each letter of the writing, and only the first is used", name)
	}
	return v, true
}

// trimEdges takes the whitespace off the two ends of the writing, which is
// what the line breaks around a `<text>` in the file are there for — they make
// the file readable and say nothing about where the writing starts and ends.
// A piece left with nothing in it is dropped, and where it would have started
// goes onto the writing after it: the place was in the file and still counts.
func trimEdges(runs []TextRun) []TextRun {
	if len(runs) == 0 {
		return runs
	}
	runs[0].Text = strings.TrimLeft(runs[0].Text, " ")
	runs[len(runs)-1].Text = strings.TrimRight(runs[len(runs)-1].Text, " ")
	if runs[0].Text == "" && len(runs) > 1 {
		head := runs[0]
		if !runs[1].HasX {
			runs[1].X, runs[1].HasX = head.X, head.HasX
		}
		if !runs[1].HasY {
			runs[1].Y, runs[1].HasY = head.Y, head.HasY
		}
		runs[1].DX += head.DX
		runs[1].DY += head.DY
	}
	out := runs[:0]
	for _, r := range runs {
		if r.Text != "" {
			out = append(out, r)
		}
	}
	return out
}

// paintText writes a `<text>`: its runs one after another with the pen carried
// from one to the next, which is what puts the writing of a `<tspan>` after the
// writing before it rather than each of them at the origin.
//
// The pen is in the drawing's own coordinates and the letters are put down as
// pixels, so each piece goes through the transform the node was given: the size
// it is drawn at is its size in those units brought onto the canvas, and where
// it starts is the pen put through the same transform. The pen then moves on by
// however wide the piece came out, so the piece after it starts where this one
// ended rather than where it began.
func paintText(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color) {
	penX, penY := 0.0, 0.0
	for _, run := range n.Runs {
		st := run.Style
		if run.HasX {
			penX = run.X
		}
		if run.HasY {
			penY = run.Y
		}
		penX, penY = penX+run.DX, penY+run.DY
		if st.Hidden || run.Text == "" {
			continue
		}
		total := m.Mul(st.Transform)
		scale := scaleOf(total)
		if scale <= 0 {
			continue
		}
		size := st.FontSize
		if size <= 0 {
			size = defaultFontSize
		}
		// The face is asked for at the size in pixels this writing covers on
		// the canvas, so a drawing painted big is written big rather than a
		// small piece of writing blown up afterwards.
		face := canvas.DefaultFace().AtSize(size * scale)
		width := float64(face.Width(run.Text)) / scale
		switch st.Anchor {
		case AnchorMiddle:
			penX -= width / 2
		case AnchorEnd:
			penX -= width
		}
		// Writing with no fill is not written at all, and a fill that is a
		// gradient rather than a colour is one this cannot write with — the
		// drawing was told so while it was read — so both come to nothing here.
		if st.HasFill {
			c := st.Fill
			if st.FillCurrent {
				c = current
			}
			if c := fade(c, st.Opacity*st.FillOpacity); c.A() > 0 {
				x, y := total.Map(penX, penY)
				face.Draw(cv, int(x+0.5), int(y+0.5), run.Text, c)
			}
		}
		penX += width
	}
}
