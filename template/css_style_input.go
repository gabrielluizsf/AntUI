package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) Input(p InputParams) {
	st := cs.styleFor(css.RoleInput, p.State)
	cs.paintBox(p.Win, st, p.X, p.Y, p.W, p.H)
	u := cs.u()
	padding := cs.base(st, st.Padding[css.SideLeft])
	savedClip := p.Win.Canvas().Clip
	p.Win.Canvas().SetClip(canvas.Area{X: p.X + padding, Y: p.Y + 1, Width: p.W - 2*padding, Height: p.H - 2})

	cursor := clampInt(p.Cursor, 0, len(p.Text))
	cursorPx := cs.measure(st, p.Text[:cursor])
	scroll := max(cursorPx-(p.W-2*padding), 0)
	// drawText adds the style's own left padding; pass the pre-scrolled x so
	// the glyphs land at x+padding-scroll — the same origin the caret uses.
	cs.drawText(p.Win, st, p.X-scroll, p.Y, p.W, p.H, p.Text)

	if p.State.Focused && p.Caret {
		cx := p.X + padding - scroll + cursorPx
		cy := p.Y + (p.H-cs.textHeight(st))/2
		p.Win.Canvas().FillRect(canvas.Area{X: cx, Y: cy, Width: u, Height: cs.textHeight(st)}, cs.ink(st))
	}
	p.Win.Canvas().Clip = savedClip
}
