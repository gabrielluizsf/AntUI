package svg

import "strings"

// localName takes the namespace off a tag: the `svg:` of `svg:circle` and the
// whole of `{http://www.w3.org/2000/svg}circle` both come back as `circle`.
func localName(name string) string {
	if i := strings.LastIndexByte(name, ':'); i >= 0 {
		return name[i+1:]
	}
	if i := strings.IndexByte(name, '}'); i >= 0 {
		return name[i+1:]
	}
	return name
}
