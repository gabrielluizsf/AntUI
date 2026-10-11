package template

import "github.com/gabrielluizsf/antui/template/css"

// animate folds every running animation of the base style over it, in the
// order the animations were declared. An animation whose @keyframes block is
// missing is skipped; so is one waiting out its delay or already finished,
// unless a fill-mode pins its frames. The animation's own timing function
// reshapes the progress the frames read.
func (cs *cssStyle) animate(base css.Style) css.Style {
	if len(base.Animations) == 0 {
		return base
	}
	out := base
	now := cs.now()
	for _, a := range base.Animations {
		kf := cs.classes.Keyframes(a.Name)
		if kf == nil {
			continue
		}
		t, on := animationFrame(a, now)
		if !on {
			continue
		}
		out = css.Animate(out, kf, a.Timing.Ease(t), cs.unitsCtx(out))
	}
	return out
}
