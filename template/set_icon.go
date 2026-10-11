package template

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/antui/svg"
)

// SetIcon registers the drawing a named icon is made of, replacing whatever was
// there, and reports a drawing that could not be read — in which case the icon
// that was registered stays as it was.
//
// The drawing is SVG text, which is the point: a developer reads the icon that
// is already there with [IconSource], changes it, and hands the changed text
// back. Recolouring it, thickening the line, squaring off the corner of a caret
// or drawing a label of one's own across the tick of a checkbox are all changes
// to the text, and none of them need anything else from the engine.
//
//	btn := template.IconSource(template.IconCheck)
//	btn = strings.Replace(btn, "stroke-width", `fill="red" stroke-width`)
//	template.SetIcon(template.IconCheck, btn)
//
// The name is what a stylesheet's `url()` and the widgets ask for, so a new
// name makes a new icon: `SetIcon("save", src)` is what a button wears when its
// stylesheet says `background-image: url(save)`.
func SetIcon(name, src string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("template: an icon needs a name")
	}
	img, err := svg.Parse(src)
	if err != nil {
		return fmt.Errorf("template: the icon %q could not be read: %w", name, err)
	}
	iconMu.Lock()
	icons[name] = &iconEntry{src: src, img: img}
	iconMu.Unlock()
	// The pictures already painted are of the icon that was there before, so they
	// go: what is on screen for a frame is the old one, which is no worse than a
	// frame painted before the icon was replaced.
	iconPaintMu.Lock()
	for key := range iconPaint {
		if key.name == name {
			delete(iconPaint, key)
		}
	}
	iconPaintMu.Unlock()
	return nil
}
