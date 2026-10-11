package template

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/antui/template/css"
)

// cssStyle paints the CSS template: every seam is a decision from a computed
// [css.Style], and everything the stylesheet did not say falls back onto the
// window's theme. It is the [Style] behind [TemplateWithCSS].
type cssStyle struct {
	win     *antui.Window
	classes *css.CSSClasses
	images  map[string]*canvas.Canvas // background-image url() targets

	// The widget timeline: one entry per widget as it appears in a frame, so
	// a widget's state changes and their transitions persist across Reset
	// calls. Entries are keyed by role, ordinal in the frame and label;
	// entryStyle blends them while they animate.
	timeline map[string]*widgetEntry
	cur      *widgetEntry   // the widget laid out right now
	seq      int            // widgets begun so far in this frame
	frame    int            // Reset calls carried out so far
	now      func() float64 // clock in seconds, antui.Now by default

	// transformed widgets first paint to a window-sized scratch layer, then
	// composite through their matrix; considered dirty every frame.
	layer  *canvas.Canvas
	layerW int
	layerH int
}
