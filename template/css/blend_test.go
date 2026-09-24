package css

import (
	"math"
	"testing"
)

// TestBlendColors colours interpolate channel by channel across the blend,
// red into cyan meeting in the middle as the colour between them.
func TestBlendColors(t *testing.T) {
	a := mustParse(t, `button { background-color: #111111; }`).Style("button", nil, StateNone, 800)
	b := mustParse(t, `button { background-color: #333333; }`).Style("button", nil, StateNone, 800)
	out := BlendStyles(a, b, 0.5, map[string]bool{"background-color": true}, Units{})
	if out.Background != 0xFF222222 {
		t.Errorf("blend = %v, want 0xFF222222", out.Background)
	}
}

// TestBlendStylesOnlyListed a blend only touches the properties it is asked
// for; everything else stays the target style's values.
func TestBlendStylesOnlyListed(t *testing.T) {
	a := mustParse(t, `button { opacity: 0; }`).Style("button", nil, StateNone, 800)
	b := mustParse(t, `button { opacity: 1; color: lime; }`).Style("button", nil, StateNone, 800)
	out := BlendStyles(a, b, 0.5, map[string]bool{"opacity": true}, Units{})
	if out.Opacity != 0.5 {
		t.Errorf("opacity = %v, want 0.5", out.Opacity)
	}
	if out.Color != 0xFF00FF00 {
		t.Errorf("unlisted color must stay b's, got %v", out.Color)
	}
}

// TestBlendOpacity a plain float property moves linearly with the progress.
func TestBlendOpacity(t *testing.T) {
	a, b := Style{Opacity: 0.2}, Style{Opacity: 0.8}
	out := BlendStyles(a, b, 0.25, map[string]bool{"opacity": true}, Units{})
	if math.Abs(out.Opacity-0.35) > 1e-9 {
		t.Errorf("opacity = %v, want 0.35", out.Opacity)
	}
}

// TestBlendLengthSameUnit two widths written in the same unit stay in that
// unit all the way, so a 50%-to-75% width reads 62.5% at the middle.
func TestBlendLengthSameUnit(t *testing.T) {
	a := mustParse(t, `button { width: 50%; }`).Style("button", nil, StateNone, 800)
	b := mustParse(t, `button { width: 75%; }`).Style("button", nil, StateNone, 800)
	out := BlendStyles(a, b, 0.5, map[string]bool{"width": true}, Units{Width: 800})
	if !out.Width.IsPct() {
		t.Fatalf("width unit = %v, want percent", out.Width.Unit())
	}
	if got := out.Width.Resolve(Units{Width: 800}); got != 500 {
		t.Errorf("middle width = %vpx, want 500", got)
	}
}

// TestBlendLengthCrossUnit a fixed value against a percentage resolves both
// to pixels and lerps the measurement, coming out a fixed length.
func TestBlendLengthCrossUnit(t *testing.T) {
	a := Style{Width: Fixed(100)}
	b := Style{Width: Pct(50)}
	out := BlendStyles(a, b, 0.5, map[string]bool{"width": true}, Units{Width: 800})
	if got := out.Width.Resolve(Units{Width: 800}); got != 250 {
		t.Errorf("cross-unit width = %vpx, want 250", got)
	}
}

// TestBlendLengthKeywordsDiscrete lengths that cannot mix — auto, none — swap
// over at the half way point instead of pretending to be numbers.
func TestBlendLengthKeywordsDiscrete(t *testing.T) {
	a, b := Style{Width: Auto()}, Style{Width: Fixed(200)}
	early := BlendStyles(a, b, 0.4, map[string]bool{"width": true}, Units{})
	late := BlendStyles(a, b, 0.6, map[string]bool{"width": true}, Units{})
	if !early.Width.Auto() || late.Width.Auto() {
		t.Errorf("auto-to-px blend = %v (early), %v (late)", early.Width, late.Width)
	}
}

// TestBlendTransitionDelayedDelay a transition holds the entry value until
// its delay has passed, walks the easing between, and lands on the target.
func TestBlendTransitionDelayedDelay(t *testing.T) {
	from := mustParse(t, `button { background-color: #111111; }`).Style("button", nil, StateNone, 800)
	to := mustParse(t, `
		button { background-color: #333333; transition: background-color 1s 0.5s; }
	`).Style("button", nil, StateNone, 800)
	if hold := BlendTransition(from, to, 300, Units{}); hold.Background != 0xFF111111 {
		t.Errorf("during the delay = %v, want the entry colour", hold.Background)
	}
	if mid := BlendTransition(from, to, 800, Units{}); mid.Background != 0xFF1B1B1B {
		t.Errorf("mid transition = %v, want 0xFF1B1B1B", mid.Background)
	}
	if done := BlendTransition(from, to, 1500, Units{}); done.Background != 0xFF333333 {
		t.Errorf("after the span = %v, want the target", done.Background)
	}
	if settled := BlendTransition(from, to, 20000, Units{}); settled.Background != 0xFF333333 {
		t.Errorf("long after = %v, want the target", settled.Background)
	}
}

