package template

import "github.com/gabrielluizsf/antui/template/css"

// entryStyle is the style a widget actually paints and lays out with: the
// cascade style animated, then transitioned. When the requested state moved
// off the one the widget last held, a transition starts toward the new
// target's style; while it runs, the entry draws the blended value.
func (cs *cssStyle) entryStyle(e *widgetEntry, role string, requested State) css.Style {
	now := cs.now()
	target := cs.animate(cs.baseStyle(role, requested))
	if e.fresh {
		// A widget that just appeared has no previous style to transition
		// from: it snaps to its resting target so the first frame of a
		// widget born mid-state (a hovered or checked element) does not
		// animate from an empty canvas.
		e.fresh = false
		e.state = requested
		e.displayed = target
		return e.displayed
	}
	if requested != e.state {
		e.state = requested
		e.from = e.displayed
		e.to = target
		e.started = now
		e.budget = css.TransitionSpan(target)
		e.active = len(target.Transitions) > 0 && e.budget > 0
	}
	if e.active {
		elapsed := (now - e.started) * 1000
		if elapsed < e.budget {
			e.displayed = css.BlendTransition(css.BlendTransitionParams{From: e.from, To: e.to, Elapsed: elapsed, Units: cs.unitsCtx(e.to)})
			return e.displayed
		}
		e.active = false
	}
	e.displayed = target
	return e.displayed
}
