package template

import "github.com/gabrielluizsf/antui/template/css"

// outlineOn reports whether the outline paints. An outline with no width or
// style none is not drawn, the way a border with no width is not.
func (cs *cssStyle) outlineOn(st css.Style) bool {
	return st.OutlineWidth > 0 && st.OutlineStyle != css.BorderNone
}
