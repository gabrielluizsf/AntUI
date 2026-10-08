package svg

import (
	"math"
	"strconv"
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
	// XList and YList are a place for each letter rather than one place for
	// the whole of the piece, which is what `x="10 14 19"` says: the first
	// entry is where the first letter goes, the second where the second goes,
	// and a letter the list ran out before is left to follow the one before it.
	// They are only there when the drawing named more than one place — one
	// number is the X and Y above, which is what carries a single place across
	// the whitespace the file wrapped the writing with.
	XList, YList []float64
	// DX and DY are how far the pen is moved along from where it already was,
	// which is what spaces a piece away from the one before it without saying
	// anything about where on the page either of them is. They add up through
	// the elements around a piece, the way SVG reads them.
	DX, DY float64
	// Path is the outline a `<textPath>` writes along rather than a pen the
	// piece hangs on, and Offset is how far along it the writing starts. Both
	// are nil and zero for writing that has a pen instead, which is every
	// piece of writing a drawing has ever had.
	Path   *canvas.Path
	Offset float64
}

// textNode builds a `<text>`: the writing flattened into runs in the order it
// was written, rather than a node for every `<tspan>` in it, because a
// `<tspan>` is not a shape of its own — it is a stretch of the writing with a
// place and a style of its own, and it only means anything beside the writing
// on either side of it.
func (img *Image) textNode(e *element, st Style, warn func(string, ...any)) *Node {
	n := &Node{Name: e.Name, Style: st}
	n.Runs = trimEdges(img.textRuns(e, st, warn), st.PreserveSpace)
	// The letters are drawn one bitmap at a time from where their pen lands, so
	// a drawing that turns the writing has each letter turned with it at paint
	// time. A turn with no way back is the one that cannot be put right: the
	// matrix would have to be undone to know where a letter's pixels land, and
	// there is no such place. Saying so is better than writing somewhere the
	// file did not ask for.
	if (st.Transform.B != 0 || st.Transform.C != 0) && len(n.Runs) > 0 {
		if _, ok := st.Transform.Inverse(); !ok {
			warn("the writing is turned by a transform with no way back, so it lands where the pen is rather than along the turn")
		}
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
		if p.Elem.Name == "textPath" {
			// The writing is not hung on a pen at all: it is laid along the
			// outline the href names, with `startOffset` saying how far along
			// it the first letter goes. What it names could not be followed,
			// and then there is nowhere to write, which is worth one sentence
			// rather than the writing turning up in the wrong place.
			along := img.textPathOutline(p.Elem, warn)
			if along == nil {
				continue
			}
			kid := st.with(p.Elem, warn)
			img.resolvePaints(&kid, warn)
			sub := img.textRuns(p.Elem, kid, warn)
			off := startOffset(p.Elem, along, warn)
			for i := range sub {
				sub[i].Path, sub[i].Offset = along, off
			}
			runs = append(runs, sub...)
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
	// lands on the pieces of that writing. A piece that a `<tspan>` closer
	// to the writing has already placed keeps the place it was given: the
	// innermost of them is the one that counts, the same way an attribute on
	// a child beats the one written on its parent.
	place(runs, e, warn)
	return runs
}

// place puts where an element says its writing starts onto the pieces of that
// writing. `x` and `y` may name a place for every letter rather than one place
// for the whole of it, and both are read here: the list runs through the
// letters of the pieces in the order they were written, taking one value for
// each letter even where the value does not land — a piece a nearer element
// has already placed keeps that place, but its letters still count against the
// list, so the value after them goes to the letter after them. `dx` and `dy`
// are a move of the pen rather than a place for it, and are never one per
// letter.
func place(runs []TextRun, e *element, warn func(string, ...any)) {
	placeAt(runs, e, "x", warn)
	placeAt(runs, e, "y", warn)
	if v, ok := coord(e, "dx", warn); ok {
		runs[0].DX += v
	}
	if v, ok := coord(e, "dy", warn); ok {
		runs[0].DY += v
	}
}

// placeAt reads one of the two names a place may be written under, either as
// one number for the whole of the writing or as one number for each letter of
// it.
func placeAt(runs []TextRun, e *element, name string, warn func(string, ...any)) {
	raw := strings.TrimSpace(e.attr(name))
	if raw == "" {
		return
	}
	args := splitArgs(raw)
	if len(args) == 0 {
		return
	}
	if len(args) == 1 {
		v, err := parseNumber(args[0])
		if err != nil {
			warn("the %s %q is not a number, so the writing stays where it was", name, raw)
			return
		}
		at := &runs[0]
		if name == "y" {
			if !at.HasY {
				at.Y, at.HasY = v, true
			}
			return
		}
		if !at.HasX {
			at.X, at.HasX = v, true
		}
		return
	}
	// One number for each letter: every one of them has to be a number, since
	// a list with a hole in it would put a letter in a place the file did not
	// name, and the letters after it would follow it there.
	vals := make([]float64, len(args))
	for i, a := range args {
		v, err := parseNumber(a)
		if err != nil {
			warn("the %s %q is not a list of numbers, so the writing stays where it was", name, raw)
			return
		}
		vals[i] = v
	}
	seen := 0
	for i := range runs {
		if seen >= len(vals) {
			return
		}
		letters := len([]rune(runs[i].Text))
		if letters == 0 {
			continue
		}
		take := min(len(vals)-seen, letters)
		placed := runs[i].HasX
		if name == "y" {
			placed = runs[i].HasY
		}
		if !placed {
			if name == "y" {
				runs[i].YList = vals[seen : seen+take]
			} else {
				runs[i].XList = vals[seen : seen+take]
			}
		}
		seen += letters
	}
}

// coord is one coordinate an element wrote as a single number, which is what
// `dx` and `dy` are: a move of the pen for the piece, never one for each
// letter. A number that is not one leaves the writing where it was, and a list
// of them is said rather than being read as though it were the first of them
// by accident.
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

// textPathOutline is the outline a `<textPath>` writes along: the href followed
// the way a `<use>` follows its own, to a `<path>` this drawing has. Anything
// else — no href, one into another drawing, an id the drawing does not have, a
// shape that is not a path — has nowhere to put the writing, and says so rather
// than leaving it to come out somewhere the file did not ask for.
func (img *Image) textPathOutline(e *element, warn func(string, ...any)) *canvas.Path {
	href := e.attr("href")
	if href == "" {
		href = e.attr("xlink:href")
	}
	if href == "" {
		warn("it has no href, so there is no outline to write along")
		return nil
	}
	id, ok := strings.CutPrefix(href, "#")
	if !ok || id == "" {
		warn("the href %q names something in another drawing, which is not read", href)
		return nil
	}
	target := img.ids[id]
	if target == nil {
		warn("the href names #%s, which the drawing does not have", id)
		return nil
	}
	if target.Name != "path" {
		warn("the href names a <%s>, and writing follows a <path>", target.Name)
		return nil
	}
	p := pathData(target.attr("d"), warn)
	if p == nil || p.Empty() {
		warn("the outline the writing follows has nothing in it")
		return nil
	}
	return p
}

// startOffset is how far along the outline the writing begins. A percentage is
// a share of how long the outline is, which only the outline can answer, and
// any other length is a number of the drawing's own units.
func startOffset(e *element, along *canvas.Path, warn func(string, ...any)) float64 {
	raw := strings.TrimSpace(e.attr("startOffset"))
	if raw == "" {
		return 0
	}
	if body, ok := strings.CutSuffix(raw, "%"); ok {
		share, err := strconv.ParseFloat(body, 64)
		if err != nil {
			warn("the startOffset %q is not a number, so the writing starts at the beginning of the outline", raw)
			return 0
		}
		return along.Length() * share / 100
	}
	return readLength(warn, e, "startOffset", "startOffset")
}

// pointAt is where a distance along an outline lands: the point itself and the
// way the outline is going there, which is the direction the letter at that
// place faces. Past the end of the outline the last point and the last
// direction are what the writing hangs on, since an outline shorter than the
// writing on it still has to put that writing somewhere, and a distance before
// the start runs back along the first stretch for the same reason.
func pointAt(p *canvas.Path, d float64) (x, y, angle float64) {
	pts, closed := p.Points()
	for i, sub := range pts {
		segs := len(sub) - 1
		if closed[i] {
			segs++
		}
		for j := 0; j < segs; j++ {
			a, b := sub[j], sub[j+1]
			if j == segs-1 && closed[i] {
				b = sub[0]
			}
			seg := math.Hypot(b.X-a.X, b.Y-a.Y)
			if seg <= 0 {
				continue
			}
			if d <= seg {
				t := d / seg
				return a.X + (b.X-a.X)*t, a.Y + (b.Y-a.Y)*t,
					math.Atan2(b.Y-a.Y, b.X-a.X)
			}
			d -= seg
			x, y, angle = b.X, b.Y, math.Atan2(b.Y-a.Y, b.X-a.X)
		}
	}
	return x, y, angle
}

// trimEdges takes the whitespace off the two ends of the writing, which is
// what the line breaks around a `<text>` in the file are there for — they make
// the file readable and say nothing about where the writing starts and ends.
// A piece left with nothing in it is dropped, and where it would have started
// goes onto the writing after it: the place was in the file and still counts.
//
// preserve is `xml:space="preserve"`, and where it is on there is nothing to
// take off: the writing is drawn the way the file wrote it, spaces at the two
// ends of it included, because that is what the drawing asked for.
func trimEdges(runs []TextRun, preserve bool) []TextRun {
	if len(runs) == 0 {
		return runs
	}
	if !preserve {
		// The places a drawing gave one per letter are kept in step with the
		// letters they belong to: the values the dropped spaces were sitting
		// on go with them, so the letter that is now first is still at the
		// value the list counted out for it.
		n := len([]rune(runs[0].Text))
		runs[0].Text = strings.TrimLeft(runs[0].Text, " ")
		runs[0].XList = dropFront(runs[0].XList, n-len([]rune(runs[0].Text)))
		runs[0].YList = dropFront(runs[0].YList, n-len([]rune(runs[0].Text)))
		last := len(runs) - 1
		n = len([]rune(runs[last].Text))
		runs[last].Text = strings.TrimRight(runs[last].Text, " ")
		keep := len([]rune(runs[last].Text))
		runs[last].XList = dropPast(runs[last].XList, keep, n)
		runs[last].YList = dropPast(runs[last].YList, keep, n)
	}
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

// dropFront takes the first k values off a list of places, which is what goes
// with the k letters dropped from the front of the writing.
func dropFront(list []float64, k int) []float64 {
	if k <= 0 {
		return list
	}
	if len(list) > k {
		return list[k:]
	}
	return nil
}

// dropPast keeps the values that still have a letter to sit on after the back
// of the writing has been trimmed: n letters were there, keep are now.
func dropPast(list []float64, keep, n int) []float64 {
	if keep >= n {
		return list
	}
	if len(list) > keep {
		return list[:keep]
	}
	return list
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
//
// The writing is walked before it is drawn rather than as it is drawn, because
// the fill may be a colour at every point — a gradient or a pattern — and
// neither can answer which colour a pixel takes until the box the writing
// covers is known. The stroke goes down before the fill the same way a shape's
// does, so the part of it that falls inside the letters is covered by the fill
// and the part that falls outside is the stroke showing.
func paintText(cv *canvas.Canvas, n *Node, m canvas.Matrix, current canvas.Color, patterning []*pattern) {
	pieces := walkRuns(n.Runs, m)
	if len(pieces) == 0 {
		return
	}
	// The box the writing covers, as the drawing wrote it rather than as the
	// transform moved it: the same box a shape's geometry is measured to, and
	// the one a gradient or a pattern in fractions of the box it paints is a
	// fraction of. Writing is the only thing here with no box of its own until
	// a font is asked, and this is the font being asked.
	measured := *n
	measured.box = writingBox(pieces)
	for i := range pieces {
		paintPiece(cv, &measured, &pieces[i], m, current, patterning)
	}
}

// piece is one run of writing carried to where it goes: the face it is written
// in at the size it comes out at on the canvas, how far it comes out in the
// drawing's own units, where the pen stood when it started, and the matrix that
// takes it from there — the whole of what putting it down needs.
type piece struct {
	run        TextRun
	face       *canvas.Face
	scale      float64
	width      float64
	penX, penY float64
	at         canvas.Matrix
}

// walkRuns carries every run of writing to where the drawing puts it. Measuring
// and painting both take this walk, so that the box a fill is measured against
// and the place the letters land cannot come out different.
func walkRuns(runs []TextRun, m canvas.Matrix) []piece {
	pieces := make([]piece, 0, len(runs))
	penX, penY := 0.0, 0.0
	for _, run := range runs {
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
		at := m.Mul(st.Transform)
		scale := scaleOf(at)
		if scale <= 0 {
			continue
		}
		size := st.FontSize
		if size <= 0 {
			size = defaultFontSize
		}
		// The face is asked for at the size in pixels this writing covers on
		// the canvas, so a drawing painted big is written big rather than a
		// small piece of writing blown up afterwards. The face itself is the
		// one the canvas draws with: a font-family list is read, but a drawing
		// carries no files of its own, so the names run into the face the
		// program already had — see [Style].
		face := canvas.DefaultFace().AtSize(size * scale)
		width := (float64(face.Width(run.Text)) +
			st.LetterSpacing*float64(spaced(st, run.Text))) / scale
		// An anchor hangs the writing on the pen, and a piece hung on an
		// outline is hung along the outline instead — the pen has no part in
		// where it lands, and does not carry the pen on afterwards either.
		along := run.Path != nil
		if !along {
			switch st.Anchor {
			case AnchorMiddle:
				penX -= width / 2
			case AnchorEnd:
				penX -= width
			}
		}
		pieces = append(pieces, piece{
			run: run, face: face, scale: scale, width: width,
			penX: penX, penY: penY, at: at,
		})
		if !along {
			penX += width
		}
	}
	return pieces
}

// writingBox is the box the writing covers in the drawing's own coordinates,
// measured where the letters are before the transform moves them.
func writingBox(pieces []piece) [4]float64 {
	var b [4]float64
	found := false
	add := func(x0, y0, x1, y1 float64) {
		if !found {
			b = [4]float64{x0, y0, x1, y1}
			found = true
			return
		}
		b[0], b[1] = min(b[0], x0), min(b[1], y0)
		b[2], b[3] = max(b[2], x1), max(b[3], y1)
	}
	for _, p := range pieces {
		if p.run.Path != nil {
			// The piece is hung along an outline rather than at a pen, so what
			// it covers is where the outline is rather than where the pen was.
			if x0, y0, x1, y1, ok := p.run.Path.Bounds(); ok {
				add(x0, y0, x1, y1)
			}
			continue
		}
		ascent := float64(p.face.Ascent()) / p.scale
		descent := float64(p.face.Descent()) / p.scale
		// A piece with a line break in it is as many lines tall as it has
		// breaks, and the ones below the first are below the pen.
		down := float64(1+strings.Count(p.run.Text, "\n")) *
			float64(p.face.Height()) / p.scale
		add(p.penX, p.penY-ascent, p.penX+p.width, p.penY+descent+down)
	}
	if !found {
		return [4]float64{}
	}
	return b
}

// putRun is one piece of the writing laid down in a canvas: in the colour it
// is given, moved by however far the caller asks. The fill asks for no move at
// all, and the stroke asks for one for each point around the circle its width
// makes — which is the whole of what a stroke of that width covers.
type putRun func(target *canvas.Canvas, c canvas.Color, dx, dy float64)

// paintPiece puts one run of writing down: its stroke first and its fill over
// the stroke, which is the order SVG paints and the order that leaves the half
// of a stroke inside the letters covered by the fill and the half outside
// showing as the outline it is.
func paintPiece(cv *canvas.Canvas, n *Node, p *piece, m canvas.Matrix,
	current canvas.Color, patterning []*pattern) {
	st := p.run.Style
	put := func(target *canvas.Canvas, c canvas.Color, dx, dy float64) {
		switch {
		case p.run.Path != nil:
			writeAlong(target, p.face, p.run, p.width, p.scale, c,
				p.at.Mul(canvas.Translate(dx, dy)))
		case len(p.run.XList) > 0 || len(p.run.YList) > 0 || st.LetterSpacing != 0:
			// A place for each letter, or a gap left after every one of them,
			// or both: either puts the letters down one at a time, because the
			// pen has to be moved between them rather than handed the whole
			// piece and told to write it in one go. The pen is not in the
			// matrix here — the places a list names are places of the drawing
			// itself, not distances from where the pen already was.
			writeLetters(target, p.face, p.run, p.at, p.scale, c,
				p.penX+dx, p.penY+dy, p.width)
		default:
			at := p.at.Mul(canvas.Translate(p.penX+dx, p.penY+dy))
			drawRun(target, p.face, p.run.Text, c, at)
			if heavy(st.FontWeight) {
				// The weight the face does not carry is the same writing a
				// second time, one pixel further on.
				drawRun(target, p.face, p.run.Text, c,
					at.Mul(canvas.Translate(1/p.scale, 0)))
			}
		}
	}
	if st.HasStroke && st.Width > 0 {
		c := st.Stroke
		if st.StrokeCurrent {
			c = current
		}
		if c := fade(c, st.Opacity*st.StrokeOpacity); c.A() > 0 {
			strokeWriting(cv, put, st.strokeWidth(m), c, p.scale)
		}
	}
	var shade func(x, y int) canvas.Color
	var c canvas.Color
	switch {
	case !st.HasFill:
		return
	case st.fillGradient != nil:
		shade = gradientShade(st.fillGradient, n, m, current, st.Opacity*st.FillOpacity)
	case st.fillPattern != nil:
		// Where the pattern has nothing to paint with — it names itself, the
		// tile has no size to it, the transform has no way back — the colour
		// after the `url(...)` is what the writing falls back on, and nothing
		// at all where it brought none.
		if s, ok := patternShade(st.fillPattern, n, m, cv.Width, cv.Height, current,
			st.Opacity*st.FillOpacity, patterning); ok {
			shade = s
		} else if fb, ok := st.fillFallbackColour(current); ok {
			c = fade(fb, st.Opacity*st.FillOpacity)
		}
	default:
		c = st.Fill
		if st.FillCurrent {
			c = current
		}
		c = fade(c, st.Opacity*st.FillOpacity)
	}
	if shade == nil {
		if c.A() > 0 {
			put(cv, c, 0, 0)
		}
		return
	}
	layShaded(cv, put, shade)
}

// strokeWriting draws what a stroke over the writing covers: the same writing
// at every point around a circle as wide across as the stroke, which is the
// shape a stroke of that width makes around letters. The circle is walked
// finely enough that its points land about half a pixel apart, and a stroke
// narrower than a pixel is taken at half a pixel rather than at nothing, so a
// hairline still shows after the picture is brought back down to size.
//
// The outline comes out round rather than at the miter a stroke takes round the
// corner of a shape: the corners a glyph has are smaller than the difference
// would be at the sizes writing is drawn at.
func strokeWriting(cv *canvas.Canvas, put putRun, width float64, c canvas.Color, scale float64) {
	r := math.Max(width/2, 0.5)
	n := min(max(int(math.Ceil(4*math.Pi*r)), 8), 128)
	for i := range n {
		a := 2 * math.Pi * float64(i) / float64(n)
		put(cv, c, math.Cos(a)*r/scale, math.Sin(a)*r/scale)
	}
}

// layShaded paints the writing where the colour comes from where each pixel
// lands rather than being one colour over all of it, which is what a gradient
// and a pattern both are. The letters go into a picture of their own first, so
// that there is a coverage per pixel to ask about, and the colour is then laid
// over what is already there one pixel at a time — the same question a shape
// filled with either asks of every pixel it covers, through FillPathFunc.
func layShaded(cv *canvas.Canvas, put putRun, shade func(x, y int) canvas.Color) {
	layer, err := canvas.NewLayer(cv.Width, cv.Height)
	if err != nil {
		return
	}
	// Any opaque colour will do: only what it covered matters, and the colour
	// comes from the shade at the pixel rather than from this.
	put(layer, canvas.White, 0, 0)
	for y := range layer.Height {
		for x := range layer.Width {
			a := layer.At(x, y).A()
			if a == 0 {
				continue
			}
			c := shade(x, y)
			if c.A() == 0 {
				continue
			}
			if a < 255 {
				c = canvas.RGBA(c.R(), c.G(), c.B(), uint8(int(c.A())*int(a)/255))
			}
			cv.Pixel(x, y, c)
		}
	}
}

// heavy is a weight the canvas draws heavier than the face it has: a face is
// cut in the weights it is cut in, and the ones it does not have are faked by
// drawing the writing a second time a pixel on.
func heavy(weight uint16) bool {
	if weight == 0 {
		weight = FontWeightNormal
	}
	return weight >= 600
}

// spaced is how many letters of a piece the drawing asked for room after: a
// line break is the end of a line rather than a letter at the end of one, and
// is left out of the count.
func spaced(st Style, text string) int {
	if st.LetterSpacing == 0 {
		return 0
	}
	n := 0
	for _, r := range text {
		if r != '\n' && r != '\r' {
			n++
		}
	}
	return n
}

// writeLetters puts a piece of writing down one letter at a time: each letter
// at the place the drawing named for it when it named a place for every letter
// of the piece, and otherwise at the pen the letter before it left, with the
// letter-spacing the drawing asked for left after every one of them. Bold is
// the letter drawn a second time a pixel further on, the same as it is when
// the whole piece goes down at once.
//
// m is the transform of the writing without the pen's place in it, because a
// place a list names is a place of the drawing rather than a distance from
// wherever the pen was; penX and penY are the pen already moved by the anchor,
// for the letters the list did not reach, and width is what the anchor worked
// that movement out from.
func writeLetters(cv *canvas.Canvas, face *canvas.Face, run TextRun, m canvas.Matrix,
	scale float64, c canvas.Color, penX, penY, width float64) {
	st := run.Style
	spacing := st.LetterSpacing
	bold := heavy(st.FontWeight)
	shift := 0.0
	switch st.Anchor {
	case AnchorMiddle:
		shift = -width / 2
	case AnchorEnd:
		shift = -width
	}
	lineStart := penX
	for i, r := range run.Text {
		switch r {
		case '\r':
			continue
		case '\n':
			penX = lineStart
			penY += float64(face.Height()) / scale
			continue
		case '\t':
			penX += 4 * float64(face.Width(" ")) / scale
			continue
		}
		x, y := penX, penY
		if i < len(run.XList) {
			x = run.XList[i] + shift
		}
		if i < len(run.YList) {
			y = run.YList[i]
		}
		at := m.Mul(canvas.Translate(x, y))
		drawRun(cv, face, string(r), c, at)
		if bold {
			drawRun(cv, face, string(r), c, at.Mul(canvas.Translate(1/scale, 0)))
		}
		penX, penY = x+float64(face.Width(string(r)))/scale+spacing, y
	}
}

// writeAlong lays one piece of writing along an outline: every letter is put
// down at the place the outline has reached by the time the letters before it
// have been written, turned to the way the outline is going there, so the
// writing curves with the line rather than lying flat across it.
func writeAlong(cv *canvas.Canvas, face *canvas.Face, run TextRun, width, scale float64, c canvas.Color, m canvas.Matrix) {
	start := run.Offset
	switch run.Style.Anchor {
	case AnchorMiddle:
		start -= width / 2
	case AnchorEnd:
		start -= width
	}
	d := start
	spacing := run.Style.LetterSpacing
	bold := heavy(run.Style.FontWeight)
	for _, r := range run.Text {
		if r == '\n' || r == '\r' {
			// An outline is one line: there is nowhere below it to go.
			continue
		}
		s := string(r)
		if r == '\t' {
			s = "    "
		}
		adv := float64(face.Width(s)) / scale
		if r != ' ' && r != '\t' {
			x, y, ang := pointAt(run.Path, d)
			at := m.Mul(canvas.Translate(x, y)).Mul(canvas.Rotate(ang))
			drawRun(cv, face, s, c, at)
			if bold {
				drawRun(cv, face, s, c, at.Mul(canvas.Translate(1/scale, 0)))
			}
		}
		d += adv + spacing
	}
}

// drawRun puts one piece of writing on the canvas. The matrix is the whole of
// where it goes: the pen's own transform with the pen's place already in it,
// and for a letter along an outline, the turn of the outline at that letter
// too. A matrix that only moves and sizes the writing puts each letter down
// where the pen landed, which is the whole of what has to happen.
//
// A matrix that turns the writing is the other case: the letters are pixels
// and a pixel cannot lean, so the piece is written as it stands into a picture
// of its own and that picture is laid over through the turn — the same thing
// the shapes do, since a shape is a path taken through the transform before it
// is rasterised. The picture is made at the resolution the face was asked for
// and the matrix carries it the rest of the way, so the letters are rasterised
// once at their own size and stretched by the turn rather than drawn twice.
func drawRun(cv *canvas.Canvas, face *canvas.Face, text string, c canvas.Color, m canvas.Matrix) {
	if m.B == 0 && m.C == 0 {
		px, py := m.Map(0, 0)
		face.Draw(cv, int(px+0.5), int(py+0.5), text, c)
		return
	}
	const pad = 2
	scale := scaleOf(m)
	// One drawing unit is `scale` of the face's pixels, so the face's own
	// measurement is already the size the picture has to be. A piece with a
	// line break in it is as many lines tall as it has breaks.
	lines := 1 + strings.Count(text, "\n")
	w := face.Width(text) + 2*pad
	h := face.Ascent() + face.Descent() + (lines-1)*face.Height() + 2*pad
	if w <= 0 || h <= 0 {
		return
	}
	layer, err := canvas.NewLayer(w, h)
	if err != nil {
		return
	}
	face.Draw(layer, pad, pad+face.Ascent(), text, c)
	// Layer pixel (lx,ly) is the drawing's own (lx-pad)/scale and
	// (ly-pad-Ascent)/scale, and the matrix takes those from there onwards.
	to := m.Mul(canvas.Translate(-float64(pad)/scale,
		-float64(pad+face.Ascent())/scale).Mul(canvas.Scale(1/scale, 1/scale)))
	cv.BlitMatrix(layer, to, canvas.Area{Width: w, Height: h})
}
