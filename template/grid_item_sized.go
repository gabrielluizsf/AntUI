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
	nw, contribW, minContribW = c.gridItemWidth(st, nw, contribW, minContribW)
	nh, contribH, minContribH = c.gridItemHeight(st, nh, contribH, minContribH, rowDefinite)

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
	c.gridItemBounds(st, &it, rowDefinite)
	c.gridItemInsets(st, &it)
	return it
}
