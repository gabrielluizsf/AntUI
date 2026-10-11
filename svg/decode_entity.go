package svg

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// decodeEntity is one name turned into the character it stands for: the five
// XML defines, or a `&#nn;` number and a `&#xhh;` one, the hex form being what
// SVG uses for the characters outside the basic plane.
func decodeEntity(name string) (rune, bool) {
	switch name {
	case "amp":
		return '&', true
	case "lt":
		return '<', true
	case "gt":
		return '>', true
	case "quot":
		return '"', true
	case "apos":
		return '\'', true
	}
	body, ok := strings.CutPrefix(name, "#")
	if !ok {
		return 0, false
	}
	digits := body
	base := 10
	if hex, isHex := strings.CutPrefix(body, "x"); isHex {
		digits, base = hex, 16
	}
	v, err := strconv.ParseInt(digits, base, 32)
	if err != nil || v <= 0 || v > utf8.MaxRune {
		return 0, false
	}
	return rune(v), true
}
