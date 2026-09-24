package css

import (
	"testing"
)

// TestKeyframesParse reads a @keyframes block: its name, its frames in order,
// and each frame's offset and declarations.
func TestKeyframesParse(t *testing.T) {
	sh := mustParse(t, `
		@keyframes fade {
			from { opacity: 0; }
			50%  { opacity: 0.5; }
			to   { opacity: 1; }
		}
	`)
	kf := sh.Keyframes("fade")
	if kf == nil {
		t.Fatal("Keyframes(fade) is nil")
	}
	if kf.Name != "fade" {
		t.Errorf("name = %q", kf.Name)
	}
	if len(kf.Frames) != 3 {
		t.Fatalf("frames = %d, want 3", len(kf.Frames))
	}
	if kf.Frames[0].Offset != 0 || kf.Frames[1].Offset != 0.5 || kf.Frames[2].Offset != 1 {
		t.Errorf("offsets = %v, %v, %v", kf.Frames[0].Offset, kf.Frames[1].Offset, kf.Frames[2].Offset)
	}
	if got := sh.Style("button", nil, StateNone, 800).Animations; got != nil {
		t.Errorf("a sheet with only keyframes must style no animations, got %v", got)
	}
}

// TestKeyframesVendorPrefix vendor-prefixed @keyframes are read too, so one
// stylesheet can carry them without tripping the parser.
func TestKeyframesVendorPrefix(t *testing.T) {
	sh := mustParse(t, `
		@-webkit-keyframes spin { to { transform: rotate(360deg); } }
	`)
	if sh.Keyframes("spin") == nil {
		t.Error("-webkit-keyframes name should be readable as spin")
	}
}

// TestKeyframesClampOffsets frames beyond the range clamp into [0,1], the
// way a percentage over 100% pins to the end.
func TestKeyframesClampOffsets(t *testing.T) {
	sh := mustParse(t, `
		@keyframes wide {
			150% { opacity: 0.5; }
			-10% { opacity: 0.2; }
		}
	`)
	kf := sh.Keyframes("wide")
	if kf.Frames[0].Offset != 1 || kf.Frames[1].Offset != 0 {
		t.Errorf("clamped offsets = %v, %v", kf.Frames[0].Offset, kf.Frames[1].Offset)
	}
}

// TestKeyframesJunkSkipped a frame the parser does not understand is consumed
// and dropped; the block and the sheet around it survive.
func TestKeyframesJunkSkipped(t *testing.T) {
	sh := mustParse(t, `
		@keyframes mixed {
			from { opacity: 0; }
			midway { opacity: 0.5; }
			to { opacity: 1; }
		}
		button { color: red; }
	`)
	kf := sh.Keyframes("mixed")
	if len(kf.Frames) != 2 {
		t.Fatalf("frames = %d, want the two valid ones", len(kf.Frames))
	}
	if color := sh.Style("button", nil, StateNone, 800).Color; color != 0xFFFF0000 {
		t.Errorf("button after keyframes = %v, want red", color)
	}
}

// TestKeyframesFrameStyle a frame's declarations fold into the base the way a
// rule does: colors resolve to their current rgb, so frameStyle serves the
// blender ready values.
func TestKeyframesFrameStyle(t *testing.T) {
	sh := mustParse(t, `@keyframes fade { from { color: rgb(255, 0, 0); } }`)
	base := Style{Custom: map[string]string{}}
	st := frameStyle(sh.Keyframes("fade").Frames[0], base, Units{})
	if st.Color != 0xFFFF0000 || !st.Set["color"] {
		t.Errorf("frame color = %v, set=%v", st.Color, st.Set)
	}
}

// TestKeyframesFrameSpan the span between two frames finds its local
// progress, pins before the first frame and after the last.
func TestKeyframesFrameSpan(t *testing.T) {
	frames := []Keyframe{{Offset: 0}, {Offset: 0.5}, {Offset: 1}}
	lo, hi, local, ok := frameSpan(frames, 0.25)
	if !ok || lo.Offset != 0 || hi.Offset != 0.5 || local != 0.5 {
		t.Errorf("mid span = %+v %+v %v %v", lo, hi, local, ok)
	}
	if lo, _, _, _ = frameSpan(frames, 0); lo.Offset != 0 {
		t.Errorf("start should pin to the first frame, got %v", lo.Offset)
	}
	if lo, _, _, _ = frameSpan(frames, 1); lo.Offset != 1 {
		t.Errorf("end should pin to the last frame, got %v", lo.Offset)
	}
}

// TestAnimateMiddle a frame at each end interpolates the animated property
// across the space between them.
func TestAnimateMiddle(t *testing.T) {
	sh := mustParse(t, `
		@keyframes fade { from { opacity: 0; } to { opacity: 1; } }
	`)
	base := Style{}
	kf := sh.Keyframes("fade")
	ctx := Units{Width: 800}
	mid := Animate(base, kf, 0.5, ctx)
	if mid.Opacity != 0.5 {
		t.Errorf("Animate(0.5).Opacity = %v, want 0.5", mid.Opacity)
	}
	start := Animate(base, kf, 0, ctx)
	if start.Opacity != 0 {
		t.Errorf("Animate(0).Opacity = %v, want 0", start.Opacity)
	}
	end := Animate(base, kf, 1, ctx)
	if end.Opacity != 1 {
		t.Errorf("Animate(1).Opacity = %v, want 1", end.Opacity)
	}
}

// TestAnimateBaseUntouched an animated frame a keyframe omitted leaves the
// base style's value standing, and the base style itself is never modified
// by the animation.
func TestAnimateBaseUntouched(t *testing.T) {
	sh := mustParse(t, `
		@keyframes pop { to { opacity: 0.5; } }
	`)
	base := Style{}
	kf := sh.Keyframes("pop")
	out := Animate(base, kf, 1, ctx(800))
	if out.Opacity != 0.5 {
		t.Errorf("Animate(1).Opacity = %v, want 0.5", out.Opacity)
	}
	if base.Opacity != 0 {
		t.Errorf("the base style was mutated: %v", base.Opacity)
	}
}

// TestAnimateMissingBlock a named animation with no @keyframes block plays
// nothing: the base style is returned untouched.
func TestAnimateMissingBlock(t *testing.T) {
	kf := mustParse(t, `button { color: red; }`).Keyframes("ghost")
	if kf != nil {
		t.Fatal("no block should answer nil")
	}
	base := mustParse(t, `button { opacity: 0.75; }`).Style("button", nil, StateNone, 800)
	out := Animate(base, kf, 1, ctx(800))
	if out.Opacity != 0.75 {
		t.Errorf("Animate(nil block).Opacity = %v, want the base 0.75", out.Opacity)
	}
}

func ctx(w int) Units { return Units{Width: w} }
