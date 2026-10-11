package template

import "github.com/gabrielluizsf/antui/template/css"

// col reports whether the container's main axis runs vertically.
func (b *flexBatch) col() bool {
	return b.st.FlexDirection == css.FlexDirectionColumn || b.st.FlexDirection == css.FlexDirectionColumnReverse
}
