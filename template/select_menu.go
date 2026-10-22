package template

import "github.com/gabrielluizsf/antui/canvas"

import "github.com/gabrielluizsf/antui"

// selectMenu draws the open menu, one clickable row per option, and reports
// whether one was picked.
func (t *uiTemplate) selectMenu(win *antui.Window, x, y, w, h int, index *int, options []string, focused bool) bool {
	type optionDraw struct {
		x, y  int
		label string
		sel   bool
		state State
	}
	draws := make([]optionDraw, 0, len(options))
	picked := false
	for i, opt := range options {
		oy := y + h*(i+1)
		oid := win.WidgetID("template:option", canvas.Area{X: x, Y: oy, Width: w, Height: h}, opt)
		oh := win.Hovered(canvas.Area{X: x, Y: oy, Width: w, Height: h})
		sel := i == *index
		if win.WidgetClick(oid, oh) && !picked {
			*index = i
			t.selectOpen = 0
			picked = true
		}
		st := State{Hovered: oh, Focused: focused, On: sel}
		if !picked && focused && t.selectCursor == i {
			st.Hovered = true
		}
		draws = append(draws, optionDraw{x: x, y: oy, label: opt, sel: sel, state: st})
	}
	// The menu is drawn dead last, above every widget the flow painted under
	// it, so an open menu never hides behind fields below — behind a scrim
	// over the whole window first, which is what tells the eye the menu is
	// the only thing asking for it right now.
	win.WidgetFrontmost(func() {
		cv := win.Canvas()
		cv.FillRect(canvas.Area{X: 0, Y: 0, Width: win.Width(), Height: win.Height()}, canvas.Fade(canvas.Black, modalScrimAlpha))
		for _, d := range draws {
			t.style.SelectOption(SelectOptionParams{
				Win:   win,
				State: d.state,
				X:     d.x, Y: d.y, W: w, H: h,
				Label:    d.label,
				Selected: d.sel,
			})
		}
	})
	return picked
}
