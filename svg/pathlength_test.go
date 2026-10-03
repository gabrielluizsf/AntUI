package svg

import "testing"

// The length a drawing declares for a path is what the stroke of that shape is
// measured out against: the dashes, the gaps and the offset into them are
// stretched by however far the path measures over it, so a pattern written in
// the drawing's own numbers comes out in those numbers wherever it is painted.

func TestPathLengthRescalesTheDashOfTheShapeItIsOn(t *testing.T) {
	// A line a hundred long that says it is twenty has every length of its
	// pattern stretched five times: the dashes, the gaps, and where along them
	// the stroke starts, all measured against what was declared.
	img, err := Parse(`<svg viewBox="0 0 100 100">
		<line x1="0" y1="0" x2="100" y2="0" fill="none" stroke="red"
			stroke-dasharray="4 2" stroke-dashoffset="3" pathLength="20"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if ws := img.Warnings(); len(ws) != 0 {
		t.Fatalf("the drawing said %v, want nothing at all", ws)
	}
	line := img.Root.find("line")
	if line == nil || line.Style.Dash == nil {
		t.Fatal("the line and the pattern on it are not both in the drawing")
	}
	d := line.Style.Dash
	if len(d.On) != 1 || d.On[0] != 20 || len(d.Off) != 1 || d.Off[0] != 10 {
		t.Errorf("the pattern is %v on and %v off, want 20 and 10", d.On, d.Off)
	}
	if d.Offset != 15 {
		t.Errorf("the offset is %v, want 15", d.Offset)
	}
}

func TestPathLengthIsNotInherited(t *testing.T) {
	// The length is of the shape it is written on: what a group declares is a
	// length of the group's own shape, and the shapes inside are measured
	// against nothing of their own unless they write a length too. A length on
	// the shape beats one on the group above it, as an attribute of its own
	// beats the one it was given.
	for _, tc := range []struct {
		name  string
		group string
		inner string
		want  [2]float64
	}{
		{"the group's length is not the shape's", "20", "", [2]float64{4, 2}},
		{"the shape's own length is what counts", "100", ` pathLength="20"`, [2]float64{20, 10}},
	} {
		img, err := Parse(`<svg viewBox="0 0 100 100">
			<g pathLength="` + tc.group + `" stroke-dasharray="4 2">
				<line x1="0" y1="0" x2="100" y2="0" fill="none" stroke="red"` + tc.inner + `/>
			</g>
		</svg>`)
		if err != nil {
			t.Fatalf("%s: parse: %v", tc.name, err)
		}
		line := img.Root.find("line")
		if line == nil || line.Style.Dash == nil {
			t.Fatalf("%s: the line and the pattern on it are not both in the drawing", tc.name)
		}
		if d := line.Style.Dash; d.On[0] != tc.want[0] || d.Off[0] != tc.want[1] {
			t.Errorf("%s: the pattern is %v on and %v off, want %v and %v",
				tc.name, d.On, d.Off, tc.want[0], tc.want[1])
		}
	}
}

func TestPathLengthMeasuresAClosedShapeAllTheWayRound(t *testing.T) {
	// A square of ten by ten goes forty round, the side that leads back to the
	// first point included, so a pattern declared against twenty is stretched
	// twice — and were the closing side left out of the measure it would come
	// out at thirty and the pattern at one and a half times.
	img, err := Parse(`<svg viewBox="0 0 100 100">
		<rect x="0" y="0" width="10" height="10" fill="none" stroke="red"
			stroke-dasharray="4 2" pathLength="20"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rect := img.Root.find("rect")
	if rect == nil || rect.Style.Dash == nil {
		t.Fatal("the rectangle and the pattern on it are not both in the drawing")
	}
	if d := rect.Style.Dash; d.On[0] != 8 || d.Off[0] != 4 {
		t.Errorf("the pattern is %v on and %v off, want 8 and 4", d.On, d.Off)
	}
}

func TestPathLengthThatIsNotAPositiveNumberSaysSo(t *testing.T) {
	// A length this cannot read is said rather than quietly taken as nothing,
	// because the drawing asked for the pattern to be measured against it; the
	// pattern then stays as it was written. An empty value says nothing at all,
	// which is what an attribute written with nothing in it means everywhere.
	for _, tc := range []struct {
		value string
		warns int
		want  [2]float64
	}{
		{"", 0, [2]float64{4, 2}},     // nothing said, nothing done
		{"20", 0, [2]float64{20, 10}}, // a line of a hundred against twenty is five of it
		{"abc", 1, [2]float64{4, 2}},
		{"0", 1, [2]float64{4, 2}},
		{"-4", 1, [2]float64{4, 2}},
		{"10px", 1, [2]float64{4, 2}},
	} {
		img, err := Parse(`<svg viewBox="0 0 100 100">
			<line x1="0" y1="0" x2="100" y2="0" fill="none" stroke="red"
				stroke-dasharray="4 2" pathLength="` + tc.value + `"/>
		</svg>`)
		if err != nil {
			t.Fatalf("pathLength=%q: parse: %v", tc.value, err)
		}
		if got := len(img.Warnings()); got != tc.warns {
			t.Errorf("pathLength=%q: %d warnings, want %d (%v)",
				tc.value, got, tc.warns, img.Warnings())
		}
		line := img.Root.find("line")
		if line == nil || line.Style.Dash == nil {
			t.Fatalf("pathLength=%q: the line and the pattern on it are not both there", tc.value)
		}
		if d := line.Style.Dash; d.On[0] != tc.want[0] || d.Off[0] != tc.want[1] {
			t.Errorf("pathLength=%q: the pattern is %v on and %v off, want %v and %v",
				tc.value, d.On, d.Off, tc.want[0], tc.want[1])
		}
	}
}

func TestPathLengthWithNoPatternToScale(t *testing.T) {
	// A length declared where there are no dashes to measure out is read and
	// has nothing to do: the shape paints as it would have anyway, and a
	// drawing that asked for nothing it could not do says nothing.
	img, err := Parse(`<svg viewBox="0 0 100 100">
		<rect x="10" y="10" width="40" height="40" fill="none" stroke="red" pathLength="20"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if ws := img.Warnings(); len(ws) != 0 {
		t.Fatalf("the drawing said %v, want nothing at all", ws)
	}
	rect := img.Root.find("rect")
	if rect == nil {
		t.Fatal("the rectangle is not in the drawing")
	}
	if rect.Style.Dash != nil {
		t.Error("the length made a pattern where there was none")
	}
	if cv := img.Render(100, 100); cv == nil {
		t.Error("nothing came back to paint on")
	}
}

func TestRenderScalesTheDashToTheDeclaredLength(t *testing.T) {
	// The stretch of the pattern reaches the pixels: the dashes of a line are
	// walked from where the shape begins in the drawing's own units, and the
	// length the drawing declared for the path is what they are measured out
	// against — so declared ten long, a line of a hundred has every dash and
	// every gap of its pattern stretched ten times over and comes out painted
	// from end to end, while a pattern declared by nothing keeps the ten units
	// it was written with and its gaps fall where those units run out.
	const head = `<svg viewBox="0 0 100 100"><line x1="0" y1="50" x2="100" y2="50" ` +
		`fill="none" stroke="red" stroke-width="4" stroke-dasharray="10 10"`
	for _, tc := range []struct {
		name  string
		extra string
		want  map[int]bool
	}{
		{"declared ten long", ` pathLength="10"`, map[int]bool{5: true, 15: true, 52: true}},
		{"declared by nothing", "", map[int]bool{5: true, 15: false, 52: false}},
	} {
		img, err := Parse(head + tc.extra + `/></svg>`)
		if err != nil {
			t.Fatalf("%s: parse: %v", tc.name, err)
		}
		cv := img.Render(100, 100)
		if cv == nil {
			t.Fatalf("%s: nothing came back to paint on", tc.name)
		}
		for x, painted := range tc.want {
			if got := pixelAt(cv, x, 50).A() != 0; got != painted {
				t.Errorf("%s: the pixel at %d is painted=%v, want %v", tc.name, x, got, painted)
			}
		}
	}
}

func TestRenderScalesTheDashOffsetToTheDeclaredLength(t *testing.T) {
	// The offset is stretched along with the rest of the pattern, so where a
	// stroke begins its first dash moves with it: a line that starts a whole
	// dash in still starts a whole dash in, only ten times further along once
	// it is declared ten long. The gap the offset opens the line with reaches
	// ten times as far, and the pixels below turn over the other way round to
	// match it.
	const head = `<svg viewBox="0 0 100 100"><line x1="0" y1="50" x2="100" y2="50" ` +
		`fill="none" stroke="red" stroke-width="4" stroke-dasharray="5 5" stroke-dashoffset="5"`
	for _, tc := range []struct {
		name  string
		extra string
		want  map[int]bool
	}{
		{"declared ten long", ` pathLength="10"`, map[int]bool{3: false, 30: false, 60: true}},
		{"declared by nothing", "", map[int]bool{3: false, 8: true, 30: false, 38: true}},
	} {
		img, err := Parse(head + tc.extra + `/></svg>`)
		if err != nil {
			t.Fatalf("%s: parse: %v", tc.name, err)
		}
		cv := img.Render(100, 100)
		if cv == nil {
			t.Fatalf("%s: nothing came back to paint on", tc.name)
		}
		for x, painted := range tc.want {
			if got := pixelAt(cv, x, 50).A() != 0; got != painted {
				t.Errorf("%s: the pixel at %d is painted=%v, want %v", tc.name, x, got, painted)
			}
		}
	}
}
