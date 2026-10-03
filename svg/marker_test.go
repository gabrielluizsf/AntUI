package svg

import (
	"strings"
	"testing"
)

// A `<marker>` is a picture drawn at the vertex of the shape that names it, and
// what puts it there is the shape rather than the marker: the vertex itself,
// the direction the path is going as it passes through that vertex, and the
// room the marker asked for, which by default is a multiple of the stroke-width
// of this shape. `marker-start` takes the first vertex, `marker-end` the last
// and `marker-mid` every vertex between the two. The marker is written once,
// far from where it is used, and comes out the same at every vertex that names
// it.

// markerDrawing is a drawing with one marker in it and one line naming it, for
// the tests that ask where the marker landed. The line has no stroke, which is
// on purpose: a marker is a picture put on a vertex, and nothing about the line
// is what puts it there.
func markerDrawing(marker, line string) string {
	return `<svg viewBox="0 0 20 20">
		<defs>
			<marker id="m" ` + marker + `>
				<rect x="0" y="0" width="4" height="4" fill="red"/>
			</marker>
		</defs>
		<line x1="10" y1="10" x2="20" y2="10" stroke="none" ` + line + `/>
	</svg>`
}

// linesOf is every line of a drawing, in the order they were written, which is
// what a test about what each of them drew asks for. [Node.find] answers the
// first one and no more.
func linesOf(img *Image) []*Node {
	var out []*Node
	var walk func(n *Node)
	walk = func(n *Node) {
		if n.Name == "line" {
			out = append(out, n)
		}
		for _, k := range n.Kids {
			walk(k)
		}
	}
	walk(img.Root)
	return out
}

func TestMarkerStartIsRead(t *testing.T) {
	// The reference is read the same way a fill's is: `url(#id)` for the whole
	// drawing to follow, `none` for no marker at all, and anything else said
	// out loud rather than quietly drawn without the marker that was asked for.
	for _, tc := range []struct {
		name   string
		value  string
		warns  int
		wantID string
	}{
		{"the marker the drawing has", "url(#m)", 0, "m"},
		{"no marker at all", "none", 0, ""},
		{"a marker the drawing does not have", "url(#nowhere)", 1, ""},
		{"a reference this package cannot follow", "#m", 1, ""},
		{"nothing written in it", "", 1, ""},
	} {
		img, err := Parse(markerDrawing(`markerUnits="userSpaceOnUse"`, `marker-start="`+tc.value+`"`))
		if err != nil {
			t.Fatalf("%s: parse: %v", tc.name, err)
		}
		if got := len(img.Warnings()); got != tc.warns {
			t.Errorf("%s: %d warnings, want %d (%v)", tc.name, got, tc.warns, img.Warnings())
		}
		line := img.Root.find("line")
		if line == nil {
			t.Fatalf("%s: the line is not in the drawing", tc.name)
		}
		got := ""
		if line.Style.markerStart != nil {
			got = line.Style.markerStart.id
		}
		if got != tc.wantID {
			t.Errorf("%s: marker-start is #%s, want #%s", tc.name, got, tc.wantID)
		}
	}
}

func TestMarkerIsReadFromTheStyleAttribute(t *testing.T) {
	// A declaration says the same thing the attribute does, and the one written
	// in `style` is the one that counts where both are written. All three are
	// read the same way.
	st := Style{}.with(mustElement(t,
		`<line marker-start="none" style="marker-start: url(#a)"/>`), quiet)
	if st.markerStartRef != "a" {
		t.Errorf("marker-start = #%s, want #a", st.markerStartRef)
	}
	st = Style{}.with(mustElement(t,
		`<line marker-start="url(#a)" style="marker-start: none"/>`), quiet)
	if st.markerStartRef != "" {
		t.Errorf("marker-start = #%s, want none", st.markerStartRef)
	}
	st = Style{}.with(mustElement(t,
		`<line marker-mid="none" style="marker-mid: url(#b); marker-end: url(#c)"/>`), quiet)
	if st.markerMidRef != "b" || st.markerEndRef != "c" {
		t.Errorf("marker-mid = #%s and marker-end = #%s, want #b and #c", st.markerMidRef, st.markerEndRef)
	}
}

