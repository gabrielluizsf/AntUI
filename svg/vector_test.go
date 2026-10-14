package svg

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// The stroke of a shape is part of the shape, so everything that scales the
// shape scales the stroke with it — the viewBox onto the canvas and the
// transforms on the way — which is what SVG does and what a drawing gets unless
// it asks otherwise. `vector-effect="non-scaling-stroke"` is the otherwise: the
// stroke comes out at the width it says, in the pixels the drawing comes out at,
// whatever scaled the shape. `none` is the scaled stroke again, and it is what
// every drawing has until one of the two is written.

// paintedRows is how many rows of one column the drawing painted on, which is
// how wide a horizontal stroke came out.
func paintedRows(cv *canvas.Canvas, x, from, to int) int {
	rows := 0
	for y := from; y <= to; y++ {
		if pixelAt(cv, x, y).A() != 0 {
			rows++
		}
	}
	return rows
}

func TestVectorEffectIsRead(t *testing.T) {
	// The values the spec names are read from the attribute and from a `style`
	// declaration alike, alone and in combination, and a value this package
	// does not know is said rather than quietly painted as though it had not
	// been asked for. An attribute written with nothing in it means nothing at
	// all, as everywhere else, and `none` is the empty value.
	for _, tc := range []struct {
		name  string
		value string
		warns int
		want  VectorEffect
	}{
		{"nothing written", "", 0, VectorEffect{}},
		{"the stroke a drawing gets by default", "none", 0, VectorEffect{}},
		{"the stroke that does not scale", "non-scaling-stroke", 0, VectorEffect{NonScalingStroke: true}},
		{"the keyword in another case", "NON-SCALING-STROKE", 0, VectorEffect{NonScalingStroke: true}},
		{"the size that does not scale", "non-scaling-size", 0, VectorEffect{NonScalingStroke: true, NonScalingSize: true}},
		{"the turn that does not happen", "non-rotation", 0, VectorEffect{NonRotation: true}},
		{"the place that does not move", "fixed-position", 0, VectorEffect{FixedPosition: true}},
		{"two of them together", "non-scaling-size non-rotation", 0, VectorEffect{NonScalingStroke: true, NonScalingSize: true, NonRotation: true}},
		{"three of them together", "fixed-position, non-rotation non-scaling-size", 0, VectorEffect{NonScalingStroke: true, NonScalingSize: true, NonRotation: true, FixedPosition: true}},
		{"a value this package does not draw", "scaled", 1, VectorEffect{}},
		{"a value it does not draw beside one it does", "non-rotation skewed", 1, VectorEffect{NonRotation: true}},
	} {
		img, err := Parse(`<svg viewBox="0 0 10 10">
			<line x1="0" y1="5" x2="10" y2="5" stroke="red" vector-effect="` + tc.value + `"/>
		</svg>`)
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
		if got := line.Style.VectorEffect; got != tc.want {
			t.Errorf("%s: vector-effect = %+v, want %+v", tc.name, got, tc.want)
		}
	}
}

func TestVectorEffectIsReadFromTheStyleAttribute(t *testing.T) {
	// A declaration says the same thing the attribute does, and the one written
	// in `style` is the one that counts where both are written.
	st := Style{}.with(mustElement(t,
		`<line vector-effect="none" style="vector-effect: non-scaling-stroke"/>`), quiet)
	if !st.VectorEffect.NonScalingStroke {
		t.Error("the declaration did not beat the attribute of the same name")
	}
	st = Style{}.with(mustElement(t,
		`<line vector-effect="non-scaling-stroke" style="vector-effect: none"/>`), quiet)
	if st.VectorEffect.NonScalingStroke {
		t.Error("the declaration did not beat the attribute of the same name")
	}
}

func TestVectorEffectIsNotInherited(t *testing.T) {
	// The spec has vector-effect as not inherited, so what a group writes is
	// about the group, and the shapes inside start from the ordinary behaviour
	// again. A stroke that does not scale on the group is about the group's own
	// outline — which it has none of — and the shape inside is scaled anyway.
	group := Style{}.with(mustElement(t, `<g vector-effect="non-scaling-stroke"/>`), quiet)
	if !group.VectorEffect.NonScalingStroke {
		t.Fatal("the group did not take the vector-effect it was given")
	}
	inner := group.with(mustElement(t, `<line/>`), quiet)
	if inner.VectorEffect != (VectorEffect{}) {
		t.Error("the vector-effect travelled down to the shape, want it to stay on the group")
	}
	own := inner.with(mustElement(t, `<line vector-effect="non-scaling-stroke"/>`), quiet)
	if !own.VectorEffect.NonScalingStroke {
		t.Error("a shape's own vector-effect was not read")
	}
}

