package template

import "github.com/gabrielluizsf/antui/template/css"

// alignItems is how the container places its items across the cross axis. A
// sheet that never declared align-items has the field at its zero value, which
// is not one of the Align* constants, so the initial value is asked for the way
// the initial keyword sets it.
func (b *flexBatch) alignItems() uint8 {
	if !b.st.Has("align-items") {
		return css.AlignStretch
	}
	return b.st.AlignItems
}