func TestMarkerStartIsInheritedAndSaysItselfOnce(t *testing.T) {
	// A marker is inherited, which is what makes a `<g marker-start>` the mark
	// on every vertex below it — unlike a clip, a mask or a filter, which are
	// spent where they are written. And a reference the drawing cannot follow
	// is said once for the group rather than once for every shape inside it,
	// because the reference is spent when it is read and what the children take
	// is the answer rather than the question.
	img, err := Parse(`<svg viewBox="0 0 40 40">
		<defs>
			<marker id="a"><rect width="4" height="4" fill="red"/></marker>
			<marker id="b"><rect width="4" height="4" fill="blue"/></marker>
		</defs>
		<g marker-start="url(#a)">
			<line x1="0" y1="4" x2="10" y2="4" stroke="none"/>
			<line x1="0" y1="8" x2="10" y2="8" stroke="none" marker-start="none"/>
			<line x1="0" y1="12" x2="10" y2="12" stroke="none" marker-start="url(#b)"/>
		</g>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if ws := img.Warnings(); len(ws) != 0 {
		t.Fatalf("the drawing said %v, want nothing at all", ws)
	}
	lines := linesOf(img)
	if len(lines) != 3 {
		t.Fatalf("the drawing has %d lines, want 3", len(lines))
	}
	if lines[0].Style.markerStart == nil || lines[0].Style.markerStart.id != "a" {
		t.Error("the first line did not take the marker the group named")
	}
	if lines[1].Style.markerStart != nil {
		t.Error("a shape with marker-start=\"none\" took the marker from the group above it anyway")
	}
	if lines[2].Style.markerStart == nil || lines[2].Style.markerStart.id != "b" {
		t.Error("a shape's own marker-start did not beat the one the group named")
	}

	missing, err := Parse(`<svg viewBox="0 0 40 40">
		<g marker-start="url(#nowhere)">
			<line x1="0" y1="4" x2="10" y2="4" stroke="none"/>
			<line x1="0" y1="8" x2="10" y2="8" stroke="none"/>
		</g>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if ws := missing.Warnings(); len(ws) != 1 {
		t.Errorf("a marker nowhere was said %d times, want once (%v)", len(ws), ws)
	}
}

func TestMarkerMidAndEndAreRead(t *testing.T) {
	// The other two markers are read the same way the first one is: the same
	// reference, the same `none`, and the same say-it-out-loud where the
	// drawing cannot follow what was written, because a shape asked for a
	// marker and got none is something an author should be able to find out.
	for _, tc := range []struct {
		name   string
		value  string
		warns  int
		wantID string
	}{
		{"the marker the drawing has", "url(#m)", 0, "m"},
		{"no marker at all", "none", 0, ""},
		{"a marker the drawing does not have", "url(#nowhere)", 1, ""},
		{"a reference this package cannot follow", "#m", 1, ""},
		{"nothing written in it", "", 1, ""},
	} {
		for _, prop := range []string{"marker-mid", "marker-end"} {
			img, err := Parse(markerDrawing(`markerUnits="userSpaceOnUse"`, prop+`="`+tc.value+`"`))
			if err != nil {
				t.Fatalf("%s of %s: parse: %v", prop, tc.name, err)
			}
			if got := len(img.Warnings()); got != tc.warns {
				t.Errorf("%s of %s: %d warnings, want %d (%v)", prop, tc.name, got, tc.warns, img.Warnings())
			}
			line := img.Root.find("line")
			if line == nil {
				t.Fatalf("%s of %s: the line is not in the drawing", prop, tc.name)
			}
			got := ""
			if prop == "marker-mid" && line.Style.markerMid != nil {
				got = line.Style.markerMid.id
			} else if prop == "marker-end" && line.Style.markerEnd != nil {
				got = line.Style.markerEnd.id
			}
			if got != tc.wantID {
				t.Errorf("%s of %s: is #%s, want #%s", prop, tc.name, got, tc.wantID)
			}
		}
	}
}

func TestEachMarkerPropertyIsTurnedOffOnItsOwn(t *testing.T) {
	// The three markers are inherited one by one, and a shape that writes one
	// of them drops only that one: `marker-start="none"` on a line under a
	// group that asked for a mid and an end leaves both of those where they
	// were. And a reference the drawing cannot follow is said once for the
	// group rather than once for every shape inside it, because the reference
	// is spent when it is read and what the children take is the answer rather
	// than the question.
	img, err := Parse(`<svg viewBox="0 0 60 40">
		<defs>
			<marker id="a"><rect width="4" height="4" fill="red"/></marker>
			<marker id="b"><rect width="4" height="4" fill="blue"/></marker>
		</defs>
		<g marker-mid="url(#a)" marker-end="url(#b)">
			<line x1="0" y1="4" x2="10" y2="4" stroke="none"/>
			<line x1="0" y1="8" x2="10" y2="8" stroke="none" marker-start="none"/>
			<line x1="0" y1="12" x2="10" y2="12" stroke="none" marker-end="none"/>
		</g>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if ws := img.Warnings(); len(ws) != 0 {
		t.Fatalf("the drawing said %v, want nothing at all", ws)
	}
	lines := linesOf(img)
	if len(lines) != 3 {
		t.Fatalf("the drawing has %d lines, want 3", len(lines))
	}
	for i, l := range lines {
		if l.Style.markerMid == nil || l.Style.markerMid.id != "a" {
			t.Errorf("line %d did not take the marker-mid the group named", i)
		}
	}
	if lines[0].Style.markerEnd == nil || lines[0].Style.markerEnd.id != "b" {
		t.Error("the first line did not take the marker-end the group named")
	}
	if lines[1].Style.markerEnd == nil || lines[1].Style.markerEnd.id != "b" {
		t.Error("turning the marker-start off took the marker-end with it")
	}
	if lines[2].Style.markerEnd != nil {
		t.Error("a shape with marker-end=\"none\" took the marker from the group above it anyway")
	}

	missing, err := Parse(`<svg viewBox="0 0 60 40">
		<g marker-end="url(#nowhere)">
			<line x1="0" y1="4" x2="10" y2="4" stroke="none"/>
			<line x1="0" y1="8" x2="10" y2="8" stroke="none"/>
		</g>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if ws := missing.Warnings(); len(ws) != 1 {
		t.Errorf("a marker nowhere was said %d times, want once (%v)", len(ws), ws)
	}
}

func TestMarkerAttributesAreRead(t *testing.T) {
	// Everything a marker says about itself, and the defaults it gets where it
	// says nothing: a room three stroke-widths square, the reference at its
	// origin, no turn on it, and the size in units of the stroke-width — which
	// is what makes a marker on a thick line a thick marker.
	for _, tc := range []struct {
		name         string
		attrs        string
		view         [4]float64
		hasView      bool
		refX, refY   float64
		w, h         float64
		strokeUnits  bool
		deg          float64
		auto         bool
		startReverse bool
		warns        int
	}{
		{"what a marker says where it says nothing", "", [4]float64{}, false, 0, 0, 3, 3, true, 0, false, false, 0},
		{"everything it may say", `viewBox="0 0 10 10" refX="2" refY="3" markerWidth="8" markerHeight="4" markerUnits="userSpaceOnUse"`,
			[4]float64{0, 0, 10, 10}, true, 2, 3, 8, 4, false, 0, false, false, 0},
		{"an angle in degrees", `orient="45deg"`, [4]float64{}, false, 0, 0, 3, 3, true, 45, false, false, 0},
		{"an angle with no unit", `orient="90"`, [4]float64{}, false, 0, 0, 3, 3, true, 90, false, false, 0},
		{"turning to the direction of the path", `orient="auto"`, [4]float64{}, false, 0, 0, 3, 3, true, 0, true, false, 0},
		{"the same turned around at the start", `orient="auto-start-reverse"`, [4]float64{}, false, 0, 0, 3, 3, true, 0, true, true, 0},
		{"an orient this package cannot read", `orient="sideways"`, [4]float64{}, false, 0, 0, 3, 3, true, 0, false, false, 1},
		{"a length this package cannot read", `markerWidth="50%"`, [4]float64{}, false, 0, 0, 3, 3, true, 0, false, false, 1},
		{"a room with nothing in it", `markerWidth="0"`, [4]float64{}, false, 0, 0, 0, 3, true, 0, false, false, 1},
		{"a viewBox this package cannot read", `viewBox="0 0 10"`, [4]float64{}, false, 0, 0, 3, 3, true, 0, false, false, 1},
	} {
		img, err := Parse(`<svg viewBox="0 0 20 20">
			<defs>
				<marker id="m" ` + tc.attrs + `><rect width="4" height="4" fill="red"/></marker>
			</defs>
			<line x1="0" y1="10" x2="10" y2="10" stroke="none" marker-start="url(#m)"/>
		</svg>`)
		if err != nil {
			t.Fatalf("%s: parse: %v", tc.name, err)
		}
		if got := len(img.Warnings()); got != tc.warns {
			t.Errorf("%s: %d warnings, want %d (%v)", tc.name, got, tc.warns, img.Warnings())
		}
		m := img.markers["m"]
		if m == nil {
			t.Fatalf("%s: the marker is not in the drawing", tc.name)
		}
		if m.hasView != tc.hasView {
			t.Errorf("%s: hasView = %v, want %v", tc.name, m.hasView, tc.hasView)
		} else if m.hasView && m.view != tc.view {
			t.Errorf("%s: viewBox = %v, want %v", tc.name, m.view, tc.view)
		}
		if m.refX != tc.refX || m.refY != tc.refY {
			t.Errorf("%s: reference = (%g, %g), want (%g, %g)", tc.name, m.refX, m.refY, tc.refX, tc.refY)
		}
		if m.w != tc.w || m.h != tc.h {
			t.Errorf("%s: room = %g by %g, want %g by %g", tc.name, m.w, m.h, tc.w, tc.h)
		}
		if m.strokeUnits != tc.strokeUnits {
			t.Errorf("%s: markerUnits in units of the stroke-width = %v, want %v", tc.name, m.strokeUnits, tc.strokeUnits)
		}
		if !nearFloat(m.deg, tc.deg, 1e-9) {
			t.Errorf("%s: orient = %g degrees, want %g", tc.name, m.deg, tc.deg)
		}
		if m.auto != tc.auto || m.startReverse != tc.startReverse {
			t.Errorf("%s: orient = auto %v / start-reverse %v, want %v / %v",
				tc.name, m.auto, m.startReverse, tc.auto, tc.startReverse)
		}
	}
}

func TestRenderDrawsAMarkerAtTheFirstVertexWithoutAStroke(t *testing.T) {
	// The line asks for no stroke at all, and the marker is drawn anyway: a
	// marker is a picture put on the vertex, not something the line has to be
	// outlined for. Nothing was written for the reference point, so the
	// top-left corner of the picture lands on the first point of the path, and
	// the room the marker asked for bounds what it draws from there.
	_, cv := painted(t, markerDrawing(
		`markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="0"`,
		`marker-start="url(#m)"`), 20, 20)
	if c := pixelAt(cv, 12, 12); c.R() < 200 || c.G() > 60 || c.B() > 60 {
		t.Errorf("the marker is not at the vertex: pixel (12,12) = %v, want red", c)
	}
	if c := pixelAt(cv, 16, 16); c.A() != 0 {
		t.Errorf("the marker drew outside the rectangle it was written as: pixel (16,16) = %v", c)
	}
	if c := pixelAt(cv, 8, 8); c.A() != 0 {
		t.Errorf("the marker drew beside the vertex rather than on it: pixel (8,8) = %v", c)
	}
}

func TestRenderTurnsAMarkerToTheDirectionOfThePath(t *testing.T) {
	// `auto` turns the marker to the direction the path is going as it leaves
	// the vertex, which is what makes an arrowhead point along its own line;
	// `auto-start-reverse` turns the one at the start of a path around, so that
	// one marker can be put at both ends of a line and point out of each. With
	// neither, the marker is drawn the way it was written — and on a line
	// running left to right, the way it was written and the way `auto` turns it
	// are the same way.
	for _, tc := range []struct {
		name   string
		orient string
		want   int
	}{
		{"the marker as it was written", "", 12},
		{"turned to the direction of the path", "auto", 12},
		{"turned around at the start of it", "auto-start-reverse", 8},
	} {
		_, cv := painted(t, markerDrawing(
			`markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="2" orient="`+tc.orient+`"`,
			`marker-start="url(#m)"`), 20, 20)
		here, other := pixelAt(cv, tc.want, 10), pixelAt(cv, 20-tc.want, 10)
		if here.A() == 0 {
			t.Errorf("%s: nothing was drawn at (%d,10), where the marker should be", tc.name, tc.want)
		}
		if other.A() != 0 {
			t.Errorf("%s: something was drawn at (%d,10), on the other side of the vertex", tc.name, 20-tc.want)
		}
	}
}

func TestRenderSizesAMarkerByTheStrokeWidth(t *testing.T) {
	// The room a marker asks for is a multiple of the stroke-width of the shape
	// it names, by default, which is what makes a marker on a thick line a
	// thick marker. Here the picture has a viewBox, so the room is what the
	// viewBox is fitted into and the marker comes out twice as big when the
	// line is twice as thick.
	for _, tc := range []struct {
		name  string
		width string
		at    int
		out   int
	}{
		{"a marker the width of a unit-thick line", "1", 15, 25},
		{"the same marker on a line twice as thick", "2", 25, 35},
	} {
		_, cv := painted(t, `<svg viewBox="0 0 40 20">
			<defs>
				<marker id="m" viewBox="0 0 5 5" refX="0" refY="0" markerWidth="10" markerHeight="10">
					<rect x="0" y="0" width="5" height="5" fill="red"/>
				</marker>
			</defs>
			<line x1="10" y1="10" x2="30" y2="10" stroke="none" stroke-width="`+tc.width+`" marker-start="url(#m)"/>
		</svg>`, 40, 20)
		if c := pixelAt(cv, tc.at, 15); c.A() == 0 {
			t.Errorf("%s: nothing was drawn at (%d,15), where the marker should reach", tc.name, tc.at)
		}
		if c := pixelAt(cv, tc.out, 15); c.A() != 0 {
			t.Errorf("%s: something was drawn at (%d,15), past the room the marker asked for", tc.name, tc.out)
		}
	}
}

func TestRenderMovesAMarkerWithTheShape(t *testing.T) {
	// The transform is put into the matrix that places the marker under the
	// same condition it is put into the path itself, so a marker lands on the
	// vertex where the vertex ended up rather than where it was written.
	_, cv := painted(t, `<svg viewBox="0 0 40 20">
		<defs>
			<marker id="m" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="0">
				<rect x="0" y="0" width="4" height="4" fill="red"/>
			</marker>
		</defs>
		<g transform="translate(10,0)">
			<line x1="10" y1="10" x2="30" y2="10" stroke="none" marker-start="url(#m)"/>
		</g>
	</svg>`, 40, 20)
	if c := pixelAt(cv, 22, 12); c.R() < 200 || c.G() > 60 || c.B() > 60 {
		t.Errorf("the marker did not follow the group: pixel (22,12) = %v, want red", c)
	}
	if c := pixelAt(cv, 12, 12); c.A() != 0 {
		t.Errorf("the marker stayed where it was written: pixel (12,12) = %v", c)
	}
}

func TestAMarkerDoesNotDrawItselfRoundForever(t *testing.T) {
	// A marker whose picture has a shape naming a marker of its own is a
	// circle: drawing it would draw it again, for ever. The reference that
	// closes the circle is left off with one warning about it, which draws that
	// vertex with no marker on it rather than drawing nothing at all — and this
	// test coming back is the point of it.
	img, cv := painted(t, `<svg viewBox="0 0 20 20">
		<defs>
			<marker id="m" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10">
				<path d="M0 0 L4 4" fill="none" stroke="none" marker-start="url(#m)"/>
			</marker>
		</defs>
		<line x1="10" y1="10" x2="20" y2="10" stroke="none" marker-start="url(#m)"/>
	</svg>`, 20, 20)
	ws := img.Warnings()
	if len(ws) != 1 || !strings.Contains(ws[0].What, "for ever") {
		t.Errorf("a marker going round said %v, want one warning about going round for ever", ws)
	}
	if c := pixelAt(cv, 12, 12); c.A() != 0 {
		t.Errorf("the vertex with the marker left off was drawn anyway: pixel (12,12) = %v", c)
	}
}

func TestTwoMarkersPointingAtEachOtherAreCutAtTheFirstRepeat(t *testing.T) {
	// The same circle, written the long way round: one marker inside the other,
	// and both of them hanging off the end rather than the start, which is the
	// same circle to the cut — it is not the property that makes a cycle, it is
	// the reference. It is cut the same way, at the first repeat of the marker
	// that is already being drawn, and the warning says which of the three it
	// was that went round.
	img, cv := painted(t, `<svg viewBox="0 0 20 20">
		<defs>
			<marker id="a" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10">
				<path d="M0 0 L4 4" fill="none" stroke="none" marker-end="url(#b)"/>
			</marker>
			<marker id="b" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10">
				<path d="M0 0 L4 4" fill="none" stroke="none" marker-end="url(#a)"/>
			</marker>
		</defs>
		<line x1="10" y1="10" x2="20" y2="10" stroke="none" marker-start="url(#a)"/>
	</svg>`, 20, 20)
	ws := img.Warnings()
	if len(ws) != 1 || !strings.Contains(ws[0].What, "marker-end") || !strings.Contains(ws[0].What, "for ever") {
		t.Errorf("two markers naming each other said %v, want one warning about the marker-end going round", ws)
	}
	if c := pixelAt(cv, 12, 12); c.A() != 0 {
		t.Errorf("the vertex with the marker left off was drawn anyway: pixel (12,12) = %v", c)
	}
}

func TestRenderDrawsTheMarkerOverTheStroke(t *testing.T) {
	// The marker comes after the whole stroke rather than in the middle of it,
	// which is what the spec puts them in: an arrowhead on the line it marks
	// the start of, not under it and half covered by it. Both cover the pixel
	// the vertex lands on, so which of them is there says which was drawn last.
	_, cv := painted(t, `<svg viewBox="0 0 20 20">
		<defs>
			<marker id="m" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="2">
				<rect x="0" y="0" width="4" height="4" fill="blue"/>
			</marker>
		</defs>
		<line x1="10" y1="10" x2="20" y2="10" stroke="red" stroke-width="2" marker-start="url(#m)"/>
	</svg>`, 20, 20)
	if c := pixelAt(cv, 12, 10); c.B() < 200 || c.R() > 60 {
		t.Errorf("the pixel at the vertex = %v, want the marker rather than the stroke", c)
	}
	if c := pixelAt(cv, 16, 10); c.R() < 200 || c.B() > 60 {
		t.Errorf("the pixel past the marker = %v, want the stroke rather than the marker", c)
	}
}

func TestRenderDrawsAMarkerAtEveryVertexItWasAskedFor(t *testing.T) {
	// The three properties between them cover every vertex of the path: the
	// first point takes marker-start, the last takes marker-end and the ones
	// between the two take marker-mid. Each marker is put down with its
	// reference point on the vertex it was asked for, so the four red pixels of
	// each picture start there and reach away from it.
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs>
			<marker id="m" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="0">
				<rect x="0" y="0" width="4" height="4" fill="red"/>
			</marker>
		</defs>
		<polyline points="5,15 15,15 15,30" fill="none" stroke="none"
			marker-start="url(#m)" marker-mid="url(#m)" marker-end="url(#m)"/>
	</svg>`, 40, 40)
	for _, at := range [][2]int{{6, 16}, {16, 16}, {16, 31}} {
		if c := pixelAt(cv, at[0], at[1]); c.R() < 200 || c.G() > 60 || c.B() > 60 {
			t.Errorf("the marker at vertex %v: pixel (%d,%d) = %v, want red", at, at[0], at[1], c)
		}
	}
	// Nothing between the markers: the line has no stroke to draw there.
	for _, at := range [][2]int{{11, 16}, {16, 22}} {
		if c := pixelAt(cv, at[0], at[1]); c.A() != 0 {
			t.Errorf("something was drawn at (%d,%d), between two markers", at[0], at[1])
		}
	}

	// And a line, which has no vertex between its two, draws no marker-mid at
	// all: marker-mid is every vertex except the first and the last.
	_, lv := painted(t, markerDrawing(`markerUnits="userSpaceOnUse"`, `marker-mid="url(#m)"`), 20, 20)
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			if c := pixelAt(lv, x, y); c.A() != 0 {
				t.Fatalf("a marker-mid was drawn on a line at (%d,%d): %v", x, y, c)
			}
		}
	}
}

func TestRenderTurnsAMidMarkerToTheMiddleOfTheTurn(t *testing.T) {
	// A vertex with a segment coming in and a segment going out has no single
	// direction to turn a marker to, and the spec's answer is the middle of
	// the turn — the two directions added as unit vectors. Here the path runs
	// east into the corner and leaves it going south, so the middle is south
	// east, and the marker is turned 45 degrees: the pixel south east of the
	// vertex has the marker on it, and the pixel east of the vertex, which the
	// marker as it was written would cover, does not.
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs>
			<marker id="m" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="2" orient="auto">
				<rect x="0" y="0" width="4" height="4" fill="red"/>
			</marker>
		</defs>
		<polyline points="5,15 15,15 15,30" fill="none" stroke="none" marker-mid="url(#m)"/>
	</svg>`, 40, 40)
	if c := pixelAt(cv, 17, 18); c.A() == 0 {
		t.Errorf("the marker did not turn to the middle of the turn: pixel (17,18) = %v, want the marker", c)
	}
	if c := pixelAt(cv, 18, 13); c.A() != 0 {
		t.Errorf("the marker stayed the way it was written: pixel (18,13) = %v, want nothing there", c)
	}
	if c := pixelAt(cv, 12, 18); c.A() != 0 {
		t.Errorf("the marker turned the wrong way round the corner: pixel (12,18) = %v, want nothing there", c)
	}
}

func TestRenderTurnsAnEndMarkerToTheWayThePathCameIn(t *testing.T) {
	// The last vertex of a path has no segment going out of it, so the only
	// direction there is the one the path arrived by — which is what points an
	// arrowhead at the end of a line out along the line. This one runs north
	// into its last point, so the marker is turned north: it is drawn above the
	// vertex rather than to the right of it, which is where it was written.
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs>
			<marker id="m" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="2" orient="auto">
				<rect x="0" y="0" width="4" height="4" fill="red"/>
			</marker>
		</defs>
		<line x1="5" y1="25" x2="5" y2="10" fill="none" stroke="none" marker-end="url(#m)"/>
	</svg>`, 40, 40)
	if c := pixelAt(cv, 4, 7); c.R() < 200 || c.G() > 60 || c.B() > 60 {
		t.Errorf("the marker did not turn to the way the path came in: pixel (4,7) = %v, want red", c)
	}
	if c := pixelAt(cv, 8, 9); c.A() != 0 {
		t.Errorf("the marker stayed the way it was written: pixel (8,9) = %v, want nothing there", c)
	}
}

func TestAutoStartReverseLeavesTheOtherMarkersAlone(t *testing.T) {
	// `auto-start-reverse` is `auto` everywhere but at the first vertex, where
	// it turns the marker around so that one marker can be put at both ends of
	// a line and point out of each of them. The one at the end keeps the
	// direction it has with plain `auto`, which is the direction the path is
	// going — east here — so it points out of its end rather than back along
	// the line.
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs>
			<marker id="m" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="2" orient="auto-start-reverse">
				<rect x="0" y="0" width="4" height="4" fill="red"/>
			</marker>
		</defs>
		<line x1="5" y1="15" x2="25" y2="15" fill="none" stroke="none"
			marker-start="url(#m)" marker-end="url(#m)"/>
	</svg>`, 40, 40)
	for _, at := range [][2]int{{3, 15}, {27, 15}} {
		if c := pixelAt(cv, at[0], at[1]); c.R() < 200 || c.G() > 60 || c.B() > 60 {
			t.Errorf("nothing at (%d,%d), where a marker should be: %v", at[0], at[1], c)
		}
	}
	for _, at := range [][2]int{{7, 15}, {22, 15}} {
		if c := pixelAt(cv, at[0], at[1]); c.A() != 0 {
			t.Errorf("a marker drew at (%d,%d) on the wrong side of its vertex: %v", at[0], at[1], c)
		}
	}
}