func TestRenderScalesTheStrokeWithTheShape(t *testing.T) {
	// The transform goes into the shape while it is built and the width the
	// stroke is widened by carries the same transform, so the group scaled
	// twice comes out a shape twice as big and a stroke twice as wide, while
	// the line beside it that nothing scaled keeps the four units it was
	// written with. This is the default — the stroke a drawing gets unless it
	// asks for one that does not scale.
	_, cv := painted(t, `<svg viewBox="0 0 40 40">
		<g transform="scale(2)">
			<line x1="0" y1="10" x2="20" y2="10" stroke="red" stroke-width="4"/>
		</g>
		<line x1="0" y1="30" x2="20" y2="30" stroke="red" stroke-width="4"/>
	</svg>`, 40, 40)
	if got := paintedRows(cv, 10, 14, 26); got < 7 || got > 9 {
		t.Errorf("the stroke under scale(2) is %d pixels wide, want 8", got)
	}
	if got := paintedRows(cv, 10, 26, 36); got < 3 || got > 5 {
		t.Errorf("the stroke nothing scaled is %d pixels wide, want 4", got)
	}
}

func TestRenderKeepsANonScalingStrokeAtTheWidthItSays(t *testing.T) {
	// A drawing of twenty units painted at two hundred scales everything ten
	// times over, and the stroke goes with it: a line one unit thick comes out
	// ten pixels wide. A stroke that does not scale stays one pixel wide
	// instead — not ten, and not one unit of the canvas it is painted on
	// either, which for this drawing would be the same ten pixels.
	for _, tc := range []struct {
		name  string
		extra string
		want  int
	}{
		{"the scaled stroke the drawing gets by default", "", 10},
		{"the stroke that does not scale", ` vector-effect="non-scaling-stroke"`, 2},
	} {
		_, cv := painted(t, `<svg viewBox="0 0 20 20">
			<line x1="0" y1="10" x2="20" y2="10" stroke="red" stroke-width="1"`+tc.extra+`/>
		</svg>`, 200, 200)
		got := paintedRows(cv, 100, 88, 112)
		if got < tc.want-1 || got > tc.want+1 {
			t.Errorf("%s: the stroke is %d pixels wide, want %d", tc.name, got, tc.want)
		}
	}
}

func TestANonScalingStrokeIgnoresTheTransformAndIsNotInherited(t *testing.T) {
	// The stroke that does not scale takes neither the transform of the shape
	// nor the viewBox onto the canvas: a line a unit thick inside a group
	// scaled twice is ten pixels wide as any other stroke would be, and two
	// pixels when it asks for one that does not scale. And vector-effect is not
	// inherited, so writing it on the group says it about the group's own
	// outline — which it has none of — and the shape inside is scaled anyway.
	for _, tc := range []struct {
		name  string
		group string
		line  string
		want  int
	}{
		{"the stroke the drawing gets by default", "", "", 20},
		{"the stroke written on the shape", "", ` vector-effect="non-scaling-stroke"`, 2},
		{"the stroke written on the group above it", ` vector-effect="non-scaling-stroke"`, "", 20},
	} {
		img, cv := painted(t, `<svg viewBox="0 0 20 20">
			<g transform="scale(2)"`+tc.group+`>
				<line x1="0" y1="5" x2="10" y2="5" stroke="red" stroke-width="1"`+tc.line+`/>
			</g>
		</svg>`, 200, 200)
		if ws := img.Warnings(); len(ws) != 0 {
			t.Fatalf("%s: the drawing said %v, want nothing at all", tc.name, ws)
		}
		got := paintedRows(cv, 50, 85, 115)
		if got < tc.want-1 || got > tc.want+1 {
			t.Errorf("%s: the stroke is %d pixels wide, want %d", tc.name, got, tc.want)
		}
	}
}

func TestRenderScalesTheStrokeOfASymbolToTheSizeItIsUsedAt(t *testing.T) {
	// The viewBox of a symbol is fitted to the size the `<use>` asks for, and
	// that fit is a transform like any other: the same symbol used twice at two
	// sizes comes out with a stroke of two units the first time and four the
	// second, because everything in it grew with the box it was put into.
	_, cv := painted(t, `<svg viewBox="0 0 40 20">
		<defs>
			<symbol id="s" viewBox="0 0 10 10">
				<line x1="0" y1="5" x2="10" y2="5" stroke="red" stroke-width="2"/>
			</symbol>
		</defs>
		<use href="#s" width="10" height="10"/>
		<use href="#s" x="20" width="20" height="20"/>
	</svg>`, 40, 20)
	if got := paintedRows(cv, 5, 0, 19); got < 1 || got > 3 {
		t.Errorf("the symbol used at its own size came out %d pixels wide, want 2", got)
	}
	if got := paintedRows(cv, 25, 0, 19); got < 3 || got > 5 {
		t.Errorf("the symbol used at twice its size came out %d pixels wide, want 4", got)
	}
}

