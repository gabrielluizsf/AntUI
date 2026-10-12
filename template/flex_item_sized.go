package template

import "github.com/gabrielluizsf/antui/template/css"

// flexItemSized collects one item the way the stylesheet asks for it, from a
// size the item wants: a widget's own natural box, or the box a nested
// container asked for. The properties a container reads are the same either
// way, so a grid drawn inside a flex is placed as a real item.
func (c *CSS) flexItemSized(st css.Style, nw, nh int) flexItem {
	b := c.flex
	nw, nh = c.flexItemSize(st, nw, nh)
	shrink := 1.0
	if st.Has("flex-shrink") {
		shrink = st.FlexShrink
	}
	basis := css.Auto()
	if st.Has("flex-basis") {
		basis = st.FlexBasis
	}
	self := uint8(css.AlignAuto)
	if st.Has("align-self") {
		self = st.AlignSelf
	}
	it := flexItem{
		st: st,
		nw: nw, nh: nh, w: nw, h: nh,
		ml:        c.style.length(st, st.Margin[3], b.contentW),
		mr:        c.style.length(st, st.Margin[1], b.contentW),
		mt:        c.style.length(st, st.Margin[0], b.contentW),
		mb:        c.style.length(st, st.Margin[2], b.contentW),
		order:     st.Order,
		grow:      st.FlexGrow,
		shrink:    shrink,
		basis:     basis,
		alignSelf: self,
		explicitW: st.Has("width"),
		explicitH: st.Has("height"),
		minMain:   -1,
		maxMain:   -1,
		inLeft:    -1,
		inTop:     -1,
		inRight:   -1,
		inBottom:  -1,
	}
	if st.OutOfFlow() {
		it.out = true
		c.flexItemInsets(st, &it)
		return it
	}
	c.flexItemMain(st, &it)
	return it
}
