package svg

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// What a shape looks like is written on the element itself, or above it for
// everything inside, or in a `style` attribute that says the same things in
// another way. These check each of the three, and that the one above is what an
// element with nothing of its own uses.

func TestStyleReadsAPaint(t *testing.T) {
	st := Style{}.with(mustElement(t, `<rect fill="red" stroke="blue"/>`), quiet)
	if !st.HasFill || st.Fill != canvas.RGB(255, 0, 0) {
		t.Errorf("fill = %v, %v, want red and to have one", st.Fill, st.HasFill)
	}
	if !st.HasStroke || st.Stroke != canvas.RGB(0, 0, 255) {
		t.Errorf("stroke = %v, %v, want blue and to have one", st.Stroke, st.HasStroke)
	}
}

func TestStylePaintOfNone(t *testing.T) {
	// `none` is not a colour but the absence of one, and it turns a fill that
	// was there off rather than painting it black.
	inherited := Style{}.with(mustElement(t, `<g fill="red"/>`), quiet)
	if !inherited.HasFill {
		t.Fatal("the group has no fill to turn off")
	}
	st := inherited.with(mustElement(t, `<rect fill="none"/>`), quiet)
	if st.HasFill {
		t.Error("fill=none still leaves a fill on it")
	}
	// A stroke that was there is turned off the same way.
	withStroke := Style{}.with(mustElement(t, `<g stroke="red" stroke-width="4"/>`), quiet)
	st = withStroke.with(mustElement(t, `<rect stroke="none"/>`), quiet)
	if st.HasStroke {
		t.Error("stroke=none still leaves a stroke on it")
	}
	if st.Width != 4 {
		t.Errorf("stroke-width = %v, want the inherited 4 left alone", st.Width)
	}
}

func TestStyleIsInheritedAndOverridden(t *testing.T) {
	// A group colours everything inside it, and one shape inside may be a colour
	// of its own without disturbing the rest.
	group := Style{}.with(mustElement(t, `<g fill="red" stroke="blue" stroke-width="3"/>`), quiet)
	one := group.with(mustElement(t, `<rect/>`), quiet)
	if one.Fill != canvas.RGB(255, 0, 0) || one.Stroke != canvas.RGB(0, 0, 255) {
		t.Errorf("a shape with nothing of its own got %v on %v, want the group's",
			one.Fill, one.Stroke)
	}
	if one.Width != 3 {
		t.Errorf("stroke-width = %v, want the group's 3", one.Width)
	}
	own := group.with(mustElement(t, `<rect fill="green"/>`), quiet)
	if own.Fill != canvas.RGB(0, 128, 0) {
		t.Errorf("fill = %v, want the green it was given", own.Fill)
	}
	if own.Stroke != canvas.RGB(0, 0, 255) {
		t.Errorf("stroke = %v, want the group's blue still", own.Stroke)
	}
}

func TestStyleAttributeBeatsTheAttributeOfTheSameName(t *testing.T) {
	// A presentation attribute says the same thing as a declaration in `style`,
	// and where both are written the declaration is the one that counts.
	st := Style{}.with(mustElement(t, `<rect fill="red" style="fill: blue; stroke-width: 5"/>`), quiet)
	if st.Fill != canvas.RGB(0, 0, 255) {
		t.Errorf("fill = %v, want the blue from the style attribute", st.Fill)
	}
	if st.Width != 5 {
		t.Errorf("stroke-width = %v, want the 5 from the style attribute", st.Width)
	}
}

