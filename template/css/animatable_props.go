package css

// animatableProps is every canonical property the engine interpolates
// between two styles. It is the cylinder a transition-property of "all" and
// a keyframe both draw from; a property absent here snaps instead. Border
// widths and the discrete layout keywords (display, position, box-sizing,
// visibility) deliberately stay away.
var animatableProps = []string{
	"background-color", "color", "opacity",
	"width", "height",
	"min-width", "max-width", "min-height", "max-height",
	"margin-top", "margin-right", "margin-bottom", "margin-left",
	"padding-top", "padding-right", "padding-bottom", "padding-left",
	"top", "right", "bottom", "left",
	"font-size", "letter-spacing", "word-spacing", "line-height",
	"transform", "transform-origin",
	"box-shadow", "text-shadow", "filter", "backdrop-filter",
}
