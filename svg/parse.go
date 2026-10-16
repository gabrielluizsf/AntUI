package svg

// Parse reads a drawing and answers it, or an error saying why it could not be
// read at all. The error is only ever structural: a drawing that does not
// close what it opens has no shape to paint, while anything inside it that
// cannot be drawn is left to the warnings, which are asked for afterwards.
//
// A drawing read this way came from no file, so a relative address inside it
// has nothing to resolve against and is left as it was written —
// [ParseFile] reads a file and keeps it.
//
// A drawing that is not an `<svg>` is an error, because there is nothing to
// paint; one whose size cannot be read falls back to what its viewBox says, and
// a drawing with neither is a hundred by a hundred.
func Parse(src string) (*Image, error) {
	return parseText(src, "", nil)
}