func TestStyleDeclarations(t *testing.T) {
	// Declarations may be separated by semicolons with any spacing at all, and
	// the last one of a name is the one that counts.
	st := Style{}.with(mustElement(t, `<rect style=" fill : red ; fill : blue ; stroke : green " style2="x"/>`), quiet)
	if st.Fill != canvas.RGB(0, 0, 255) {
		t.Errorf("fill = %v, want the last blue declared", st.Fill)
	}
	if st.Stroke != canvas.RGB(0, 128, 0) {
		t.Errorf("stroke = %v, want the green", st.Stroke)
	}
	// A declaration that is not one is left out rather than guessed at.
	st = Style{}.with(mustElement(t, `<rect style="; ; nonsense; : red; fill:;"/>`), quiet)
	if st.HasFill {
		t.Errorf("fill = %v, want the empty declaration to be left out", st.Fill)
	}
}

func TestStyleLengths(t *testing.T) {
	// The units a length may be written in are turned into pixels, and a number
	// with no unit is the drawing's own.
	for _, tc := range []struct {
		in   string
		want float64
	}{
		{"10", 10},
		{"10px", 10},
		{"1in", 96},
		{"12pt", 16},
		{"1pc", 16},
		{"2.54cm", 96},
		{"25.4mm", 96},
		{"-4", -4},
		{"1e2", 100},
	} {
		got, ok := parseLength(tc.in)
		if !ok {
			t.Errorf("parseLength(%q) did not read it", tc.in)
			continue
		}
		if !nearFloat(got, tc.want, 0.001) {
			t.Errorf("parseLength(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
	for _, s := range []string{"", "  ", "10em", "10 px", "abc", "px", "10%" + "0"} {
		if v, ok := parseLength(s); ok {
			t.Errorf("parseLength(%q) = %v, want it not read", s, v)
		}
	}
}

func TestStyleOpacities(t *testing.T) {
	// An opacity is a number or a percentage, and one outside the range from
	// nothing to opaque is held at an end of it rather than turned round.
	for _, tc := range []struct {
		in   string
		want float64
	}{
		{"0", 0}, {"1", 1}, {"0.5", 0.5}, {"50%", 0.5},
		{"2", 1}, {"-1", 0}, {"200%", 1},
	} {
		got, ok := parseAlpha(tc.in)
		if !ok {
			t.Errorf("parseAlpha(%q) did not read it", tc.in)
			continue
		}
		if !nearFloat(got, tc.want, 0.001) {
			t.Errorf("parseAlpha(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
	for _, s := range []string{"", "   ", "half", "50 percent"} {
		if v, ok := parseAlpha(s); ok {
			t.Errorf("parseAlpha(%q) = %v, want it not read", s, v)
		}
	}
}

func TestStyleCapsJoinsAndLimit(t *testing.T) {
	st := Style{}.with(mustElement(t, `<rect stroke-linecap="round" stroke-linejoin="bevel" stroke-miterlimit="2.5"/>`), quiet)
	if st.Cap != canvas.CapRound {
		t.Errorf("cap = %v, want round", st.Cap)
	}
	if st.Join != canvas.JoinBevel {
		t.Errorf("join = %v, want bevel", st.Join)
	}
	if !nearFloat(st.MiterLimit, 2.5, 0.001) {
		t.Errorf("miter limit = %v, want 2.5", st.MiterLimit)
	}
	// A name it does not know falls back to the one SVG uses by default rather
	// than to nothing at all, and a limit that is not a positive number leaves
	// the inherited one alone.
	inherited := Style{MiterLimit: 4}
	st = inherited.with(mustElement(t, `<rect stroke-linecap="nonsense" stroke-linejoin="nonsense" stroke-miterlimit="0"/>`), quiet)
	if st.Cap != canvas.CapButt {
		t.Errorf("cap = %v, want butt", st.Cap)
	}
	if st.Join != canvas.JoinMiter {
		t.Errorf("join = %v, want miter", st.Join)
	}
	if st.MiterLimit != 4 {
		t.Errorf("miter limit = %v, want the inherited 4 kept", st.MiterLimit)
	}
}

func TestStyleHiding(t *testing.T) {
	// A shape that is not to be drawn says so with display or with visibility,
	// and one inside a hidden group stays hidden however it is written.
	for _, name := range []string{"display", "visibility"} {
		st := styled(mustElement(t, `<rect `+name+`="none"/>`), quiet)
		if !st.Hidden {
			t.Errorf("%s=none did not hide it", name)
		}
		shown := styled(mustElement(t, `<rect `+name+`="inline"/>`), quiet)
		if shown.Hidden {
			t.Errorf("%s=inline hid it", name)
		}
		group := Style{}.with(mustElement(t, `<g display="none"/>`), quiet)
		if !group.with(mustElement(t, `<rect display="inline"/>`), quiet).Hidden {
			t.Errorf("%s=inline inside a hidden group showed it again", name)
		}
	}
}

func TestStyleDashes(t *testing.T) {
	// A pattern of dashes and gaps, and a gap that is not there is no pattern.
	st := Style{}.with(mustElement(t, `<rect stroke-dasharray="4 2"/>`), quiet)
	if st.Dash == nil {
		t.Fatal("the pattern was not read")
	}
	st = Style{}.with(mustElement(t, `<rect stroke-dasharray="4,2"/>`), quiet)
	if st.Dash == nil {
		t.Error("a pattern written with commas was not read")
	}
	for _, s := range []string{`stroke-dasharray="none"`, `stroke-dasharray=""`, `stroke-dasharray="abc"`} {
		got := styled(mustElement(t, `<rect `+s+`/>`), quiet)
		if got.Dash != nil {
			t.Errorf("%s gave a pattern of %v, want none", s, got.Dash)
		}
	}
	// An offset on its own says nothing without a pattern to shift along.
	bare := Style{}.with(mustElement(t, `<rect stroke-dashoffset="3"/>`), quiet)
	if bare.Dash != nil {
		t.Error("an offset on its own made a pattern")
	}
}

func TestStyleDashOffsetReadsWithThePattern(t *testing.T) {
	// The offset shifts the pattern along the line, and it only means anything
	// once there is a pattern. This is the same whether the two are written as
	// attributes or as declarations, and in either order they are found in.
	for _, src := range []string{
		`<rect stroke-dasharray="4 2" stroke-dashoffset="3"/>`,
		`<rect stroke-dashoffset="3" stroke-dasharray="4 2"/>`,
		`<rect style="stroke-dashoffset: 3; stroke-dasharray: 4 2"/>`,
		`<rect style="stroke-dasharray: 4 2; stroke-dashoffset: 3"/>`,
	} {
		for range 20 {
			st := styled(mustElement(t, src), quiet)
			if st.Dash == nil {
				t.Fatalf("%s: no pattern was read", src)
			}
			if st.Dash.Offset != 3 {
				t.Fatalf("%s: offset = %v, want 3", src, st.Dash.Offset)
			}
		}
	}
}

func TestStyleTransformsComposeOutwards(t *testing.T) {
	// A transform on a child happens to the shape before the one above it, so
	// the two compose the way a list of them does: the one above is the outer
	// one, and a point goes through the child's first.
	group := Style{}.with(mustElement(t, `<g transform="translate(10, 0)"/>`), quiet)
	child := group.with(mustElement(t, `<rect transform="scale(2)"/>`), quiet)
	if !child.HasTransform {
		t.Fatal("the child's own transform was lost")
	}
	// A unit point is scaled by the child to two and then moved by the group.
	if x := mapX(child.Transform, 1, 0); !nearFloat(x, 12, 0.001) {
		t.Errorf("(1,0) went to x=%v, want 12", x)
	}
	// A transform that cannot be read leaves the inherited one alone.
	broken := group.with(mustElement(t, `<rect transform="nonsense"/>`), quiet)
	if !nearMatrix(broken.Transform, group.Transform) {
		t.Error("an unreadable transform threw the inherited one away")
	}
}

// styled is the style one element ends up with, and the warnings it makes are
// counted by the tests that care what was said rather than ending the test
// where they are found.
func styled(e *element, warn func(string, ...any)) Style { return Style{}.with(e, warn) }
