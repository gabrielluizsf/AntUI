package css

import (
	"strings"
)

// parseBackImage reads one layer: a url() or one of the gradient functions.
func parseBackImage(part string, ctx Units) (BackImage, bool) {
	open := strings.IndexByte(part, '(')
	if open < 0 || !strings.HasSuffix(part, ")") {
		return BackImage{}, false
	}
	name := strings.TrimSpace(part[:open])
	args := strings.TrimSpace(part[open+1 : len(part)-1])
	switch name {
	case "url":
		u := strings.TrimSpace(strings.Trim(args, `"'`))
		if u == "" {
			return BackImage{}, false
		}
		return BackImage{URL: u}, true
	case "linear-gradient", "radial-gradient", "conic-gradient":
		g, ok := parseGradient(name, args, ctx)
		return BackImage{Grad: g}, ok
	}
	return BackImage{}, false
}
