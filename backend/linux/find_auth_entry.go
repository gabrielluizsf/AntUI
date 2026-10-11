//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"encoding/binary"
	"io"
)

// findAuthEntry walks the authority file's entries until it finds the cookie
// for the wanted display. An entry with no display number applies to every
// display. Any short read ends the walk with no cookie.
func findAuthEntry(f io.Reader, wanted string) xAuth {
	for {
		var family uint16
		if err := binary.Read(f, binary.BigEndian, &family); err != nil {
			return xAuth{}
		}
		if _, err := readAuthBlock(f); err != nil { // address
			return xAuth{}
		}
		number, err := readAuthBlock(f)
		if err != nil {
			return xAuth{}
		}
		name, err := readAuthBlock(f)
		if err != nil {
			return xAuth{}
		}
		data, err := readAuthBlock(f)
		if err != nil {
			return xAuth{}
		}
		if string(name) == "MIT-MAGIC-COOKIE-1" &&
			(len(number) == 0 || string(number) == wanted) {
			return xAuth{name: string(name), data: data}
		}
	}
}
