package svg

import "testing"

// FuzzParse feeds the drawing reader text nobody vetted — the doctype, the
// entities, the style sheet, the paths — and paints what it read: whatever a
// drawing says must end in an error or a picture, never a panic.
func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		`<svg viewBox="0 0 10 10"><rect width="5" height="5" fill="red"/></svg>`,
		`<!DOCTYPE svg [<!ENTITY a "x">]><svg><text>&a;</text></svg>`,
		`<svg><style>rect{--c:red;fill:var(--c)}</style><defs><g id="g"><rect width="1" height="1"/></g></defs><use href="#g"/></svg>`,
		`<svg><path d="M0 0 L10 10 A5 5 0 1 1 3 3 Z" stroke="blue" stroke-dasharray="1 2"/></svg>`,
		`<svg><defs><linearGradient id="l"><stop offset="0" stop-color="red"/></linearGradient><filter id="f"><feGaussianBlur stdDeviation="2"/></filter></defs><circle r="4" fill="url(#l)" filter="url(#f)"/></svg>`,
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, src string) {
		if len(src) > 4096 {
			return
		}
		img, err := Parse(src)
		if err != nil || img == nil {
			return
		}
		img.Render(16, 16)
	})
}
