package svg

import "os"

// ParseFile reads a drawing from an SVG file and keeps the file as the base a
// relative address inside it resolves against: a picture written beside the
// drawing is found beside it, however the program is launched. What the file
// cannot say for itself is reported the way [Parse] reports it, and a drawing
// that could not be opened at all is an error about the opening.
//
// An `xml:base` anywhere below the root takes over from this file for that
// element and everything under it, as XML Base has it.
func ParseFile(path string) (*Image, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseText(string(data), path)
}
