package antui

import "io"

// writeCOFF writes the object file: one section, its relocations, and the two
// symbols a relocation needs to point at something.
func writeCOFF(w io.Writer, section resourceSection) error {
	return writeCOFFFor(w, "amd64", section)
}
