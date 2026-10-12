package template

import "github.com/gabrielluizsf/antui/template/css"

// gridItemAlignment resolves the item's justify and align, falling back to the
// container's and then to stretch.
func (b *gridBatch) gridItemAlignment(it *gridItem) (uint8, uint8) {
	justify := b.st.JustifyItems
	if !b.st.Has("justify-items") {
		justify = css.AlignStretch
	}
	if it.justifySelf != css.AlignAuto {
		justify = it.justifySelf
	}
	align := b.st.AlignItems
	if !b.st.Has("align-items") {
		align = css.AlignStretch
	}
	if it.alignSelf != css.AlignAuto {
		align = it.alignSelf
	}
	return justify, align
}