// TestBlendTransitionAll property "all" transitions everything the engine can
// interpolate, holding even unlisted properties of the target.
func TestBlendTransitionAll(t *testing.T) {
	from := mustParse(t, `button { background-color: #111111; }`).Style("button", nil, StateNone, 800)
	to := mustParse(t, `
		button { background-color: #333333; transition: 0.5s; }
	`).Style("button", nil, StateNone, 800)
	if len(to.Transitions) != 1 || to.Transitions[0].Prop != "all" {
		t.Fatalf("transitions = %+v", to.Transitions)
	}
	mid := BlendTransition(from, to, 250, Units{})
	if mid.Background != 0xFF222222 {
		t.Errorf("all-transition at half = %v, want 0xFF222222", mid.Background)
	}
}

// TestBlendTransitionNoTransition a style with no transitions hands the
// target back untouched, whatever the clock says.
func TestBlendTransitionNoTransition(t *testing.T) {
	from := mustParse(t, `button { background-color: #111111; }`).Style("button", nil, StateNone, 800)
	to := mustParse(t, `button { background-color: #333333; }`).Style("button", nil, StateNone, 800)
	if got := BlendTransition(from, to, 1000, Units{}); got.Background != 0xFF333333 {
		t.Errorf("no transition = %v, want the target", got.Background)
	}
}

// TestBlendTransformSameKind matching transform functions interpolate their
// arguments; a listen that cannot align swaps at the half way point.
func TestBlendTransformSameKind(t *testing.T) {
	a := Style{Transform: []TransformFunc{{Kind: TransformTranslate, Dx: Fixed(0)}}}
	b := Style{Transform: []TransformFunc{{Kind: TransformTranslate, Dx: Fixed(100)}}}
	out := BlendStyles(a, b, 0.5, map[string]bool{"transform": true}, Units{})
	if len(out.Transform) != 1 || out.Transform[0].Dx.Px(1) != 50 {
		t.Errorf("translate blend = %+v", out.Transform)
	}

	mixed := BlendStyles(
		Style{Transform: []TransformFunc{{Kind: TransformScale, Sx: 1}}},
		Style{Transform: []TransformFunc{{Kind: TransformRotate, Ax: Deg(45)}}},
		0.6, map[string]bool{"transform": true}, Units{},
	)
	if len(mixed.Transform) != 1 || mixed.Transform[0].Kind != TransformRotate {
		t.Errorf("mismatched transform = %+v, want the target past the swap", mixed.Transform)
	}
}

// TestBlendShadowDiscrete a shadow list that grows or shrinks cannot align,
// so it swaps wholesale at the boundary instead of inventing values.
func TestBlendShadowDiscrete(t *testing.T) {
	a := Style{}
	b := Style{BoxShadow: []Shadow{{X: 0, Y: 2, Blur: 4, Color: 0x80000000}}}
	early := BlendStyles(a, b, 0.4, map[string]bool{"box-shadow": true}, Units{})
	late := BlendStyles(a, b, 0.6, map[string]bool{"box-shadow": true}, Units{})
	if len(early.BoxShadow) != 0 || len(late.BoxShadow) != 1 {
		t.Errorf("shadow swap = %v (early), %v (late)", early.BoxShadow, late.BoxShadow)
	}
}

// TestCopyStyleOwnsSet a copied style has its own Set map, so marking a value
// in the copy can never reach back into the cascade's cached style.
func TestCopyStyleOwnsSet(t *testing.T) {
	a := mustParse(t, `button { opacity: 0; }`).Style("button", nil, StateNone, 800)
	b := CopyStyle(a)
	b.Set["opacity"] = true
	b.Opacity = 1
	if a.Opacity == b.Opacity {
		t.Error("the copy must hold its own value")
	}
	if a.Set["opacity"] != b.Set["opacity"] {
		t.Error("shared maps must not lie")
	}
	d := CopyStyle(a)
	d.Set["flavor"] = true
	if a.Set["flavor"] {
		t.Error("copy's Set map aliases the original")
	}
}

// TestBlendKeyframesBaseFallback a frame that omitted a property the other
// frame set travels from the base style's value, not from emptiness.
func TestBlendKeyframesBaseFallback(t *testing.T) {
	base := mustParse(t, `button { background-color: #222222; }`).Style("button", nil, StateNone, 800)
	kf := mustParse(t, `
		@keyframes fade { from { background-color: #111111; } to { opacity: 1; } }
	`).Keyframes("fade")
	out := Animate(base, kf, 0.25, Units{})
	if out.Background != 0xFF151515 {
		t.Errorf("background = %v, want 0xFF151515 (quarter way from #111 to base #222)", out.Background)
	}
}
