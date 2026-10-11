package template

import "github.com/gabrielluizsf/antui/template/css"

// widgetEntry is the timeline of one widget on the page: the interaction
// state it painted last, the style it showed when the state changed, the
// style it is heading toward, and when in the clock the change happened. It
// is what turns a hover from a jump into a glide, because the entry keeps
// hold of the whole motion between frames.
type widgetEntry struct {
	role      string
	state     State
	from      css.Style
	to        css.Style
	started   float64 // timeline seconds when the transition began
	budget    float64 // how long the longest transition runs, in milliseconds
	active    bool
	displayed css.Style
	seen      int  // the frame the widget was last drawn in
	fresh     bool // the entry has never been painted before
}