// boxSpan is how many pixels wide and tall a box of painted pixels came out,
// which is how a test says how big a shape was held at.
func boxSpan(b box) (int, int) {
	return int(b.MaxX-b.MinX) + 1, int(b.MaxY-b.MinY) + 1
}

func TestRenderHoldsANonScalingSizeAgainstTheScaleAndTheViewBox(t *testing.T) {
	// The host of the effect is the canvas the drawing comes out on, so the
	// viewBox fit is part of what holds the size still: a rectangle of ten
	// units inside a group scaled twice comes out ten pixels at forty, where
	// the viewBox maps unit for unit, and ten pixels at eighty, where it maps
	// two. Without the property it comes out twice as big the wider way round.
	for _, size := range []int{40, 80} {
		for _, tc := range []struct {
			name string
			ve   string
			want int
		}{
			{"the scale the drawing gets by default", "", size / 2},
			{"the size that does not scale", ` vector-effect="non-scaling-size"`, 10},
		} {
			_, cv := painted(t, `<svg viewBox="0 0 40 40">
				<g transform="scale(2)">
					<rect width="10" height="10" fill="red"`+tc.ve+`/>
				</g>
			</svg>`, size, size)
			w, h := boxSpan(writtenBox(cv))
			if w < tc.want-1 || w > tc.want+1 || h < tc.want-1 || h > tc.want+1 {
				t.Errorf("%s at %d: the rectangle is %d by %d pixels, want %d", tc.name, size, w, h, tc.want)
			}
		}
	}
}

func TestRenderDropsTheTurnForANonRotation(t *testing.T) {
	// A rectangle ten by four written in a group that turned it a quarter turn
	// and scaled it twice comes out eight pixels wide and twenty tall; the turn
	// that does not happen drops the rotation, keeps the move and keeps the
	// scale, so it is twenty wide and eight tall. The size that does not scale
	// drops the scale instead and keeps the turn: four wide and ten tall.
	for _, tc := range []struct {
		name         string
		ve           string
		wantW, wantH int
	}{
		{"the turn and the scale the drawing gets by default", "", 8, 20},
		{"the turn that does not happen", ` vector-effect="non-rotation"`, 20, 8},
		{"the size that does not scale under the turn", ` vector-effect="non-scaling-size"`, 4, 10},
	} {
		_, cv := painted(t, `<svg viewBox="0 0 40 40">
			<g transform="translate(10,10) rotate(90) scale(2)">
				<rect width="10" height="4" fill="red"`+tc.ve+`/>
			</g>
		</svg>`, 40, 40)
		w, h := boxSpan(writtenBox(cv))
		if w < tc.wantW-1 || w > tc.wantW+1 || h < tc.wantH-1 || h > tc.wantH+1 {
			t.Errorf("%s: the rectangle is %d by %d pixels, want %d by %d", tc.name, w, h, tc.wantW, tc.wantH)
		}
	}
}

func TestRenderKeepsAFixedPositionWhereItsOwnTransformWritesIt(t *testing.T) {
	// fixed-position spends the shift the element's own transform writes as the
	// place it sits, whatever the elements above it did: the rectangle written
	// at four and eight inside two groups that carried the drawing twenty and
	// three to the right sits at four, and without the property at twenty-seven.
	for _, tc := range []struct {
		name string
		ve   string
		want int
	}{
		{"the place the drawing gets by default", "", 27},
		{"the place that does not move", ` vector-effect="fixed-position"`, 4},
	} {
		_, cv := painted(t, `<svg viewBox="0 0 40 40">
			<g transform="translate(20,0)">
				<g transform="translate(3,0)">
					<rect width="6" height="6" fill="red" transform="translate(4,8)"`+tc.ve+`/>
				</g>
			</g>
		</svg>`, 40, 40)
		if got := int(writtenBox(cv).MinX); got < tc.want-1 || got > tc.want+1 {
			t.Errorf("%s: the rectangle starts at x %d, want %d", tc.name, got, tc.want)
		}
	}
}

