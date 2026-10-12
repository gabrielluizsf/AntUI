package svg

// filterElement is a `<filter>` found in the drawing, with the id it is named
// by.
type filterElement struct {
	id string
	e  *element
}
