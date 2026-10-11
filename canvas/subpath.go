package canvas

// subpath is one run of points and whether the last one joins the first.
type subpath struct {
	pts    []Point
	closed bool
}