func TestRenderMarksEveryVertexOfAShapeThatCloses(t *testing.T) {
	// A shape that closes travels back to where it started, so its last vertex
	// is its first one again — the point the path arrives at when it closes —
	// and it is marked as the last one, with marker-end, while the four corners
	// between the two are marker-mid. The closing vertex is therefore marked
	// twice over, once as the first vertex and once as the last, and each time
	// turned to the way the path is going there: east where it leaves, north
	// where it arrives. The corners in between are turned to the middle of the
	// turn each of them makes, which is why the markers stand on the diagonal
	// out of every corner rather than square with the rectangle.
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<defs>
			<marker id="m" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="0" orient="auto">
				<rect x="0" y="0" width="4" height="4" fill="red"/>
			</marker>
		</defs>
		<rect x="0" y="10" width="10" height="10" fill="none" stroke="none"
			marker-start="url(#m)" marker-mid="url(#m)" marker-end="url(#m)"/>
	</svg>`, 40, 40)
	for _, at := range [][2]int{
		{1, 11},  // the first vertex, where the path leaves going east
		{11, 12}, // the first corner, where the path turns from east to south
		{7, 20},  // the second corner, from south to west
		{1, 17},  // the third, from west to north
		{1, 8},   // the last vertex again, where the path arrives going north
	} {
		if c := pixelAt(cv, at[0], at[1]); c.R() < 200 || c.G() > 60 || c.B() > 60 {
			t.Errorf("no marker at %v: pixel (%d,%d) = %v, want red", at, at[0], at[1], c)
		}
	}
	// Nothing between the corners: only the vertices are marked, and the shape
	// itself has no stroke and no fill to draw with.
	for _, at := range [][2]int{{6, 11}, {6, 15}} {
		if c := pixelAt(cv, at[0], at[1]); c.A() != 0 {
			t.Errorf("something was drawn at (%d,%d), between the corners", at[0], at[1])
		}
	}
}

func TestAMarkerSurvivesAPathThatTurnsBackOnItself(t *testing.T) {
	// Two vertices no other marker has an answer for. A path that turns
	// straight back on itself has no middle of the turn to stand in — the two
	// directions add up to nothing — and the spec's answer is the direction it
	// came in by, so the marker keeps facing the way the path was going rather
	// than turning round with it. A vertex written twice over is a turn of no
	// size at all, and its marker is put down the same way rather than coming
	// out sideways or not at all.
	_, back := painted(t, `<svg viewBox="0 0 40 40">
		<defs>
			<marker id="m" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="2" orient="auto">
				<rect x="0" y="0" width="4" height="4" fill="red"/>
			</marker>
		</defs>
		<path d="M5 15 L20 15 L5 15" fill="none" stroke="none" marker-mid="url(#m)"/>
	</svg>`, 40, 40)
	if c := pixelAt(back, 21, 14); c.A() == 0 {
		t.Errorf("the marker did not keep the direction the path came in by: pixel (21,14) = %v", c)
	}
	if c := pixelAt(back, 17, 15); c.A() != 0 {
		t.Errorf("the marker turned round with the path: pixel (17,15) = %v, want nothing there", c)
	}

	_, twice := painted(t, `<svg viewBox="0 0 40 40">
		<defs>
			<marker id="m" markerUnits="userSpaceOnUse" markerWidth="10" markerHeight="10" refX="0" refY="2" orient="auto">
				<rect x="0" y="0" width="4" height="4" fill="red"/>
			</marker>
		</defs>
		<path d="M5 15 L20 15 L20 15 L5 30" fill="none" stroke="none" marker-mid="url(#m)"/>
	</svg>`, 40, 40)
	// The first of the two coinciding vertices has a segment of no length
	// leaving it, so the only direction there is the one it came in by.
	if c := pixelAt(twice, 21, 14); c.A() == 0 {
		t.Errorf("no marker at the vertex with no segment leaving it: pixel (21,14) = %v", c)
	}
	// The second has one arriving of no length, and turns to the segment that
	// leaves it instead — north west, back the way the path goes.
	if c := pixelAt(twice, 17, 16); c.A() == 0 {
		t.Errorf("no marker at the vertex with no segment arriving at it: pixel (17,16) = %v", c)
	}
}
