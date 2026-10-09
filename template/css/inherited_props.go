package css

// inheritedProps names the properties CSS passes down the tree by default:
// the textual ones. A template's widgets have no nested parent, so the
// inheritance source is the body's computed style — the box properties are
// deliberately absent, they do not inherit in CSS either. Every property
// implemented later that inherits by default must be added here.
var inheritedProps = map[string]bool{
	"color":          true,
	"font-size":      true,
	"font-weight":    true,
	"font-style":     true,
	"font-family":    true,
	"line-height":    true,
	"letter-spacing": true,
	"word-spacing":   true,
	"text-align":     true,
	"text-transform": true,
	"white-space":    true,
	"overflow-wrap":  true,
	"visibility":     true,
	"cursor":         true,
	"text-shadow":    true,
}
