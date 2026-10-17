package svg

import (
	"fmt"
	"unicode/utf8"
)

// checkChars says every character of a document is one XML lets stand in it:
// the control characters XML keeps out — everything below a space but tab, a
// line break and a carriage return — and any byte that is not part of the
// UTF-8 the text is written in are errors wherever they stand, in a tag, in
// the text between tags or inside a comment. A reader that passed one on would
// paint a document that is not the one the file holds, and a character
// reference to such a character is refused where it is written for the same
// reason.
func checkChars(src string) error {
	for i := 0; i < len(src); {
		r, size := utf8.DecodeRuneInString(src[i:])
		if r == utf8.RuneError && size == 1 {
			return fmt.Errorf("antui/svg: the drawing holds a byte at offset %d that is not part of the UTF-8 it is written in", i)
		}
		if !validXMLChar(r) {
			return fmt.Errorf("antui/svg: the drawing holds %q, which XML keeps out of a document", r)
		}
		i += size
	}
	return nil
}
