package template

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/svg"
)

// resetIcons puts the built-in icons back, which is what a program wants after
// it has been changing them, and what a test wants between cases.
func resetIcons() {
	iconMu.Lock()
	defer iconMu.Unlock()
	icons = make(map[string]*iconEntry, len(builtinIcons))
	for name, src := range builtinIcons {
		// A built-in icon is known to be readable, since it is the one this
		// package ships; one that is not would be a mistake here rather than
		// something to report, and it would draw as nothing.
		img, err := svg.Parse(src)
		if err != nil {
			continue
		}
		icons[name] = &iconEntry{src: src, img: img}
	}
	iconPaintMu.Lock()
	iconPaint = map[iconPaintKey]*canvas.Canvas{}
	iconPaintMu.Unlock()
}
