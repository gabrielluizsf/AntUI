package canvas

// Path is a sequence of subpaths. Each subpath is a run of points that may or
// may not be closed, and a fill treats every subpath as a boundary while a
// stroke walks each one as a line of its own. Curves are flattened into points
// when they are added, so a Path is a polygon set by the time anything reads it.
type Path struct {
	subs []subpath
}
