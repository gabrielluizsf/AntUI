package template

import "github.com/gabrielluizsf/antui/template/audio"

// Builtin is the template that paints exactly like the built-in widgets — the
// light and dark themes, the built-in face — voiced by
// [github.com/gabrielluizsf/antui/template/audio.Simple]. It is the zero-cost
// starting point, and the Cyberpunk template's plainspoken neighbour.
func Builtin() Template {
	return New(builtinStyle{}, audio.Simple)
}
