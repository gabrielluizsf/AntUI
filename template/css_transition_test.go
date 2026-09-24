package template

import (
	"testing"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// newCSSTemplate wires a template over a stylesheet and a fake clock, so a
// transition can be read at exact moments instead of waiting real seconds.
// The returned now is the clock's cell: write the seconds into it between
// frames.
func newCSSTemplate(t *testing.T, sheet string) (*CSS, *float64) {
	t.Helper()
	win, _, err := antui.Offscreen(400, 200)
	if err != nil {
		t.Fatal(err)
	}
	tpl := TemplateWithCSS(win)
	classes, path := cssTable(t, sheet)
	if err := tpl.SetStyle(path, classes); err != nil {
		t.Fatal(err)
	}
	var clock float64
	tpl.style.now = func() float64 { return clock }
	return tpl, &clock
}

// TestCSSTransitionGlideOnHover a hover starts a transition from the resting
// colour to the hovered one: halfway through the box shows the colour in
// between, and at the end the hovered colour itself.
func TestCSSTransitionGlideOnHover(t *testing.T) {
	tpl, clock := newCSSTemplate(t, `
		button { background-color: #111111; transition: background-color 1s; }
		button:hover { background-color: #333333; }
	`)

	*clock = 0
	_, _, _, _, rest, ok := tpl.layout(css.RoleButton, "Go")
	if !ok {
		t.Fatal("button should be visible")
	}
	if rest.Background != 0xFF111111 {
		t.Errorf("resting background = %v, want #111111", rest.Background)
	}

	// The pointer arrives: the transition starts now, so at this instant
	// the box still shows the colour it settled on.
	*clock = 0
	early := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})
	if early.Background != 0xFF111111 {
		t.Errorf("background at transition start = %v, want #111111", early.Background)
	}

	*clock = 0.5
	mid := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})
	if want := canvas.Mix(0xFF111111, 0xFF333333, 0.5); mid.Background != want {
		t.Errorf("background halfway = %v, want %v", mid.Background, want)
	}

	*clock = 1
	done := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})
	if done.Background != 0xFF333333 {
		t.Errorf("background after 1s = %v, want #333333", done.Background)
	}
}

// TestCSSTransitionDelayWaits the delay runs before any blending: the box
// keeps the resting colour while it waits, then glides once the duration
// begins.
func TestCSSTransitionDelayWaits(t *testing.T) {
	tpl, clock := newCSSTemplate(t, `
		button { background-color: #111111; transition: background-color 1s 0.5s; }
		button:hover { background-color: #333333; }
	`)

	*clock = 0
	tpl.layout(css.RoleButton, "Go")
	tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})

	// A quarter second in: still inside the half-second delay.
	*clock = 0.25
	waiting := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})
	if waiting.Background != 0xFF111111 {
		t.Errorf("background during delay = %v, want #111111", waiting.Background)
	}

	// 0.75s in: a quarter of the duration has passed, so a quarter of the
	// way from resting to hovered.
	*clock = 0.75
	gone := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})
	if want := canvas.Mix(0xFF111111, 0xFF333333, 0.25); gone.Background != want {
		t.Errorf("background a quarter in = %v, want %v", gone.Background, want)
	}

	// The full budget (delay + duration) reaches the target.
	*clock = 1.5
	finished := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})
	if finished.Background != 0xFF333333 {
		t.Errorf("background after the delay and duration = %v, want #333333", finished.Background)
	}
}