func TestRenderHoldsTheShapesInsideAnElementThatAsked(t *testing.T) {
	// A vector effect is not inherited, but the coordinate system it makes of
	// an element is the one the shapes inside it are drawn in: the outer group
	// asking for a size that does not scale holds the outer scale of two off
	// the rectangle inside, while the scale of two written on the inner group,
	// which asked for nothing, still happens to it.
	for _, tc := range []struct {
		name string
		ve   string
		want int
	}{
		{"the scale the drawing gets by default", "", 20},
		{"the outer scale held while the inner one still scales", ` vector-effect="non-scaling-size"`, 10},
	} {
		_, cv := painted(t, `<svg viewBox="0 0 40 40">
			<g transform="scale(2)"`+tc.ve+`>
				<g transform="scale(2)">
					<rect width="5" height="5" fill="red"/>
				</g>
			</g>
		</svg>`, 40, 40)
		w, h := boxSpan(writtenBox(cv))
		if w < tc.want-1 || w > tc.want+1 || h < tc.want-1 || h > tc.want+1 {
			t.Errorf("%s: the rectangle is %d by %d pixels, want %d", tc.name, w, h, tc.want)
		}
	}
}

func TestRenderHoldsTheWritingInsideAnElementThatAsked(t *testing.T) {
	// Writing is put down piece by piece from where the pen lands, and the pen
	// goes through the same transform a shape's geometry does, so the
	// coordinate system a non-scaling-size makes of a group holds the writing
	// inside it as well as the shapes.
	plain := boxOf(t, `<text x="2" y="16" font-size="8">Hi</text>`)
	held := boxOf(t, `<g transform="scale(2)" vector-effect="non-scaling-size"><text x="2" y="16" font-size="8">Hi</text></g>`)
	scaled := boxOf(t, `<g transform="scale(2)"><text x="2" y="16" font-size="8">Hi</text></g>`)
	if held != plain {
		t.Errorf("the writing inside the held group came out %s, want the same as the writing nothing scaled, %s", held, plain)
	}
	if scaled == plain {
		t.Error("the writing inside the scaled group did not grow, so the test says nothing")
	}
}

func TestRenderHoldsTheMarkersInsideAnElementThatAsked(t *testing.T) {
	// A marker is a picture put down at the vertex the shape built, through the
	// same transform the shape's geometry went through, so a coordinate system
	// held still holds the markers drawn on the shape too.
	draw := func(group, line string) box {
		_, cv := painted(t, `<svg viewBox="0 0 20 20">
			<defs>
				<marker id="m"><rect width="2" height="2" fill="red"/></marker>
			</defs>
			<g `+group+`>
				<line x1="5" y1="5" x2="8" y2="5" stroke="none" marker-start="url(#m)" `+line+`/>
			</g>
		</svg>`, 20, 20)
		return writtenBox(cv)
	}
	plain := draw("", "")
	held := draw(`transform="scale(2)"`, `vector-effect="non-scaling-size"`)
	scaled := draw(`transform="scale(2)"`, "")
	if held != plain {
		t.Errorf("the marker inside the held group came out %s, want the same as the one nothing scaled, %s", held, plain)
	}
	if scaled.Empty || scaled == plain {
		t.Errorf("the marker inside the scaled group came out %s, want it moved out of the one nothing scaled", scaled)
	}
}

