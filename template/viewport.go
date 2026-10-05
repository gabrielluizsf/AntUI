package template

import "github.com/gabrielluizsf/antui/template/css"

// viewport is the window the styles of this moment are read against: the two
// edges @media measures, the display's density and the color scheme the
// system paints in, and the size the viewport units resolve to. It is read
// off the window itself rather than remembered, so a resize the event loop
// delivered before the frame reached the painter — or a theme the platform
// changed — is what the whole frame draws with, from the background clear to
// the last widget, and the style cache behind GetStyleViewport is keyed by
// it, so neither ever leaves a stale answer hanging.
//
// A display that did not say how dense it is and a platform that did not say
// what scheme it is in stay as they were read: unknown, which is what keeps
// prefers-color-scheme and resolution from guessing at an answer.
func (cs *cssStyle) viewport() css.Viewport {
	vp := css.Viewport{Width: cs.win.Width(), Height: cs.win.Height()}
	if scale, ok := cs.win.DisplayScale(); ok {
		vp.Scale = scale
	}
	if dark, ok := cs.win.SystemDark(); ok {
		vp.Scheme = css.SchemeLight
		if dark {
			vp.Scheme = css.SchemeDark
		}
	}
	return vp
}
