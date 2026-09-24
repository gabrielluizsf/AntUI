package template

import (
	"math"
	"strconv"

	"github.com/gabrielluizsf/antui/template/css"
)

// widgetEntry is the timeline of one widget on the page: the interaction
// state it painted last, the style it showed when the state changed, the
// style it is heading toward, and when in the clock the change happened. It
// is what turns a hover from a jump into a glide, because the entry keeps
// hold of the whole motion between frames.
type widgetEntry struct {
	role      string
	state     State
	from      css.Style
	to        css.Style
	started   float64 // timeline seconds when the transition began
	budget    float64 // how long the longest transition runs, in milliseconds
	active    bool
	displayed css.Style
	seen      int  // the frame the widget was last drawn in
	fresh     bool // the entry has never been painted before
}

// beginWidget opens the timeline of the widget being laid out. The key is
// the role plus its ordinal in this frame plus its label, which pins a
// widget through its frame and back to the same entry next frame. The entry
// becomes cs.cur, so the painter's styleFor calls land on it.
func (cs *cssStyle) beginWidget(role, label string) *widgetEntry {
	cs.seq++
	key := role + "\x00" + strconv.Itoa(cs.seq) + "\x00" + label
	e := cs.timeline[key]
	if e == nil {
		e = &widgetEntry{role: role, fresh: true}
		cs.timeline[key] = e
	}
	e.seen = cs.frame
	cs.cur = e
	return e
}

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
			e.displayed = css.BlendTransition(e.from, e.to, elapsed, cs.unitsCtx(e.to))
			return e.displayed
		}
		e.active = false
	}
	e.displayed = target
	return e.displayed
}

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

// animationFrame turns the clock into the 0..1 progress read of one
// animation: the delay is skipped first, then the iteration and its
// direction pick the frame. on reports whether the animation is playing or
// pinned there by a fill-mode; pinned animations still return frames.
func animationFrame(a css.Animation, now float64) (float64, bool) {
	dur := a.Duration.Sec()
	if dur <= 0 {
		return 0, false
	}
	el := now - a.Delay.Sec()
	if el < 0 {
		if a.Fill == css.FillBackwards || a.Fill == css.FillBoth {
			return 0, true
		}
		return 0, false
	}
	if !a.Infinite && a.Iterations >= 0 && el >= dur*a.Iterations {
		if a.Fill == css.FillForwards || a.Fill == css.FillBoth {
			return 1, true
		}
		return 0, false
	}
	cycle := int(el / dur)
	local := math.Mod(el, dur) / dur
	switch a.Direction {
	case css.AnimReverse:
		return 1 - local, true
	case css.AnimAlternate:
		if cycle%2 == 1 {
			return 1 - local, true
		}
		return local, true
	case css.AnimAlternateReverse:
		if cycle%2 == 1 {
			return local, true
		}
		return 1 - local, true
	}
	return local, true
}

// resetTimeline opens a new frame: the widget ordinals restart at zero and
// the current entry reference clears, so the next layout begins a fresh
// run. An entry not drawn for two frames is dropped, so a widget that
// stopped appearing (or changed role and header) does not linger in the map.
func (cs *cssStyle) resetTimeline() {
	cs.frame++
	cs.seq = 0
	cs.cur = nil
	for k, e := range cs.timeline {
		if e.seen+2 <= cs.frame {
			delete(cs.timeline, k)
		}
	}
}
