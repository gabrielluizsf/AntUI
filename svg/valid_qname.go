package svg

import "strings"

// validQName reports whether s is a name the namespaces of a document allow:
// at most one colon, splitting a prefix from a local part, and each part a
// name of its own without one. `a:b` is such a name, while `a:b:c`, `:a` and
// `a:` are not, however much like one they look.
func validQName(s string) bool {
	if strings.Count(s, ":") > 1 {
		return false
	}
	prefix, local, found := strings.Cut(s, ":")
	if !found {
		return validNCName(s)
	}
	return validNCName(prefix) && validNCName(local)
}

// validNCName reports whether s is a name a prefix or a local part may be: one
// that opens with a letter or an underscore and holds only letters, digits,
// dashes, dots and underscores — no colon, which parts a name in two. Bytes at
// or above 0x80 are let through untouched, because a name may be written in
// any script and this is not the place to argue with one.
func validNCName(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 0x80:
			continue
		case c == ':':
			return false
		case i == 0:
			if !isNameStart(c) {
				return false
			}
		case isNameChar(c):
			continue
		default:
			return false
		}
	}
	return true
}
