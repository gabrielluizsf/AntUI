package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

func (cs *cssStyle) drawTextLine(win *antui.Window, st css.Style, x, y, w, h int, text string, justify bool) {
	cv := win.Canvas()
	o := cs.textOpts(st)
	text = transformText(text, st.TextTransform)

	padL, padR := cs.padding(st, 3), cs.padding(st, 1)
	padT, padB := cs.padding(st, 2), cs.padding(st, 0)
	avail := w - padL - padR
	th := o.Height()

	if justify {
		o.WordSpacing += justifyGap(cv, text, o, avail)
	}
	if st.TextOverflow == css.TextOverflowEllipsis {
		text = fitText(cv, text, o, avail)
	}
	tw := cv.TextWidthStyled(text, o)
	tx := x + padL
	if !justify {
		switch st.TextAlign {
		case 1:
			tx = x + w/2 - tw/2
		case 2:
			tx = x + w - padL - tw
		}
	}
	ty := y + padT + (h-padT-padB-th)/2 + cs.vAlignShift(st, th, h, padT, padB)
	if len(st.TextShadow) > 0 {
		cs.paintTextShadows(win, st, tx, ty, text, o)
	}
	o.Color = cs.ink(st)
	cv.DrawStyled(canvas.Point{X: float64(tx), Y: float64(ty)}, text, o)
	cs.decorate(cv, st, tx, ty, tw, th, o.Scale)
}
