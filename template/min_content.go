package template

import "github.com/gabrielluizsf/antui/template/css"

// minContent is the size a widget takes when nothing else is asked of it: the
// widest word it cannot break, plus whatever it draws either way.
func (cs *cssStyle) minContent(role, label string, st css.Style, u int) (w, h int) {
	return cs.sized(role, st, u, cs.widestWord(st, label))
}
