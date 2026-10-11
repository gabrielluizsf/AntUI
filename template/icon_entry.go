package template

import "github.com/gabrielluizsf/antui/svg"

// An icon as it is registered: the text it was written in, so a developer can
// read it and change it, and the drawing that text became, read once so that
// painting it is not parsing it every frame.
type iconEntry struct {
	src string
	img *svg.Image
}
