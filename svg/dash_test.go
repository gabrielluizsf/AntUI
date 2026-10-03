package svg

import "testing"

// The lengths of a stroke-dasharray are lengths of the drawing: ten on and ten
// off is ten units of the coordinates the file wrote, and what says how many
// pixels those ten come out at is the viewBox and the transform the shape goes
// through, the same as for everything else in the picture. The pattern is cut
// where the shape was written and the cut travels with the shape from there.

func TestRenderScalesTheDashWithTheViewBox(t *testing.T) {
	// A drawing of a hundred units and one of fifty both ask for the same ten
	// on and ten off, and both are painted into a hundred pixels: a unit of the
	// first is one pixel and a unit of the second is two, so the gaps of the
	// second fall twice as far along. Where they fall is what the pixels below
	// say, and the same numbers in both drawings would put the gaps in the same
	// place — which is what happens when the pattern is cut in pixels instead
	// of in the units the drawing wrote.
	for _, tc := range []struct {
		name string
		draw string
		row  int
		want map[int]bool
	}{
		{
			"a hundred units into a hundred pixels",
			`<svg viewBox="0 0 100 100">
				<line x1="0" y1="50" x2="100" y2="50" fill="none" stroke="red"
					stroke-width="4" stroke-dasharray="10 10"/>
			</svg>`,
			50,
			map[int]bool{5: true, 12: false, 45: true, 95: false},
		},
		{
			"fifty units into a hundred pixels",
			`<svg viewBox="0 0 50 50">
				<line x1="0" y1="25" x2="50" y2="25" fill="none" stroke="red"
					stroke-width="4" stroke-dasharray="10 10"/>
			</svg>`,
			50,
			map[int]bool{5: true, 30: false, 45: true, 70: false},
		},
	} {
		img, err := Parse(tc.draw)
		if err != nil {
			t.Fatalf("%s: parse: %v", tc.name, err)
		}
		if ws := img.Warnings(); len(ws) != 0 {
			t.Fatalf("%s: the drawing said %v, want nothing at all", tc.name, ws)
		}
		cv := img.Render(100, 100)
		if cv == nil {
			t.Fatalf("%s: nothing came back to paint on", tc.name)
		}
		for x, painted := range tc.want {
			if got := pixelAt(cv, x, tc.row).A() != 0; got != painted {
				t.Errorf("%s: the pixel at %d is painted=%v, want %v", tc.name, x, got, painted)
			}
		}
	}
}

func TestRenderScalesTheDashWithTheTransform(t *testing.T) {
	// The transform of the element is put on the shape while it is built and
	// the pattern is cut before that, so what the transform does to the shape
	// it does to the cut as well: a group scaled twice draws dashes twice as
	// long on the screen, one scaled half as much draws them half as long, and
	// one that is not scaled at all draws the ten units as ten pixels. A
	// pattern cut after the transform would come out the same in all three.
	for _, tc := range []struct {
		name  string
		group string
		row   int
		want  map[int]bool
	}{
		{
			"where it was written",
			"",
			25,
			map[int]bool{5: true, 15: false, 25: true, 35: false},
		},
		{
			"scaled twice",
			` transform="scale(2)"`,
			50,
			map[int]bool{5: true, 15: true, 25: false, 45: true, 75: false},
		},
		{
			"scaled half",
			` transform="scale(0.5)"`,
			12,
			map[int]bool{2: true, 7: false, 12: true, 17: false},
		},
	} {
		img, err := Parse(`<svg viewBox="0 0 100 100">
			<g` + tc.group + `>
				<line x1="0" y1="25" x2="40" y2="25" fill="none" stroke="red"
					stroke-width="4" stroke-dasharray="10 10"/>
			</g>
		</svg>`)
		if err != nil {
			t.Fatalf("%s: parse: %v", tc.name, err)
		}
		cv := img.Render(100, 100)
		if cv == nil {
			t.Fatalf("%s: nothing came back to paint on", tc.name)
		}
		for x, painted := range tc.want {
			if got := pixelAt(cv, x, tc.row).A() != 0; got != painted {
				t.Errorf("%s: the pixel at %d is painted=%v, want %v", tc.name, x, got, painted)
			}
		}
	}
}

func TestRenderCutsTheDashBeforeANonUniformTransform(t *testing.T) {
	// The pattern is cut in the units the drawing wrote and only then taken
	// through the transform, so a stretch of two to one doubles what a dash
	// measures along the stretch — ten units of it become twenty pixels —
	// while the shape is cut at ten units no matter which way each part of it
	// ends up running. Cutting after the transform, or cutting by the one
	// scale the whole matrix comes to, would put a dash that runs the length
	// of the stretch somewhere in between the two, which is nowhere the
	// drawing asked for.
	img, err := Parse(`<svg viewBox="0 0 100 100">
		<g transform="scale(2,1)">
			<line x1="0" y1="25" x2="40" y2="25" fill="none" stroke="red"
				stroke-width="4" stroke-dasharray="10 10"/>
		</g>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if ws := img.Warnings(); len(ws) != 0 {
		t.Fatalf("the drawing said %v, want nothing at all", ws)
	}
	cv := img.Render(100, 100)
	if cv == nil {
		t.Fatal("nothing came back to paint on")
	}
	// The dash runs from ten units to twenty along the stretched line, and
	// every pixel of it is painted; the gap that follows runs to forty.
	for _, tc := range []struct {
		x       int
		painted bool
	}{
		{5, true},
		{15, true},
		{25, false},
		{45, true},
		{75, false},
	} {
		if got := pixelAt(cv, tc.x, 25).A() != 0; got != tc.painted {
			t.Errorf("the pixel at %d is painted=%v, want %v", tc.x, got, tc.painted)
		}
	}
}

func TestRenderDashesTheStrokeAndNotTheFill(t *testing.T) {
	// The cut is of the outline the stroke is widened from and not of the
	// shape itself: what is inside the rectangle is filled whole, while the
	// stroke over its edge has the gaps the pattern asks for. Cutting the
	// shape would fill the dashes and leave the rest of it empty.
	img, err := Parse(`<svg viewBox="0 0 100 100">
		<rect x="10" y="10" width="40" height="40" fill="red" stroke="red"
			stroke-width="4" stroke-dasharray="10 10"/>
	</svg>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	cv := img.Render(100, 100)
	if cv == nil {
		t.Fatal("nothing came back to paint on")
	}
	if got := pixelAt(cv, 30, 30); got.A() == 0 {
		t.Errorf("the middle of the rectangle is %v, want it filled", got)
	}
	// The stroke reaches two units past the edge, and the row above the top
	// edge is that and nothing else: the pattern starts at the corner of the
	// shape, so ten in and ten out from there.
	for _, tc := range []struct {
		x       int
		painted bool
	}{
		{15, true},
		{25, false},
	} {
		if got := pixelAt(cv, tc.x, 8).A() != 0; got != tc.painted {
			t.Errorf("the pixel at %d is painted=%v, want %v", tc.x, got, tc.painted)
		}
	}
}
