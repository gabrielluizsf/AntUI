package template

import "github.com/gabrielluizsf/antui/template/css"

// gridItemSized collects one item the way the stylesheet asks for it, from the
// size the widget wants and the narrowest it can be. The solver gets to see
// both, so a track can ask for either.
func (c *CSS) gridItemSized(role, label string, st css.Style, naturalW, naturalH, minNaturalW, minNaturalH int) gridItem {
	b := c.grid
	nw, nh := naturalW, naturalH
	contribW, contribH := naturalW, naturalH
	minContribW, minContribH := minNaturalW, minNaturalH
	rowDefinite := b.st.Has("height") && gridLengthDefinite(b.st.Height)
	if st.Has("width") && gridLengthDefinite(st.Width) {
		nw = c.style.length(st, st.Width, b.contentW)
		if !st.Width.IsPct() {
			contribW, minContribW = nw, nw
		}
	}
	if st.Has("max-width") && gridLengthDefinite(st.MaxWidth) {
		limit := c.style.length(st, st.MaxWidth, b.contentW)
		nw = min(nw, limit)
		if !st.MaxWidth.IsPct() {
			contribW = min(contribW, limit)
			minContribW = min(minContribW, limit)
		}
	}
	if st.Has("min-width") && gridLengthDefinite(st.MinWidth) {
		limit := c.style.length(st, st.MinWidth, b.contentW)
		nw = max(nw, limit)
		if !st.MinWidth.IsPct() {
			contribW = max(contribW, limit)
			minContribW = max(minContribW, limit)
		}
	}
	if st.Has("height") && gridLengthDefinite(st.Height) && (!st.Height.IsPct() || rowDefinite) {
		nh = c.style.length(st, st.Height, b.contentH)
		if !st.Height.IsPct() {
			contribH, minContribH = nh, nh
		}
	}
	if st.Has("max-height") && gridLengthDefinite(st.MaxHeight) && (!st.MaxHeight.IsPct() || rowDefinite) {
		limit := c.style.length(st, st.MaxHeight, b.contentH)
		nh = min(nh, limit)
		if !st.MaxHeight.IsPct() {
			contribH = min(contribH, limit)
			minContribH = min(minContribH, limit)
		}
	}
	if st.Has("min-height") && gridLengthDefinite(st.MinHeight) && (!st.MinHeight.IsPct() || rowDefinite) {
		limit := c.style.length(st, st.MinHeight, b.contentH)
		nh = max(nh, limit)
		if !st.MinHeight.IsPct() {
			contribH = max(contribH, limit)
			minContribH = max(minContribH, limit)
		}
	}

	justifySelf := uint8(css.AlignAuto)
	if st.Has("justify-self") {
		justifySelf = st.JustifySelf
	}
	alignSelf := uint8(css.AlignAuto)
	if st.Has("align-self") {
		alignSelf = st.AlignSelf
	}
	it := gridItem{
		role:        role,
		label:       label,
		st:          st,
		nw:          nw,
		nh:          nh,
		contribW:    contribW,
		contribH:    contribH,
		contribMinW: minContribW,
		contribMinH: minContribH,
		w:           nw,
		h:           nh,
		ml:          c.style.length(st, st.Margin[3], b.contentW),
		mr:          c.style.length(st, st.Margin[1], b.contentW),
		mt:          c.style.length(st, st.Margin[0], b.contentW),
		mb:          c.style.length(st, st.Margin[2], b.contentW),
		explicitW:   st.Has("width") && gridLengthDefinite(st.Width),
		explicitH:   st.Has("height") && gridLengthDefinite(st.Height),
		minW:        -1,
		maxW:        -1,
		minH:        -1,
		maxH:        -1,
		column:      st.GridColumn,
		row:         st.GridRow,
		justifySelf: justifySelf,
		alignSelf:   alignSelf,
		order:       st.Order,
		inLeft:      -1,
		inTop:       -1,
		inRight:     -1,
		inBottom:    -1,
	}
	if st.Has("min-width") && gridLengthDefinite(st.MinWidth) {
		it.minW = c.style.length(st, st.MinWidth, b.contentW)
	}
	if st.Has("max-width") && gridLengthDefinite(st.MaxWidth) {
		it.maxW = c.style.length(st, st.MaxWidth, b.contentW)
	}
	if st.Has("min-height") && gridLengthDefinite(st.MinHeight) && (!st.MinHeight.IsPct() || rowDefinite) {
		it.minH = c.style.length(st, st.MinHeight, b.contentH)
	}
	if st.Has("max-height") && gridLengthDefinite(st.MaxHeight) && (!st.MaxHeight.IsPct() || rowDefinite) {
		it.maxH = c.style.length(st, st.MaxHeight, b.contentH)
	}
	if st.OutOfFlow() {
		it.out = true
		if st.Has("left") {
			it.inLeft = c.style.length(st, st.Left, b.contentW)
		}
		if st.Has("right") {
			it.inRight = c.style.length(st, st.Right, b.contentW)
		}
		if st.Has("top") {
			it.inTop = c.style.length(st, st.Top, b.contentW)
		}
		if st.Has("bottom") {
			it.inBottom = c.style.length(st, st.Bottom, b.contentW)
		}
	}
	return it
}
