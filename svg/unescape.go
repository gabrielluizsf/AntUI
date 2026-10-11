package svg

import (
	"strings"
	"unicode/utf8"
)

// unescape turns the entities inside an attribute value or a piece of text into
// the characters they stand for. XML has five of its own and a drawing uses
// them constantly, and anything else it writes as `&#…;` — a number, or the
// number of a character — which is how a label with an ampersand in it or a
// path with a degree sign in a comment gets through.
//
// An entity this does not know is left exactly as it was written. Guessing at
// what `&foo;` was meant to be would turn a typo into a different drawing, and
// an unknown entity is almost always a nameless one that the reader can see.
func unescape(s string) string {
	if !strings.ContainsRune(s, '&') {
		return s
	}
	var out strings.Builder
	out.Grow(len(s))
	for i := 0; i < len(s); {
		if s[i] != '&' {
			_, size := utf8.DecodeRuneInString(s[i:])
			out.WriteString(s[i : i+size])
			i += size
			continue
		}
		end := strings.IndexByte(s[i:], ';')
		if end < 0 || end > 12 {
			// No semicolon, so this is a bare ampersand and not an entity.
			out.WriteByte('&')
			i++
			continue
		}
		name := s[i+1 : i+end]
		if r, ok := decodeEntity(name); ok {
			out.WriteRune(r)
		} else {
			out.WriteString(s[i : i+end+1])
		}
		i += end + 1
	}
	return out.String()
}