func TestVectorEffectMatrixHoldsEveryCombination(t *testing.T) {
	// Every way the spec's three host-relative values combine, each as the
	// device point the element's own (1,0) lands on. The host turned the
	// drawing a quarter turn and scaled it twice — {0,2,-2,0} — and set it down
	// at (106,-4), the canvas is supersampled twice on top of that, and the
	// element's own transform writes the shift (6,-4), which is not where the
	// host put it. The shift is what separates fixed-position: without it the
	// element lands where the host carried it, with it where the file wrote it.
	//
	//	                        the point (1,0) of the element lands on
	//	none                    212,-4   the host's reading, the effect at rest
	//	non-scaling-size        212,-6   the turn and skew kept, the scale one
	//	non-rotation            216,-8   the scale kept, the turn and skew gone
	//	both                    214,-8   the host's linear part and nothing of it
	//	fixed-position           12,-4   the file's shift, the host's linear part
	//	size and fixed           12,-6   the file's shift, held at scale one
	//	rotation and fixed       16,-8   the file's shift, turned upright
	//	all three                14,-8   the file's shift alone
	//
	// The first row is the effect at rest, which [paintNodes] never asks for
	// because a value that holds nothing still does not constrain; the rest are
	// the values the spec names.
	turned := canvas.Matrix{B: 2, C: -2, E: 106, F: -4}
	for _, tc := range []struct {
		name string
		ve   VectorEffect
		want [2]float64
	}{
		{"none", VectorEffect{}, [2]float64{212, -4}},
		{"non-scaling-size", VectorEffect{NonScalingSize: true}, [2]float64{212, -6}},
		{"non-rotation", VectorEffect{NonRotation: true}, [2]float64{216, -8}},
		{"size and rotation", VectorEffect{NonScalingSize: true, NonRotation: true}, [2]float64{214, -8}},
		{"fixed-position", VectorEffect{FixedPosition: true}, [2]float64{12, -4}},
		{"size and fixed-position", VectorEffect{NonScalingSize: true, FixedPosition: true}, [2]float64{12, -6}},
		{"rotation and fixed-position", VectorEffect{NonRotation: true, FixedPosition: true}, [2]float64{16, -8}},
		{"all three", VectorEffect{NonScalingSize: true, NonRotation: true, FixedPosition: true}, [2]float64{14, -8}},
	} {
		st := Style{Transform: turned}
		st.VectorEffect = tc.ve
		st.VectorEffect.OriginX, st.VectorEffect.OriginY = 6, -4
		at := vectorEffectMatrix(canvas.Scale(2, 2), st).Mul(st.Transform)
		x, y := at.Map(1, 0)
		if x != tc.want[0] || y != tc.want[1] {
			t.Errorf("%s: the element's (1,0) landed on %g,%g, want %g,%g", tc.name, x, y, tc.want[0], tc.want[1])
		}
	}
}

func TestVectorEffectMatrixLeavesADegenerateHostAlone(t *testing.T) {
	// A host that flattened one axis has no determinant to divide by, and a
	// transform with no way back has no space to hold the shape in: both leave
	// the element as SVG would have drawn it, rather than making an infinity of
	// it. The second is a matrix the determinant does not quite catch, which is
	// why the way back is asked for as well.
	flat := Style{Transform: canvas.Scale(1, 0)}
	flat.VectorEffect.NonScalingSize = true
	if got := vectorEffectMatrix(canvas.Scale(2, 2), flat); got != canvas.Scale(2, 2) {
		t.Errorf("a host with no determinant came out %+v, want the drawing's own matrix", got)
	}
	tiny := Style{Transform: canvas.Scale(1e-13, 1e-13)}
	tiny.VectorEffect.NonRotation = true
	if got := vectorEffectMatrix(canvas.Scale(2, 2), tiny); got != canvas.Scale(2, 2) {
		t.Errorf("a transform with no way back came out %+v, want the drawing's own matrix", got)
	}
}

func TestRenderFixesThePositionWithoutDroppingTheTurn(t *testing.T) {
	// fixed-position spends the shift the element's own transform writes as the
	// place it sits and keeps the host's turn and scale: the rectangle written
	// at (4,8) inside a group that turned the drawing a quarter turn comes out
	// turned, four wide and ten tall, hanging from its own origin at four and
	// eight. The same rectangle without the property comes out where the turn
	// carried it, and asking for the turn that does not happen as well uprights
	// it there.
	for _, tc := range []struct {
		name         string
		ve           string
		wantX, wantY int
		wantW, wantH int
	}{
		{"the turn and the place the drawing gets by default", "", 8, 24, 4, 10},
		{"the place that does not move", ` vector-effect="fixed-position"`, 0, 8, 4, 10},
		{"the place that does not move, turned upright", ` vector-effect="fixed-position non-rotation"`, 4, 8, 10, 4},
	} {
		_, cv := painted(t, `<svg viewBox="0 0 40 40">
			<g transform="translate(20,20) rotate(90)">
				<rect width="10" height="4" fill="red" transform="translate(4,8)"`+tc.ve+`/>
			</g>
		</svg>`, 40, 40)
		b := writtenBox(cv)
		w, h := boxSpan(b)
		if int(b.MinX) < tc.wantX-1 || int(b.MinX) > tc.wantX+1 ||
			int(b.MinY) < tc.wantY-1 || int(b.MinY) > tc.wantY+1 ||
			w < tc.wantW-1 || w > tc.wantW+1 || h < tc.wantH-1 || h > tc.wantH+1 {
			t.Errorf("%s: the rectangle is %d by %d pixels at (%g,%g), want %d by %d at (%d,%d)",
				tc.name, w, h, b.MinX, b.MinY, tc.wantW, tc.wantH, tc.wantX, tc.wantY)
		}
	}
}
