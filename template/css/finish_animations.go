package css

import (
	"math"
)

// finishAnimations merges the animation longhands (and the shorthand's own
// parallel lists) into the style's final Animations slice. The name list
// sets the length, and — like CSS — every shorter list cycles to cover the
// whole row, or leaves its CSS initial where the longhand never appeared.
func finishAnimations(st *Style) {
	if len(st.AnimationNames) == 0 {
		st.Animations = nil
		return
	}
	for i, name := range st.AnimationNames {
		a := Animation{Name: name, Iterations: 1, Direction: AnimNormal, Fill: FillNone}
		if len(st.AnimationDurs) > 0 {
			a.Duration = st.AnimationDurs[i%len(st.AnimationDurs)]
		}
		if len(st.AnimationTims) > 0 {
			a.Timing = st.AnimationTims[i%len(st.AnimationTims)]
		}
		if len(st.AnimationDels) > 0 {
			a.Delay = st.AnimationDels[i%len(st.AnimationDels)]
		}
		if len(st.AnimationIters) > 0 {
			if v := st.AnimationIters[i%len(st.AnimationIters)]; math.IsInf(v, 1) {
				a.Infinite = true
			} else {
				a.Iterations = v
			}
		}
		if len(st.AnimationDirs) > 0 {
			a.Direction = st.AnimationDirs[i%len(st.AnimationDirs)]
		}
		if len(st.AnimationFills) > 0 {
			a.Fill = st.AnimationFills[i%len(st.AnimationFills)]
		}
		st.Animations = append(st.Animations, a)
	}
}
