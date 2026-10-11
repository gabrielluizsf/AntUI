package svg

// pathScanner reads the path data one thing at a time, skipping the punctuation
// between them, which is what lets the numbers be written with or without
// commas and spaces between them.
type pathScanner struct {
	src string
	pos int
}
