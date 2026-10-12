package svg

import "github.com/gabrielluizsf/antui/canvas"

// pathState is what one command of a `d` needs to know from the ones before it:
// the path being built, the string being read, and the pen. The two curve flags
// say which kind of curve came last, because a smooth curve mirrors the control
// point of a curve of its own kind and uses the pen itself after any other
// command.
type pathState struct {
	p    *canvas.Path
	sc   *pathScanner
	warn func(string, ...any)

	cmd       byte
	cur       canvas.Point // where the pen is
	start     canvas.Point // where the current subpath began, for Z
	ctl       canvas.Point // the control point of a curve, for S and T
	prevCubic bool         // the last curve was cubic, so S may mirror
	prevQuad  bool         // the last curve was quadratic, so T may mirror
	started   bool
	rel       bool
}
