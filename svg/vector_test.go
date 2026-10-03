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
	// The two values that mean something are read from the attribute and from a
	// `style` declaration alike, and any other vector effect is said rather than
	// quietly painted as though it had not been asked for. An attribute written
	// with nothing in it means nothing at all, as everywhere else.
	for _, tc := range []struct {
		name  string
		value string
		warns int
		want  bool
	}{
		{"nothing written", "", 0, false},
		{"the stroke a drawing gets by default", "none", 0, false},
		{"the stroke that does not scale", "non-scaling-stroke", 0, true},
		{"the keyword in another case", "NON-SCALING-STROKE", 0, true},
		{"a vector effect this package does not do", "non-scaling-size", 1, false},
		{"another one of them", "non-rotation", 1, false},
		{"and a third", "fixed-position", 1, false},
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
		if line.Style.NonScalingStroke != tc.want {
			t.Errorf("%s: non-scaling stroke = %v, want %v", tc.name, line.Style.NonScalingStroke, tc.want)
		}
	}
}

func TestVectorEffectIsReadFromTheStyleAttribute(t *testing.T) {
	// A declaration says the same thing the attribute does, and the one written
	// in `style` is the one that counts where both are written.
	st := Style{}.with(mustElement(t,
		`<line vector-effect="none" style="vector-effect: non-scaling-stroke"/>`), quiet)
	if !st.NonScalingStroke {
		t.Error("the declaration did not beat the attribute of the same name")
	}
	st = Style{}.with(mustElement(t,
		`<line vector-effect="non-scaling-stroke" style="vector-effect: none"/>`), quiet)
	if st.NonScalingStroke {
		t.Error("the declaration did not beat the attribute of the same name")
	}
}

func TestVectorEffectIsNotInherited(t *testing.T) {
	// The spec has vector-effect as not inherited, so what a group writes is
	// about the group — which has no outline of its own to compute — and the
	// shapes inside take the scaled stroke every drawing gives them.
	group := Style{}.with(mustElement(t, `<g vector-effect="non-scaling-stroke"/>`), quiet)
	if !group.NonScalingStroke {
		t.Fatal("the group did not take the vector-effect it was given")
	}
	inner := group.with(mustElement(t, `<line/>`), quiet)
	if inner.NonScalingStroke {
		t.Error("the vector-effect travelled down to the shape, want it to stay on the group")
	}
	own := inner.with(mustElement(t, `<line vector-effect="non-scaling-stroke"/>`), quiet)
	if !own.NonScalingStroke {
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
