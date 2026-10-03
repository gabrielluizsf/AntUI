package svg

import (
	"strconv"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// readPathLength reads the total length an element says its path has, in the
// drawing's own units. A drawing writes one so that a stroke can be laid out
// along a path in numbers of its own — ten dashes across a line, whatever the
// line comes out at — rather than in the ones the geometry measures, which are
// only known once the path is built. It is a number and not a length: no unit
// goes on it, and one that is not a positive number is nothing to measure a
// path against, so it is said and left out rather than taken as zero.
func readPathLength(raw string, warn func(string, ...any)) float64 {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v <= 0 {
		warn("the path length %q is not a positive number, so the dash pattern is not scaled to it", raw)
		return 0
	}
	return v
}

// applyPathLength puts the pathLength an element wrote to work on the stroke
// that element paints: every length of the pattern — the dashes, the gaps and
// the offset into them — is multiplied by how far the path measures over what
// the drawing said it did. A line of a hundred that declares itself to be
// twenty has its pattern stretched five times, which is what keeps a pattern
// written as ten on and ten off covering half of the line however long the
// line really is.
//
// The pattern is replaced rather than changed in place, because the style it
// came in is shared with every element that inherited it: this is the stretch
// of one shape's own path, and the pattern the shapes around it are drawn with
// is the one that was written.
//
// There is nothing to do where there is nothing to do it with: no pathLength on
// the element, no pattern to stretch, and a path that measures nothing at all.
func applyPathLength(st *Style, path *canvas.Path) {
	if st.pathLength <= 0 || st.Dash == nil {
		return
	}
	by := path.Length() / st.pathLength
	if !(by > 0) {
		return
	}
	on := make([]float64, len(st.Dash.On))
	off := make([]float64, len(st.Dash.Off))
	for i, l := range st.Dash.On {
		on[i] = l * by
	}
	for i, l := range st.Dash.Off {
		off[i] = l * by
	}
	// The run of dashes and the length of a whole turn of it are worked out of
	// these when the pattern is walked, so the copy carries the lengths alone
	// and nothing left over from the pattern it came from.
	st.Dash = &canvas.Dash{On: on, Off: off, Offset: st.Dash.Offset * by}
}
