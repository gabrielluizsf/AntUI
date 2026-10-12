package svg

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
)

// pathData is the `d` of a `<path>`: a list of commands, each followed by the
// numbers it takes. The letters may be run together and the numbers written
// without spaces between them, as the grammar allows, so the whole string is
// scanned for both at once.
func pathData(d string, warn func(string, ...any)) *canvas.Path {
	if strings.TrimSpace(d) == "" {
		return nil
	}
	st := &pathState{p: canvas.NewPath(), sc: &pathScanner{src: d}, warn: warn}
	for {
		// Everything read is the whole of it: a command repeated by a number
		// runs out of numbers as soon as the string does.
		if st.sc.eof() {
			break
		}
		// A letter is a command; after one, its own letter repeats when the
		// next thing in the string is a number.
		if st.sc.atLetter() {
			st.cmd = st.sc.letter()
		} else if st.cmd == 0 {
			warn("the path data starts with something that is not a command")
			return nil
		} else if st.cmd == 'M' {
			// A number after a moveto is another moveto, and after anything
			// else it is a lineto, which is what SVG calls an implicit one.
			st.cmd = 'L'
		} else if st.cmd == 'm' {
			st.cmd = 'l'
		} else if st.cmd == 'Z' || st.cmd == 'z' {
			st.cmd = 0
			st.sc.skipSpace()
			if st.sc.eof() {
				break
			}
			continue
		}
		if !st.command() {
			return nil
		}
	}
	return st.p
}