// TestCSSTransitionReverse a state change mid-flight starts a new transition
// from wherever the box is now, heading the other way.
func TestCSSTransitionReverse(t *testing.T) {
	tpl, clock := newCSSTemplate(t, `
		button { background-color: #111111; transition: background-color 1s; }
		button:hover { background-color: #333333; }
	`)

	*clock = 0
	tpl.layout(css.RoleButton, "Go")
	tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})

	*clock = 0.5
	half := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})
	// The pointer leaves halfway through: the new transition starts from
	// this colour, not from either endpoint.
	*clock = 0.5
	back := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{})
	if back.Background != half.Background {
		t.Errorf("reverse begins where the box was: got %v, want %v", back.Background, half.Background)
	}

	*clock = 0.75
	toward := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{})
	if want := canvas.Mix(half.Background, 0xFF111111, 0.25); toward.Background != want {
		t.Errorf("background heading back = %v, want %v", toward.Background, want)
	}

	*clock = 1.5
	rest := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{})
	if rest.Background != 0xFF111111 {
		t.Errorf("background after the return = %v, want #111111", rest.Background)
	}
}

// TestCSSTransitionEndsAfterBudget once the budget runs out the entry stops
// blending: every later frame reads the target itself, whether the clock
// keeps running or not.
func TestCSSTransitionEndsAfterBudget(t *testing.T) {
	tpl, clock := newCSSTemplate(t, `
		button { background-color: #111111; transition: background-color 1s; }
		button:hover { background-color: #333333; }
	`)

	*clock = 0
	tpl.layout(css.RoleButton, "Go")
	tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})
	if !tpl.style.cur.active {
		t.Fatal("hover should be mid-transition")
	}

	*clock = 2
	after := tpl.style.entryStyle(tpl.style.cur, css.RoleButton, State{Hovered: true})
	if tpl.style.cur.active {
		t.Error("a finished transition must stop being active")
	}
	if after.Background != 0xFF333333 {
		t.Errorf("background past the budget = %v, want #333333", after.Background)
	}

	// The same answer holds when the widget is laid out again next frame.
	tpl.Reset()
	*clock = 2
	_, _, _, _, again, ok := tpl.layout(css.RoleButton, "Go")
	if !ok {
		t.Fatal("button should still be visible")
	}
	if again.Background != 0xFF333333 {
		t.Errorf("background after a reset = %v, want #333333", again.Background)
	}
}

// TestCSSTransitionFirstAppearanceSnaps a widget that is born already in a
// mid-state (hovered, checked, focused) has no resting frame to glide from:
// the timeline snaps to the target so a freshly drawn widget never washes in
// from an empty canvas.
func TestCSSTransitionFirstAppearanceSnaps(t *testing.T) {
	tpl, clock := newCSSTemplate(t, `
		button { background-color: #111111; transition: background-color 1s; }
		button:hover { background-color: #333333; }
	`)
	e := tpl.style.beginWidget(css.RoleButton, "Go")
	*clock = 0.5
	st := tpl.style.entryStyle(e, css.RoleButton, State{Hovered: true})
	if st.Background != 0xFF333333 {
		t.Errorf("first-appear hovered background = %v, want #333333", st.Background)
	}
}

// TestCSSTransitionReappearSnaps a widget that left the page and comes back
// is a new widget as far as the timeline goes: whatever it was before the
// disappearance, its reappearance in a state is a jump, not a glide.
func TestCSSTransitionReappearSnaps(t *testing.T) {
	tpl, clock := newCSSTemplate(t, `
		button { background-color: #111111; transition: background-color 1s; }
		button:hover { background-color: #333333; }
	`)
	*clock = 0
	if _, _, _, _, _, ok := tpl.layout(css.RoleButton, "Go"); !ok {
		t.Fatal("button should be visible")
	}
	tpl.Reset()
	if _, _, _, _, _, ok := tpl.layout(css.RoleButton, "Go"); !ok {
		t.Fatal("button should be visible again")
	}
	// Three empty frames prune the entry, so the next appearance is born.
	tpl.style.resetTimeline()
	tpl.style.resetTimeline()
	tpl.style.resetTimeline()
	e := tpl.style.beginWidget(css.RoleButton, "Go")
	*clock = 0.5
	st := tpl.style.entryStyle(e, css.RoleButton, State{Hovered: true})
	if st.Background != 0xFF333333 {
		t.Errorf("reappeared hovered background = %v, want #333333", st.Background)
	}
}
