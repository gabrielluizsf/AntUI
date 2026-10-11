package antui

import (
	"path/filepath"
)

// systemFontPaths is the fallback, for a machine where reg is not on the
// path — a stripped container, or a program run somewhere odd.
func systemFontPaths() []string {
	var out []string
	for _, dir := range fontDirs() {
		for _, name := range []string{"segoeui.ttf", "tahoma.ttf", "arial.ttf", "verdana.ttf"} {
			out = append(out, filepath.Join(dir, name))
		}
	}
	return out
}
